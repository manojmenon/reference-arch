package service

import (
	"testing"

	"github.com/enterprise/enterprise-3tier/backend/internal/repository"
)

func TestUserRepository_implements_UserStore(t *testing.T) {
	var _ UserStore = (*repository.UserRepository)(nil)
}
