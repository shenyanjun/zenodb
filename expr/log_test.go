package expr

import (
	"math"
	"testing"
)

func TestLN(t *testing.T) {
	doTestLog(t, LN(CONST(math.E)))
}

func TestLog2(t *testing.T) {
	doTestLog(t, LOG2(CONST(2)))
}

func TestLog10(t *testing.T) {
	doTestLog(t, LOG10(CONST(10)))
}

func doTestLog(t *testing.T, e Expr) {
	val, _, _ := msgpacked(t, e).Get(nil)
	assertFloatEquals(t, 1, val)
}
