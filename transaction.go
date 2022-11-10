package postgres

import "context"

type Transaction interface {
	Reader
	Writer
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
