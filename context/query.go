package context

import "context"

type PageToken string

type Iter interface {
	Next(ctx context.Context, out interface{}) bool
	Err() error
}

type PagingIter interface {
	Iter
	NextToken(ctx context.Context) (PageToken, error)
}

type Query interface {
	Filter(expr string, args ...interface{}) Query
	Limit(limit int) Query
	One(ctx context.Context, ret interface{}) (bool, error)
	All(ctx context.Context, ret interface{}) error

	// Optional future support
	// Iter() PagingIter
	// AllWithNextToken(ctx context.Context, ret interface{}) (PageToken, error)
	// Count(ctx context.Context) (int, error)
}
