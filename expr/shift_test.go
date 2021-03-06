package expr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShiftRegular(t *testing.T) {
	params := Map{
		"a": 4.4,
	}
	s := msgpacked(t, SHIFT(SUM(FIELD("a")), 1*time.Hour))
	b1 := make([]byte, s.EncodedWidth()*2)
	b2 := make([]byte, s.EncodedWidth()*2)
	b3 := make([]byte, s.EncodedWidth()*2)
	_, val, _ := s.Update(b1, params, nil)
	assert.EqualValues(t, 4.4, val)
	s.Update(b2, params, nil)
	val, _, _ = s.Get(b2)
	assert.EqualValues(t, 4.4, val)
	s.Merge(b3, b1, b2)
	val, _, _ = s.Get(b3)
	assert.EqualValues(t, 8.8, val)
}

func TestShiftSubMerge(t *testing.T) {
	res := 1 * time.Hour
	periods := 10

	fa := msgpacked(t, SUM(FIELD("a")))
	fs := msgpacked(t, SUB(SHIFT(SHIFT(SUM(FIELD("a")), -2*res), -1*res), SUM(FIELD("a"))))
	assert.EqualValues(t, -3*res, fs.Shift())

	a := make([]byte, fa.EncodedWidth()*periods)
	s := make([]byte, fs.EncodedWidth()*periods)

	for i := 0; i < periods; i++ {
		fa.Update(a[i*fa.EncodedWidth():], Map{"a": float64(i)}, nil)
	}

	subs := fs.SubMergers([]Expr{fa})
	for i := 0; i < periods; i++ {
		for _, sub := range subs {
			sub(s[i*fs.EncodedWidth():], a[i*fa.EncodedWidth():], res, nil)
		}
	}
	for i := 0; i < periods; i++ {
		expected := 3
		if i >= 7 {
			expected = -1 * i
		}
		actual, _, _ := fs.Get(s[i*fs.EncodedWidth():])
		assert.EqualValues(t, expected, actual, "Wrong value at position %d", i)
	}
}
