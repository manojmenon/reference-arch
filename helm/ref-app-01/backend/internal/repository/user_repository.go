package repository

import (
	"context"
	"errors"
	"time"

	"github.com/enterprise/enterprise-3tier/backend/internal/domain"
	"github.com/enterprise/enterprise-3tier/backend/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository persists users via PostgreSQL.
type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create inserts a user and returns the stored row.
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return database.RetryTransient(ctx, 5, 200*time.Millisecond, func(ctx context.Context) error {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO users (id, name, email, created_at) VALUES ($1, $2, $3, $4)`,
			u.ID, u.Name, u.Email, u.CreatedAt,
		)
		return err
	})
}

// List returns all users ordered by created_at.
func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	var out []domain.User
	err := database.RetryTransient(ctx, 5, 200*time.Millisecond, func(ctx context.Context) error {
		rows, err := r.pool.Query(ctx,
			`SELECT id, name, email, created_at FROM users ORDER BY created_at DESC`,
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		out = nil
		for rows.Next() {
			var u domain.User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
				return err
			}
			out = append(out, u)
		}
		return rows.Err()
	})
	if out == nil {
		out = []domain.User{}
	}
	return out, err
}

// GetByID returns a user by primary key.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var u domain.User
	err := database.RetryTransient(ctx, 5, 200*time.Millisecond, func(ctx context.Context) error {
		return r.pool.QueryRow(ctx,
			`SELECT id, name, email, created_at FROM users WHERE id = $1`, id,
		).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	})
	return u, err
}

// Update modifies name/email for the given id.
func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, name, email *string) (domain.User, error) {
	var u domain.User
	err := database.RetryTransient(ctx, 5, 200*time.Millisecond, func(ctx context.Context) error {
		return r.pool.QueryRow(ctx, `
			UPDATE users SET
				name = COALESCE($2, name),
				email = COALESCE($3, email)
			WHERE id = $1
			RETURNING id, name, email, created_at
		`, id, name, email).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	})
	return u, err
}

// Delete removes a user by id.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return database.RetryTransient(ctx, 5, 200*time.Millisecond, func(ctx context.Context) error {
		tag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
		return nil
	})
}

// IsNotFound reports whether err is a missing row.
func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

// IsUniqueViolation reports duplicate key on unique index (e.g. email).
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
