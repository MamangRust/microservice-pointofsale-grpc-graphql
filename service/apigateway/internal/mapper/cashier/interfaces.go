package cashiergraphqlmapper

import (
	pb "github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/monolith-graphql-pointofsale-apigateway/internal/model"
)

type CashierGraphqlMapper interface {
	ToGraphqlResponseCashier(res *pb.ApiResponseCashier) *model.APIResponseCashier
	ToGraphqlResponsesCashier(res *pb.ApiResponsesCashier) *model.APIResponsesCashier
	ToGraphqlResponseCashierDeleteAt(res *pb.ApiResponseCashierDeleteAt) *model.APIResponseCashierDeleteAt
	ToGraphqlResponseCashierDelete(res *pb.ApiResponseCashierDelete) *model.APIResponseCashierDelete
	ToGraphqlResponseCashierAll(res *pb.ApiResponseCashierAll) *model.APIResponseCashierAll
	ToGraphqlResponsePaginationCashier(res *pb.ApiResponsePaginationCashier) *model.APIResponsePaginationCashier
	ToGraphqlResponsePaginationCashierDeleteAt(res *pb.ApiResponsePaginationCashierDeleteAt) *model.APIResponsePaginationCashierDeleteAt
	ToGraphqlResponseMonthlyTotalSales(res *pb.ApiResponseCashierMonthlyTotalSales) *model.APIResponseCashierMonthlyTotalSales
	ToGraphqlResponseMonthlySales(res *pb.ApiResponseCashierMonthSales) *model.APIResponseCashierMonthSales
	ToGraphqlResponseYearlySales(res *pb.ApiResponseCashierYearSales) *model.APIResponseCashierYearSales
	ToGraphqlResponseYearlyTotalSales(res *pb.ApiResponseCashierYearlyTotalSales) *model.APIResponseCashierYearlyTotalSales
}
