package protomapper

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"
)

type AuthProtoMapper interface {
	ToProtoResponseVerifyCode(status string, message string) *pb.ApiResponseVerifyCode
	ToProtoResponseForgotPassword(status string, message string) *pb.ApiResponseForgotPassword
	ToProtoResponseResetPassword(status string, message string) *pb.ApiResponseResetPassword
	ToProtoResponseLogin(status string, message string, response *response.TokenResponse) *pb.ApiResponseLogin
	ToProtoResponseRegister(status string, message string, response *response.UserResponse) *pb.ApiResponseRegister
	ToProtoResponseRefreshToken(status string, message string, response *response.TokenResponse) *pb.ApiResponseRefreshToken
	ToProtoResponseGetMe(status string, message string, response *response.UserResponse) *pb.ApiResponseGetMe
}

type UserProtoMapper interface {
	ToProtoResponseUserDeleteAt(status string, message string, pbResponse *response.UserResponseDeleteAt) *pb.ApiResponseUserDeleteAt
	ToProtoResponsesUser(status string, message string, pbResponse []*response.UserResponse) *pb.ApiResponsesUser
	ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pb.ApiResponseUser
	ToProtoResponseUserDelete(status string, message string) *pb.ApiResponseUserDelete
	ToProtoResponseUserAll(status string, message string) *pb.ApiResponseUserAll
	ToProtoResponsePaginationUserDeleteAt(pagination *pb.PaginationMeta, status string, message string, users []*response.UserResponseDeleteAt) *pb.ApiResponsePaginationUserDeleteAt
	ToProtoResponsePaginationUser(pagination *pb.PaginationMeta, status string, message string, users []*response.UserResponse) *pb.ApiResponsePaginationUser
}

type RoleProtoMapper interface {
	ToProtoResponseRoleAll(status string, message string) *pb.ApiResponseRoleAll
	ToProtoResponseRoleDelete(status string, message string) *pb.ApiResponseRoleDelete
	ToProtoResponseRole(status string, message string, pbResponse *response.RoleResponse) *pb.ApiResponseRole
	ToProtoResponsesRole(status string, message string, pbResponse []*response.RoleResponse) *pb.ApiResponsesRole
	ToProtoResponsePaginationRole(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.RoleResponse) *pb.ApiResponsePaginationRole
	ToProtoResponsePaginationRoleDeleteAt(pagination *pb.PaginationMeta, status string, message string, pbResponse []*response.RoleResponseDeleteAt) *pb.ApiResponsePaginationRoleDeleteAt
}

type CategoryProtoMapper interface {
	ToProtoResponseMonthlyTotalPrice(status string, message string, row []*response.CategoriesMonthlyTotalPriceResponse) *pb.ApiResponseCategoryMonthlyTotalPrice
	ToProtoResponseYearlyTotalPrice(status string, message string, row []*response.CategoriesYearlyTotalPriceResponse) *pb.ApiResponseCategoryYearlyTotalPrice
	ToProtoResponseCategoryMonthlyPrice(status string, message string, row []*response.CategoryMonthPriceResponse) *pb.ApiResponseCategoryMonthPrice
	ToProtoResponseCategoryYearlyPrice(status string, message string, row []*response.CategoryYearPriceResponse) *pb.ApiResponseCategoryYearPrice

	ToProtoResponsesCategory(status string, message string, pbResponse []*response.CategoryResponse) *pb.ApiResponsesCategory
	ToProtoResponseCategoryDeleteAt(status string, message string, pbResponse *response.CategoryResponseDeleteAt) *pb.ApiResponseCategoryDeleteAt

	ToProtoResponseCategoryAll(status string, message string) *pb.ApiResponseCategoryAll
	ToProtoResponseCategory(status string, message string, pbResponse *response.CategoryResponse) *pb.ApiResponseCategory
	ToProtoResponseCategoryDelete(status string, message string) *pb.ApiResponseCategoryDelete
	ToProtoResponsePaginationCategoryDeleteAt(pagination *pb.PaginationMeta, status string, message string, categories []*response.CategoryResponseDeleteAt) *pb.ApiResponsePaginationCategoryDeleteAt
	ToProtoResponsePaginationCategory(pagination *pb.PaginationMeta, status string, message string, categories []*response.CategoryResponse) *pb.ApiResponsePaginationCategory
}

