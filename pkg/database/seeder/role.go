package seeder

import (
	roledb "github.com/MamangRust/microservice-point-of-sale-role/database/schema"
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"

	"go.uber.org/zap"
)

type roleSeeder struct {
	roledb *roledb.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewRoleSeeder(roledb *roledb.Queries, ctx context.Context, logger logger.LoggerInterface) *roleSeeder {
	return &roleSeeder{
		roledb: roledb,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *roleSeeder) Seed() error {
	randomRoles := []string{
		"ROLE_ADMIN",
		"Admin Access 1",
		"Super Admin",
		"Admin",
		"Store Manager",
		"Cashier",
		"Inventory Staff",
		"Support",
		"Auditor",
		"Viewer",
	}

	totalRoles := len(randomRoles)

	for i, roleName := range randomRoles {
		_, err := r.roledb.CreateRole(r.ctx, roleName)
		if err != nil {
			r.logger.Error("failed to seed role", zap.Int("role", i+1), zap.String("roleName", roleName), zap.Error(err))
			return fmt.Errorf("failed to seed role %d (%s): %w", i+1, roleName, err)
		}
	}

	r.logger.Info("role seeded successfully", zap.Int("totalRoles", totalRoles))
	return nil
}
