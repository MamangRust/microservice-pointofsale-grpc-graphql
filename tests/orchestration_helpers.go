package tests

import (
	"bytes"
	"context"
	authdb "github.com/MamangRust/microservice-point-of-sale-auth/database/schema"
	cashierdb "github.com/MamangRust/microservice-point-of-sale-cashier/database/schema"
	categorydb "github.com/MamangRust/microservice-point-of-sale-category/database/schema"
	merchantdb "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	order_itemdb "github.com/MamangRust/microservice-point-of-sale-order-item/database/schema"
	orderdb "github.com/MamangRust/microservice-point-of-sale-order/database/schema"
	productdb "github.com/MamangRust/microservice-point-of-sale-product/database/schema"
	roledb "github.com/MamangRust/microservice-point-of-sale-role/database/schema"
	transactiondb "github.com/MamangRust/microservice-point-of-sale-transacton/database/schema"
	userdb "github.com/MamangRust/microservice-point-of-sale-user/database/schema"
	"mime/multipart"

	"github.com/MamangRust/microservice-point-of-sale-pkg/auth"
	"github.com/MamangRust/microservice-point-of-sale-pkg/hash"
	"github.com/MamangRust/microservice-point-of-sale-shared/cache"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Role
	role_cache "github.com/MamangRust/microservice-point-of-sale-role/cache"
	role_handler "github.com/MamangRust/microservice-point-of-sale-role/handler"
	role_repo "github.com/MamangRust/microservice-point-of-sale-role/repository"
	role_service "github.com/MamangRust/microservice-point-of-sale-role/service"

	// User
	user_cache "github.com/MamangRust/microservice-point-of-sale-user/cache"
	user_handler "github.com/MamangRust/microservice-point-of-sale-user/handler"
	user_repo "github.com/MamangRust/microservice-point-of-sale-user/repository"
	user_service "github.com/MamangRust/microservice-point-of-sale-user/service"

	// Auth
	auth_cache "github.com/MamangRust/microservice-point-of-sale-auth/cache"
	auth_handler "github.com/MamangRust/microservice-point-of-sale-auth/handler"
	auth_repo "github.com/MamangRust/microservice-point-of-sale-auth/repository"
	auth_service "github.com/MamangRust/microservice-point-of-sale-auth/service"

	// Category
	category_cache "github.com/MamangRust/microservice-point-of-sale-category/cache"
	category_handler "github.com/MamangRust/microservice-point-of-sale-category/handler"
	category_repo "github.com/MamangRust/microservice-point-of-sale-category/repository"
	category_service "github.com/MamangRust/microservice-point-of-sale-category/service"

	// Product
	product_cache "github.com/MamangRust/microservice-point-of-sale-product/cache"
	product_handler "github.com/MamangRust/microservice-point-of-sale-product/handler"
	product_repo "github.com/MamangRust/microservice-point-of-sale-product/repository"
	product_service "github.com/MamangRust/microservice-point-of-sale-product/service"

	// Merchant
	merchant_cache "github.com/MamangRust/microservice-point-of-sale-merchant/cache"
	merchant_handler "github.com/MamangRust/microservice-point-of-sale-merchant/handler"
	merchant_repo "github.com/MamangRust/microservice-point-of-sale-merchant/repository"
	merchant_service "github.com/MamangRust/microservice-point-of-sale-merchant/service"

	// Order
	order_cache "github.com/MamangRust/microservice-point-of-sale-order/cache"
	order_handler "github.com/MamangRust/microservice-point-of-sale-order/handler"
	order_repo "github.com/MamangRust/microservice-point-of-sale-order/repository"
	order_service "github.com/MamangRust/microservice-point-of-sale-order/service"

	// Transaction
	transaction_cache "github.com/MamangRust/microservice-point-of-sale-transacton/cache"
	transaction_handler "github.com/MamangRust/microservice-point-of-sale-transacton/handler"
	transaction_repo "github.com/MamangRust/microservice-point-of-sale-transacton/repository"
	transaction_service "github.com/MamangRust/microservice-point-of-sale-transacton/service"

	// Order Item
	order_item_cache "github.com/MamangRust/microservice-point-of-sale-order-item/cache"
	order_item_handler "github.com/MamangRust/microservice-point-of-sale-order-item/handler"
	order_item_repo "github.com/MamangRust/microservice-point-of-sale-order-item/repository"
	order_item_service "github.com/MamangRust/microservice-point-of-sale-order-item/service"

	// Cashier
	mencache "github.com/MamangRust/microservice-point-of-sale-cashier/cache"
	cashier_handler "github.com/MamangRust/microservice-point-of-sale-cashier/handler"
	cashier_repo "github.com/MamangRust/microservice-point-of-sale-cashier/repository"
	cashier_service "github.com/MamangRust/microservice-point-of-sale-cashier/service"
)

