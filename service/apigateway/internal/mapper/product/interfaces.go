package productgraphqlmapper

import (
	pb "github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/monolith-graphql-pointofsale-apigateway/internal/model"
)

type ProductGraphqlMapper interface {
	ToGraphqlResponseProduct(res *pb.ApiResponseProduct) *model.APIResponseProduct
	ToGraphqlResponsesProduct(res *pb.ApiResponsesProduct) *model.APIResponsesProduct
	ToGraphqlResponseProductDeleteAt(res *pb.ApiResponseProductDeleteAt) *model.APIResponseProductDeleteAt
	ToGraphqlResponseProductDelete(res *pb.ApiResponseProductDelete) *model.APIResponseProductDelete
	ToGraphqlResponseProductAll(res *pb.ApiResponseProductAll) *model.APIResponseProductAll
	ToGraphqlResponsePaginationProduct(res *pb.ApiResponsePaginationProduct) *model.APIResponsePaginationProduct
	ToGraphqlResponsePaginationProductDeleteAt(res *pb.ApiResponsePaginationProductDeleteAt) *model.APIResponsePaginationProductDeleteAt
}
