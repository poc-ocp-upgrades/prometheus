package storage

import (
	"context"
	"errors"
	"github.com/prometheus/prometheus/pkg/labels"
)

var (
	ErrNotFound						= errors.New("not found")
	ErrOutOfOrderSample				= errors.New("out of order sample")
	ErrDuplicateSampleForTimestamp	= errors.New("duplicate sample for timestamp")
	ErrOutOfBounds					= errors.New("out of bounds")
)

type Storage interface {
	Queryable
	StartTime() (int64, error)
	Appender() (Appender, error)
	Close() error
}
type Queryable interface {
	Querier(ctx context.Context, mint, maxt int64) (Querier, error)
}
type Querier interface {
	Select(*SelectParams, ...*labels.Matcher) (SeriesSet, Warnings, error)
	LabelValues(name string) ([]string, error)
	LabelNames() ([]string, error)
	Close() error
}
type SelectParams struct {
	Start	int64
	End		int64
	Step	int64
	Func	string
}
type QueryableFunc func(ctx context.Context, mint, maxt int64) (Querier, error)

func (f QueryableFunc) Querier(ctx context.Context, mint, maxt int64) (Querier, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return f(ctx, mint, maxt)
}

type Appender interface {
	Add(l labels.Labels, t int64, v float64) (uint64, error)
	AddFast(l labels.Labels, ref uint64, t int64, v float64) error
	Commit() error
	Rollback() error
}
type SeriesSet interface {
	Next() bool
	At() Series
	Err() error
}
type Series interface {
	Labels() labels.Labels
	Iterator() SeriesIterator
}
type SeriesIterator interface {
	Seek(t int64) bool
	At() (t int64, v float64)
	Next() bool
	Err() error
}
type Warnings []error
