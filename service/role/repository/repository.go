package repository

import (
	db "github.com/MamangRust/microservice-point-of-sale-role/database/schema"
)

type Repositories struct {
	RoleCommand RoleCommandRepository
	RoleQuery   RoleQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		RoleCommand: NewRoleCommandRepository(DB),
		RoleQuery:   NewRoleQueryRepository(DB),
	}
}