func (s *BaseTestSuite) SetupRoleService() {
	cacheStore := s.GetCacheStore()
	queries := roledb.New(s.ts.DBPool())

	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(queries)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repositories:  roleRepos,
		Logger:        s.Log,
		Mencache:      roleMencache,
		Observability: s.Obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterRoleServiceServer(server, roleGapi.Role)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["role"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupUserService() {
	cacheStore := s.GetCacheStore()
	queries := userdb.New(s.ts.DBPool())
	hasher := hash.NewHashingPassword()

	userMencache := user_cache.NewMencache(cacheStore)
	userRepos := user_repo.NewRepositories(queries)
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        s.Log,
		Hash:          hasher,
		Mencache:      userMencache,
		Observability: s.Obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, userGapi.User)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["user"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupAuthService() {
	cacheStore := s.GetCacheStore()
	queries := authdb.New(s.ts.DBPool())
	hasher := hash.NewHashingPassword()
	tokenManager, _ := auth.NewManager("mysecret")

	authRepos := auth_repo.NewRepositories(queries)
	authMencache := auth_cache.NewMencache(cacheStore)
	authSvc := auth_service.NewService(&auth_service.Deps{
		Repositories:  authRepos,
		Logger:        s.Log,
		Mencache:      authMencache,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: s.Obs,
	})
	authGapi := auth_handler.NewAuthHandleGrpc(authSvc, s.Log)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["auth"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCategoryService() {
	cacheStore := s.GetCacheStore()
	queries := categorydb.New(s.ts.DBPool())

	catMencache := category_cache.NewMencache(cacheStore)
	catRepos := category_repo.NewRepositories(queries)
	catSvc := category_service.NewService(&category_service.Deps{
		Mencache:      catMencache,
		Repositories:  catRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	catGapi := category_handler.NewHandler(&category_handler.Deps{
		Service: catSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterCategoryServiceServer(server, catGapi.Category)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["category"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupProductService() {
	cacheStore := s.GetCacheStore()
	queries := productdb.New(s.ts.DBPool())

	prodMencache := product_cache.NewMencache(cacheStore)
	prodRepos := product_repo.NewRepositories(queries)
	prodSvc := product_service.NewService(&product_service.Deps{
		Mencache:      prodMencache,
		Repositories:  prodRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	prodGapi := product_handler.NewHandler(&product_handler.Deps{
		Service: prodSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, prodGapi.Product)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["product"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupMerchantService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}

	cacheStore := s.GetCacheStore()
	queries := merchantdb.New(s.ts.DBPool())

	merchantMencache := merchant_cache.NewMencache(cacheStore)
	userQueryClient := pb.NewUserServiceClient(s.Conns["user"])
	merchantRepos := merchant_repo.NewRepositories(queries, userQueryClient)
	merchantSvc := merchant_service.NewService(&merchant_service.Deps{
		Mencache:      merchantMencache,
		Repositories:  merchantRepos,
		Logger:        s.Log,
		Observability: s.Obs,
		Kafka:         nil,
	})
	merchantGapi := merchant_handler.NewHandler(&merchant_handler.Deps{
		Service: merchantSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterMerchantServiceServer(server, merchantGapi.Merchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["merchant"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderService() {
	if _, ok := s.Conns["cashier"]; !ok {
		s.SetupCashierService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}
	if _, ok := s.Conns["product"]; !ok {
		s.SetupProductService()
	}
	if _, ok := s.Conns["order-item"]; !ok {
		s.SetupOrderItemService()
	}

	cacheStore := s.GetCacheStore()
	queries := orderdb.New(s.ts.DBPool())

	cashierClient := pb.NewCashierServiceClient(s.Conns["cashier"])
	merchantClient := pb.NewMerchantServiceClient(s.Conns["merchant"])
	productClient := pb.NewProductServiceClient(s.Conns["product"])
	orderItemClient := pb.NewOrderItemServiceClient(s.Conns["order-item"])

	orderMencache := order_cache.NewMencache(cacheStore)
	orderRepos := order_repo.NewRepositories(queries, cashierClient, merchantClient, productClient, orderItemClient)
	orderSvc := order_service.NewService(&order_service.Deps{
		Mencache:      orderMencache,
		Repositories:  orderRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	orderGapi := order_handler.NewHandler(&order_handler.Deps{
		Service: orderSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, orderGapi.Order)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupTransactionService() {
	if _, ok := s.Conns["cashier"]; !ok {
		s.SetupCashierService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}
	if _, ok := s.Conns["order"]; !ok {
		s.SetupOrderService()
	}
	if _, ok := s.Conns["order-item"]; !ok {
		s.SetupOrderItemService()
	}

	cacheStore := s.GetCacheStore()
	queries := transactiondb.New(s.ts.DBPool())

	cashierClient := pb.NewCashierServiceClient(s.Conns["cashier"])
	merchantClient := pb.NewMerchantServiceClient(s.Conns["merchant"])
	orderClient := pb.NewOrderServiceClient(s.Conns["order"])
	orderItemClient := pb.NewOrderItemServiceClient(s.Conns["order-item"])

	transactionMencache := transaction_cache.NewMencache(cacheStore)
	transactionRepos := transaction_repo.NewRepositories(queries, cashierClient, merchantClient, orderClient, orderItemClient)
	transactionSvc := transaction_service.NewService(&transaction_service.Deps{
		Mencache:      transactionMencache,
		Repositories:  transactionRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	transactionGapi := transaction_handler.NewHandler(&transaction_handler.Deps{
		Service: transactionSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterTransactionServiceServer(server, transactionGapi.Transaction)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["transaction"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderItemService() {
	cacheStore := s.GetCacheStore()
	queries := order_itemdb.New(s.ts.DBPool())

	itemMencache := order_item_cache.NewMencache(cacheStore)
	itemRepos := order_item_repo.NewRepositories(queries)
	itemSvc := order_item_service.NewService(&order_item_service.Deps{
		Mencache:      itemMencache,
		Repositories:  itemRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	itemGapi := order_item_handler.NewHandler(&order_item_handler.Deps{
		Service: itemSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterOrderItemServiceServer(server, itemGapi.OrderItem)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order-item"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCashierService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}

	cacheStore := s.GetCacheStore()
	queries := cashierdb.New(s.ts.DBPool())

	cashierMencache := mencache.NewMencache(cacheStore)
	userClient := pb.NewUserServiceClient(s.Conns["user"])
	merchantClient := pb.NewMerchantServiceClient(s.Conns["merchant"])
	cashierRepos := cashier_repo.NewRepositories(queries, userClient, merchantClient)
	cashierSvc := cashier_service.NewService(&cashier_service.Deps{
		Ctx:           context.Background(),
		Mencache:      cashierMencache,
		Repositories:  cashierRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	cashierGapi := cashier_handler.NewHandler(&cashier_handler.Deps{
		Service: cashierSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterCashierServiceServer(server, cashierGapi.Cashier)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["cashier"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) dial(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) GetCacheStore() *cache.CacheStore {
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	return cache.NewCacheStore(s.ts.RedisClient(), s.Log, cacheMetrics)
}

func (s *BaseTestSuite) BuildMultipartRequestBody(fields map[string]string, fieldName, fileName string) ([]byte, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range fields {
		fw, _ := w.CreateFormField(key)
		fw.Write([]byte(r))
	}
	fw, _ := w.CreateFormFile(fieldName, fileName)
	fw.Write([]byte("dummy image content"))
	w.Close()
	return b.Bytes(), w.FormDataContentType()
}
