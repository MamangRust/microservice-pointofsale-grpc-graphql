package repository

import (
	db "github.com/MamangRust/microservice-point-of-sale-auth/database/schema"
)

type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		User:         NewUserRepository(DB),
		RefreshToken: NewRefreshTokenRepository(DB),
		UserRole:     NewUserRoleRepository(DB),
		Role:         NewRoleRepository(DB),
		ResetToken:   NewResetTokenRepository(DB),
	}
}
