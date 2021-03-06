package common

import (
	"context"
	"fmt"
	"time"

	"github.com/getlantern/bytemap"
	"github.com/getlantern/wal"

	"github.com/getlantern/zenodb/encoding"
)

const (
	keyIncludeMemStore = "zenodb.includeMemStore"

	nanosPerMilli = 1000000
)

type Partition struct {
	Keys   []string
	Tables []*PartitionTable
}

type PartitionTable struct {
	Name   string
	Offset wal.Offset
}

type Follow struct {
	Stream          string
	EarliestOffset  wal.Offset
	PartitionNumber int
	Partitions      map[string]*Partition
}

type QueryRemote func(sqlString string, includeMemStore bool, isSubQuery bool, subQueryResults [][]interface{}, onValue func(bytemap.ByteMap, []encoding.Sequence)) (hasReadResult bool, err error)

type QueryMetaData struct {
	FieldNames []string
	AsOf       time.Time
	Until      time.Time
	Resolution time.Duration
	Plan       string
}

// QueryStats captures stats about query
type QueryStats struct {
	NumPartitions           int
	NumSuccessfulPartitions int
	LowestHighWaterMark     int64
	HighestHighWaterMark    int64
	MissingPartitions       string
}

// Retriable is a marker for retriable errors
type Retriable interface {
	error

	Retriable() bool
}

type retriable struct {
	wrapped error
}

func (err *retriable) Error() string {
	return fmt.Sprintf("%v (retriable)", err.wrapped.Error())
}

func (err *retriable) Retriable() bool {
	return true
}

// MarkRetriable marks the given error as retriable
func MarkRetriable(err error) Retriable {
	return &retriable{err}
}

func WithIncludeMemStore(ctx context.Context, includeMemStore bool) context.Context {
	return context.WithValue(ctx, keyIncludeMemStore, includeMemStore)
}

func ShouldIncludeMemStore(ctx context.Context) bool {
	include := ctx.Value(keyIncludeMemStore)
	return include != nil && include.(bool)
}

func NanosToMillis(nanos int64) int64 {
	return nanos / nanosPerMilli
}

func TimeToMillis(ts time.Time) int64 {
	return NanosToMillis(ts.UnixNano())
}
