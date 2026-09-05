package seeder

import (
	roledb "github.com/MamangRust/microservice-point-of-sale-role/database/schema"
	userdb "github.com/MamangRust/microservice-point-of-sale-user/database/schema"
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type userRoleSeeder struct {
	roledb *roledb.Queries
	userdb *userdb.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserRoleSeeder(roledb *roledb.Queries, userdb *userdb.Queries, ctx context.Context, logger logger.LoggerInterface) *userRoleSeeder {
	return &userRoleSeeder{
		roledb: roledb,
		userdb: userdb,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *userRoleSeeder) Seed() error {
	users, err := r.userdb.GetUsers(r.ctx, userdb.GetUsersParams{
		Column1: "",
		Limit:   int32(20),
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("failed to fetch users", zap.Error(err))
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	roles, err := r.roledb.GetRoles(r.ctx, roledb.GetRolesParams{
		Column1: "",
		Limit:   4,
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("failed to fetch roles", zap.Error(err))
		return fmt.Errorf("failed to fetch roles: %w", err)
	}

	if len(users) == 0 || len(roles) == 0 {
		r.logger.Debug("no users or roles available for seeding")
		return nil
	}

	rand.Seed(uint64(time.Now().UnixNano()))

	for _, user := range users {
		role := roles[rand.Intn(len(roles))]

		_, err := r.roledb.AssignRoleToUser(r.ctx, roledb.AssignRoleToUserParams{
			UserID: user.UserID,
			RoleID: role.RoleID,
		})
		if err != nil {
			r.logger.Error("failed to assign role to user", zap.String("user", user.Email), zap.String("role", role.RoleName), zap.Error(err))
			return fmt.Errorf("failed to assign role %s to user %s: %w", role.RoleName, user.Email, err)
		}
	}

	r.logger.Info("user roles assigned successfully", zap.Int("users", len(users)), zap.Int("roles", len(roles)))
	return nil
}
