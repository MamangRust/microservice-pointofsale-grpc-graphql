package orderitemgraphqlmapper

import (
	pb "github.com/MamangRust/microservice-point-of-sale-shared/pb"
	"github.com/MamangRust/monolith-graphql-pointofsale-apigateway/internal/model"
)

type OrderItemGraphqlMapper interface {
	ToGraphqlResponseOrderItem(res *pb.ApiResponseOrderItem) *model.APIResponseOrderItem
	ToGraphqlResponsesOrderItem(res *pb.ApiResponsesOrderItem) *model.APIResponsesOrderItem
	ToGrapqhlResponseOrderItemDelete(res *pb.ApiResponseOrderItemDelete) *model.APIResponseOrderItemDelete
	ToGrapqhlResponseOrderItemAll(res *pb.ApiResponseOrderItemAll) *model.APIResponseOrderItemAll
	ToGraphqlResponsePaginationOrderItem(res *pb.ApiResponsePaginationOrderItem) *model.APIResponsePaginationOrderItem
	ToGraphqlResponsePaginationOrderItemDeleteAt(res *pb.ApiResponsePaginationOrderItemDeleteAt) *model.APIResponsePaginationOrderItemDeleteAt
}
