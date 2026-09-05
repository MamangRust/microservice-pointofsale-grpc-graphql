package ordergraphqlmapper

import (
	pb "github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/monolith-graphql-pointofsale-apigateway/internal/model"
)

type OrderGraphqlMapper interface {
	ToGraphqlResponseOrder(res *pb.ApiResponseOrder) *model.APIResponseOrder
	ToGraphqlResponsesOrder(res *pb.ApiResponsesOrder) *model.APIResponsesOrder
	ToGraphqlResponseOrderDeleteAt(res *pb.ApiResponseOrderDeleteAt) *model.APIResponseOrderDeleteAt
	ToGraphqlResponseOrderDelete(res *pb.ApiResponseOrderDelete) *model.APIResponseOrderDelete
	ToGraphqlResponseOrderAll(res *pb.ApiResponseOrderAll) *model.APIResponseOrderAll
	ToGraphqlResponsePaginationOrder(res *pb.ApiResponsePaginationOrder) *model.APIResponsePaginationOrder
	ToGraphqlResponsePaginationOrderDeleteAt(res *pb.ApiResponsePaginationOrderDeleteAt) *model.APIResponsePaginationOrderDeleteAt
	ToGraphqlResponseMonthlyRevenue(res *pb.ApiResponseOrderMonthly) *model.APIResponseOrderMonthly
	ToGraphqlResponseYearlyRevenue(res *pb.ApiResponseOrderYearly) *model.APIResponseOrderYearly
	ToGraphqlResponseMonthlyTotalRevenue(res *pb.ApiResponseOrderMonthlyTotalRevenue) *model.APIResponseOrderMonthlyTotalRevenue
	ToGraphqlResponseYearlyTotalRevenue(res *pb.ApiResponseOrderYearlyTotalRevenue) *model.APIResponseOrderYearlyTotalRevenue
}