type CashierProtoMapper interface {
	ToProtoMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthTotalSales) *pb.ApiResponseCashierMonthlyTotalSales
	ToProtoYearlyTotalSales(status, message string, row []*response.CashierResponseYearTotalSales) *pb.ApiResponseCashierYearlyTotalSales

	ToProtoResponseMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthSales) *pb.ApiResponseCashierMonthSales
	ToProtoResponseYearlyTotalSales(status, message string, row []*response.CashierResponseYearSales) *pb.ApiResponseCashierYearSales

	ToProtoResponseCashier(status string, message string, pbResponse *response.CashierResponse) *pb.ApiResponseCashier
	ToProtoResponseCashierDeleteAt(status string, message string, pbResponse *response.CashierResponseDeleteAt) *pb.ApiResponseCashierDeleteAt
	ToProtoResponsesCashier(status string, message string, pbResponse []*response.CashierResponse) *pb.ApiResponsesCashier
	ToProtoResponseCashierDelete(status string, message string) *pb.ApiResponseCashierDelete
	ToProtoResponseCashierAll(status string, message string) *pb.ApiResponseCashierAll
	ToProtoResponsePaginationCashierDeleteAt(pagination *pb.PaginationMeta, status string, message string, users []*response.CashierResponseDeleteAt) *pb.ApiResponsePaginationCashierDeleteAt
	ToProtoResponsePaginationCashier(pagination *pb.PaginationMeta, status string, message string, users []*response.CashierResponse) *pb.ApiResponsePaginationCashier
}

type MerchantProtoMapper interface {
	ToProtoResponseMerchant(status string, message string, pbResponse *response.MerchantResponse) *pb.ApiResponseMerchant
	ToProtoResponseMerchantDeleteAt(status string, message string, pbResponse *response.MerchantResponseDeleteAt) *pb.ApiResponseMerchantDeleteAt

	ToProtoResponsesMerchant(status string, message string, pbResponse []*response.MerchantResponse) *pb.ApiResponsesMerchant
	ToProtoResponseMerchantDelete(status string, message string) *pb.ApiResponseMerchantDelete
	ToProtoResponseMerchantAll(status string, message string) *pb.ApiResponseMerchantAll
	ToProtoResponsePaginationMerchantDeleteAt(pagination *pb.PaginationMeta, status string, message string, merchants []*response.MerchantResponseDeleteAt) *pb.ApiResponsePaginationMerchantDeleteAt
	ToProtoResponsePaginationMerchant(pagination *pb.PaginationMeta, status string, message string, merchants []*response.MerchantResponse) *pb.ApiResponsePaginationMerchant
}

type MerchantDocumentProtoMapper interface {
	ToProtoResponseMerchantDocument(status string, message string, doc *response.MerchantDocumentResponse) *pb.ApiResponseMerchantDocument
	ToProtoResponsesMerchantDocument(status string, message string, docs []*response.MerchantDocumentResponse) *pb.ApiResponsesMerchantDocument

	ToProtoResponsePaginationMerchantDocument(pagination *pb.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponse) *pb.ApiResponsePaginationMerchantDocument
	ToProtoResponsePaginationMerchantDocumentDeleteAt(pagination *pb.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponseDeleteAt) *pb.ApiResponsePaginationMerchantDocumentAt

	ToProtoResponseMerchantDocumentDelete(status string, message string) *pb.ApiResponseMerchantDocumentDelete

	ToProtoResponseMerchantDocumentAll(status string, message string) *pb.ApiResponseMerchantDocumentAll
}

type OrderItemProtoMapper interface {
	ToProtoResponseOrderItem(status string, message string, pbResponse *response.OrderItemResponse) *pb.ApiResponseOrderItem
	ToProtoResponsesOrderItem(status string, message string, pbResponse []*response.OrderItemResponse) *pb.ApiResponsesOrderItem
	ToProtoResponseOrderItemDelete(status string, message string) *pb.ApiResponseOrderItemDelete
	ToProtoResponseOrderItemAll(status string, message string) *pb.ApiResponseOrderItemAll
	ToProtoResponsePaginationOrderItemDeleteAt(pagination *pb.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponseDeleteAt) *pb.ApiResponsePaginationOrderItemDeleteAt
	ToProtoResponsePaginationOrderItem(pagination *pb.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponse) *pb.ApiResponsePaginationOrderItem
}

