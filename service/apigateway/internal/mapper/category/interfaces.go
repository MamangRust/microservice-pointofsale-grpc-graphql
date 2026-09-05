package categorygraphqlmapper

import (
	pb "github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/monolith-graphql-pointofsale-apigateway/internal/model"
)

type CategoryGraphqlMapper interface {
	ToGraphqlResponseCategory(res *pb.ApiResponseCategory) *model.APIResponseCategory
	ToGraphqlResponsesCategory(res *pb.ApiResponsesCategory) *model.APIResponsesCategory
	ToGraphqlResponseCategoryDeleteAt(res *pb.ApiResponseCategoryDeleteAt) *model.APIResponseCategoryDeleteAt
	ToGraphqlResponseCategoryDelete(res *pb.ApiResponseCategoryDelete) *model.APIResponseCategoryDelete
	ToGraphqlResponseCategoryAll(res *pb.ApiResponseCategoryAll) *model.APIResponseCategoryAll
	ToGraphqlResponsePaginationCategory(res *pb.ApiResponsePaginationCategory) *model.APIResponsePaginationCategory
	ToGraphqlResponseCategoryMonthlyTotalPrice(res *pb.ApiResponseCategoryMonthlyTotalPrice) *model.APIResponseCategoryMonthlyTotalPrice
	ToGraphqlResponseCategoryYearlyTotalPrice(res *pb.ApiResponseCategoryYearlyTotalPrice) *model.APIResponseCategoryYearlyTotalPrice
	ToGraphqlResponseCategoryMonthlyPrice(res *pb.ApiResponseCategoryMonthPrice) *model.APIResponseCategoryMonthPrice
	ToGraphqlResponseCategoryYearlyPrice(res *pb.ApiResponseCategoryYearPrice) *model.APIResponseCategoryYearPrice
	ToGraphqlResponsePaginationCategoryDeleteAt(res *pb.ApiResponsePaginationCategoryDeleteAt) *model.APIResponsePaginationCategoryDeleteAt
}
