package service

import (
	"context"
	"testing"
	"time"

	"github.com/enterprise/enterprise-3tier/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeRepo struct {
	users map[uuid.UUID]domain.User
	err   error
	dup   bool
}

func (f *fakeRepo) Create(ctx context.Context, u domain.User) error {
	if f.dup {
		return &pgconn.PgError{Code: "23505"}
	}
	if f.err != nil {
		return f.err
	}
	if f.users == nil {
		f.users = map[uuid.UUID]domain.User{}
	}
	f.users[u.ID] = u
	return nil
}

func (f *fakeRepo) List(ctx context.Context) ([]domain.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	out := make([]domain.User, 0, len(f.users))
	for _, u := range f.users {
		out = append(out, u)
	}
	return out, nil
}

func (f *fakeRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	if f.err != nil {
		return domain.User{}, f.err
	}
	u, ok := f.users[id]
	if !ok {
		return domain.User{}, pgx.ErrNoRows
	}
	return u, nil
}

func (f *fakeRepo) Update(ctx context.Context, id uuid.UUID, name, email *string) (domain.User, error) {
	if f.dup {
		return domain.User{}, &pgconn.PgError{Code: "23505"}
	}
	if f.err != nil {
		return domain.User{}, f.err
	}
	u, ok := f.users[id]
	if !ok {
		return domain.User{}, pgx.ErrNoRows
	}
	if name != nil {
		u.Name = *name
	}
	if email != nil {
		u.Email = *email
	}
	f.users[id] = u
	return u, nil
}

func (f *fakeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if f.err != nil {
		return f.err
	}
	if _, ok := f.users[id]; !ok {
		return pgx.ErrNoRows
	}
	delete(f.users, id)
	return nil
}

func TestUserService_CreateAndGet(t *testing.T) {
	fr := &fakeRepo{users: map[uuid.UUID]domain.User{}}
	svc := NewUserService(fr)

	in := domain.CreateUserInput{Name: "Ada", Email: "ada@example.com"}
	u, err := svc.Create(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != in.Name || u.Email != in.Email {
		t.Fatalf("user mismatch: %+v", u)
	}

	got, err := svc.Get(context.Background(), u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != u.ID {
		t.Fatalf("expected id %v, got %v", u.ID, got.ID)
	}
}

func TestUserService_Get_NotFound(t *testing.T) {
	svc := NewUserService(&fakeRepo{users: map[uuid.UUID]domain.User{}})
	_, err := svc.Get(context.Background(), uuid.New())
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	svc := NewUserService(&fakeRepo{dup: true})
	_, err := svc.Create(context.Background(), domain.CreateUserInput{Name: "A", Email: "a@b.c"})
	if err != ErrEmailTaken {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserService_Update(t *testing.T) {
	id := uuid.New()
	fr := &fakeRepo{users: map[uuid.UUID]domain.User{
		id: {ID: id, Name: "Old", Email: "old@example.com", CreatedAt: time.Now().UTC()},
	}}
	svc := NewUserService(fr)
	n := "New"
	u, err := svc.Update(context.Background(), id, domain.UpdateUserInput{Name: &n})
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "New" {
		t.Fatalf("expected name New, got %q", u.Name)
	}
}
