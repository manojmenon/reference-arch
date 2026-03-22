package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// RetryTransient runs fn and retries on common transient errors.
func RetryTransient(ctx context.Context, max int, backoff time.Duration, fn func(ctx context.Context) error) error {
	var last error
	for attempt := 0; attempt <= max; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		last = fn(ctx)
		if last == nil {
			return nil
		}
		if attempt == max {
			break
		}
		if !isTransient(last) {
			return last
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff * time.Duration(attempt+1)):
		}
	}
	return last
}

func isTransient(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "40001", "40P01", "08000", "08003", "08006", "57P01", "57P02", "57P03":
			return true
		}
	}
	return false
}