type OrderProtoMapper interface {
	ToProtoResponseMonthlyTotalRevenue(status string, message string, row []*response.OrderMonthlyTotalRevenueResponse) *pb.ApiResponseOrderMonthlyTotalRevenue
	ToProtoResponseYearlyTotalRevenue(status string, message string, row []*response.OrderYearlyTotalRevenueResponse) *pb.ApiResponseOrderYearlyTotalRevenue

	ToProtoResponseMonthlyRevenue(status string, message string, row []*response.OrderMonthlyResponse) *pb.ApiResponseOrderMonthly
	ToProtoResponseYearlyRevenue(status string, message string, row []*response.OrderYearlyResponse) *pb.ApiResponseOrderYearly

	ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pb.ApiResponseOrder
	ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pb.ApiResponseOrderDeleteAt
	ToProtoResponsesOrder(status string, message string, pbResponse []*response.OrderResponse) *pb.ApiResponsesOrder
	ToProtoResponseOrderDelete(status string, message string) *pb.ApiResponseOrderDelete
	ToProtoResponseOrderAll(status string, message string) *pb.ApiResponseOrderAll
	ToProtoResponsePaginationOrderDeleteAt(pagination *pb.PaginationMeta, status string, message string, orders []*response.OrderResponseDeleteAt) *pb.ApiResponsePaginationOrderDeleteAt
	ToProtoResponsePaginationOrder(pagination *pb.PaginationMeta, status string, message string, orders []*response.OrderResponse) *pb.ApiResponsePaginationOrder
}

type ProductProtoMapper interface {
	ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pb.ApiResponseProduct
	ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pb.ApiResponseProductDeleteAt

	ToProtoResponsesProduct(status string, message string, pbResponse []*response.ProductResponse) *pb.ApiResponsesProduct
	ToProtoResponseProductDelete(status string, message string) *pb.ApiResponseProductDelete
	ToProtoResponseProductAll(status string, message string) *pb.ApiResponseProductAll
	ToProtoResponsePaginationProductDeleteAt(pagination *pb.PaginationMeta, status string, message string, products []*response.ProductResponseDeleteAt) *pb.ApiResponsePaginationProductDeleteAt
	ToProtoResponsePaginationProduct(pagination *pb.PaginationMeta, status string, message string, products []*response.ProductResponse) *pb.ApiResponsePaginationProduct
}

type TransactionProtoMapper interface {
	ToProtoResponseMonthAmountSuccess(status string, message string, row []*response.TransactionMonthlyAmountSuccessResponse) *pb.ApiResponseTransactionMonthAmountSuccess
	ToProtoResponseYearAmountSuccess(status string, message string, row []*response.TransactionYearlyAmountSuccessResponse) *pb.ApiResponseTransactionYearAmountSuccess
	ToProtoResponseMonthAmountFailed(status string, message string, row []*response.TransactionMonthlyAmountFailedResponse) *pb.ApiResponseTransactionMonthAmountFailed
	ToProtoResponseYearAmountFailed(status string, message string, row []*response.TransactionYearlyAmountFailedResponse) *pb.ApiResponseTransactionYearAmountFailed
	ToProtoResponseMonthMethod(status string, message string, row []*response.TransactionMonthlyMethodResponse) *pb.ApiResponseTransactionMonthPaymentMethod
	ToProtoResponseYearMethod(status string, message string, row []*response.TransactionYearlyMethodResponse) *pb.ApiResponseTransactionYearPaymentmethod

	ToProtoResponseTransaction(status string, message string, trans *response.TransactionResponse) *pb.ApiResponseTransaction
	ToProtoResponseTransactionDeleteAt(status string, message string, trans *response.TransactionResponseDeleteAt) *pb.ApiResponseTransactionDeleteAt
	ToProtoResponsesTransaction(status string, message string, transList []*response.TransactionResponse) *pb.ApiResponsesTransaction
	ToProtoResponseTransactionDelete(status string, message string) *pb.ApiResponseTransactionDelete
	ToProtoResponseTransactionAll(status string, message string) *pb.ApiResponseTransactionAll
	ToProtoResponsePaginationTransactionDeleteAt(pagination *pb.PaginationMeta, status string, message string, transactions []*response.TransactionResponseDeleteAt) *pb.ApiResponsePaginationTransactionDeleteAt
	ToProtoResponsePaginationTransaction(pagination *pb.PaginationMeta, status string, message string, transactions []*response.TransactionResponse) *pb.ApiResponsePaginationTransaction
}
