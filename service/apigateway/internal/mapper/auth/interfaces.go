package authgraphqlmapper

import (
	pb "github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/monolith-graphql-pointofsale-apigateway/internal/model"
)

type AuthGraphqlMapper interface {
	ToGraphqlVerifyCode(res *pb.ApiResponseVerifyCode) *model.APIResponseVerifyCode
	ToGraphqlForgotPassword(res *pb.ApiResponseForgotPassword) *model.APIResponseForgotPassword
	ToGraphqlResetPassword(res *pb.ApiResponseResetPassword) *model.APIResponseResetPassword
	ToGraphqlResponseLogin(res *pb.ApiResponseLogin) *model.APIResponseLogin
	ToGraphqlResponseRegister(res *pb.ApiResponseRegister) *model.APIResponseRegister
	ToGraphqlResponseRefreshToken(res *pb.ApiResponseRefreshToken) *model.APIResponseRefreshToken
	ToGraphqlResponseGetMe(res *pb.ApiResponseGetMe) *model.APIResponseGetMe
}
