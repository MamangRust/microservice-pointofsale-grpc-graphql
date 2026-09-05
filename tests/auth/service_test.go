package auth_test

import (
	"context"
	authdb "github.com/MamangRust/microservice-point-of-sale-auth/database/schema"
	"testing"

	mencache "github.com/MamangRust/microservice-point-of-sale-auth/cache"
	"github.com/MamangRust/microservice-point-of-sale-auth/repository"
	"github.com/MamangRust/microservice-point-of-sale-auth/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/auth"
	"github.com/MamangRust/microservice-point-of-sale-pkg/hash"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/cache"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	tests "github.com/MamangRust/microservice-point-of-sale-test"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	pool        *pgxpool.Pool
	authService *service.Service
	email       string
	password    string
}

func (s *AuthServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.pool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(opts)

	authQueries := authdb.New(pool)
	repos := repository.NewRepositories(authQueries)

	log, _ := logger.NewLogger("test", nil)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)
	mencacheService := mencache.NewMencache(cacheStore)

	tokenManager, _ := auth.NewManager("mysecretkey")

	obs, _ := observability.NewObservability("test", log)
	s.authService = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencacheService,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: obs,
	})

	s.email = "auth.service.test@example.com"
	s.password = "password123"
}

func (s *AuthServiceTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *AuthServiceTestSuite) TestAuthLifecycle() {
	ctx := context.Background()

	// 1. Register
	regReq := &requests.RegisterRequest{
		FirstName:       "Auth",
		LastName:        "Service",
		Email:           s.email,
		Password:        s.password,
		ConfirmPassword: s.password,
	}

	created, err := s.authService.Register.Register(ctx, regReq)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	s.Equal(s.email, created.Email)

	// 1b. Verify email (login hanya menerima user is_verified = true)
	var verifyCode string
	err = s.pool.QueryRow(ctx,
		"SELECT verification_code FROM users WHERE email = $1", s.email).Scan(&verifyCode)
	s.Require().NoError(err)

	verified, err := s.authService.PasswordReset.VerifyCode(ctx, verifyCode)
	s.Require().NoError(err)
	s.True(verified)

	// 2. Login
	loginReq := &requests.AuthRequest{
		Email:    s.email,
		Password: s.password,
	}

	tokenRes, err := s.authService.Login.Login(ctx, loginReq)
	s.Require().NoError(err)
	s.Require().NotNil(tokenRes)
	s.NotEmpty(tokenRes.AccessToken)
	s.NotEmpty(tokenRes.RefreshToken)

	// 3. ForgotPassword
	success, err := s.authService.PasswordReset.ForgotPassword(ctx, s.email)
	s.Require().NoError(err)
	s.True(success)
}

func TestAuthServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthServiceTestSuite))
}
