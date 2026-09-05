package protomapper

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	"github.com/MamangRust/microservice-point-of-sale-shared/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type categoryProtoMapper struct {
}

func NewCategoryProtoMapper() *categoryProtoMapper {
	return &categoryProtoMapper{}
}

func (c *categoryProtoMapper) ToProtoResponseCategory(status string, message string, pbResponse *response.CategoryResponse) *pb.ApiResponseCategory {
	return &pb.ApiResponseCategory{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategory(pbResponse),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryDeleteAt(status string, message string, pbResponse *response.CategoryResponseDeleteAt) *pb.ApiResponseCategoryDeleteAt {
	return &pb.ApiResponseCategoryDeleteAt{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategoryDeleteAt(pbResponse),
	}
}

func (c *categoryProtoMapper) ToProtoResponsesCategory(status string, message string, pbResponse []*response.CategoryResponse) *pb.ApiResponsesCategory {
	return &pb.ApiResponsesCategory{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCategory(pbResponse),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryDelete(status string, message string) *pb.ApiResponseCategoryDelete {
	return &pb.ApiResponseCategoryDelete{
		Status:  status,
		Message: message,
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryAll(status string, message string) *pb.ApiResponseCategoryAll {
	return &pb.ApiResponseCategoryAll{
		Status:  status,
		Message: message,
	}
}

func (c *categoryProtoMapper) ToProtoResponsePaginationCategoryDeleteAt(pagination *pb.PaginationMeta, status string, message string, categories []*response.CategoryResponseDeleteAt) *pb.ApiResponsePaginationCategoryDeleteAt {
	return &pb.ApiResponsePaginationCategoryDeleteAt{
		Status:     status,
		Message:    message,
		Data:       c.mapResponsesCategoryDeleteAt(categories),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (c *categoryProtoMapper) ToProtoResponsePaginationCategory(pagination *pb.PaginationMeta, status string, message string, categories []*response.CategoryResponse) *pb.ApiResponsePaginationCategory {
	return &pb.ApiResponsePaginationCategory{
		Status:     status,
		Message:    message,
		Data:       c.mapResponsesCategory(categories),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryMonthlyPrice(status string, message string, row []*response.CategoryMonthPriceResponse) *pb.ApiResponseCategoryMonthPrice {
	return &pb.ApiResponseCategoryMonthPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCategoryMonthlyPrices(row),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryYearlyPrice(status string, message string, row []*response.CategoryYearPriceResponse) *pb.ApiResponseCategoryYearPrice {
	return &pb.ApiResponseCategoryYearPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCategoryYearlyPrices(row),
	}
}

func (c *categoryProtoMapper) ToProtoResponseMonthlyTotalPrice(status string, message string, row []*response.CategoriesMonthlyTotalPriceResponse) *pb.ApiResponseCategoryMonthlyTotalPrice {
	return &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategoryMonthlyTotalPrices(row),
	}
}

func (c *categoryProtoMapper) ToProtoResponseYearlyTotalPrice(status string, message string, row []*response.CategoriesYearlyTotalPriceResponse) *pb.ApiResponseCategoryYearlyTotalPrice {
	return &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategoryYearlyTotalPrices(row),
	}
}

func (c *categoryProtoMapper) mapResponseCategory(category *response.CategoryResponse) *pb.CategoryResponse {
	return &pb.CategoryResponse{
		Id:            int32(category.ID),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	}
}

func (c *categoryProtoMapper) mapResponsesCategory(categories []*response.CategoryResponse) []*pb.CategoryResponse {
	var mappedCategories []*pb.CategoryResponse

	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.mapResponseCategory(category))
	}

	return mappedCategories
}

func (c *categoryProtoMapper) mapResponseCategoryDeleteAt(category *response.CategoryResponseDeleteAt) *pb.CategoryResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt != nil {
		deletedAt = wrapperspb.String(*category.DeletedAt)
	}

	return &pb.CategoryResponseDeleteAt{
		Id:            int32(category.ID),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}

func (c *categoryProtoMapper) mapResponsesCategoryDeleteAt(categories []*response.CategoryResponseDeleteAt) []*pb.CategoryResponseDeleteAt {
	var mappedCategories []*pb.CategoryResponseDeleteAt

	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.mapResponseCategoryDeleteAt(category))
	}

	return mappedCategories
}

func (s *categoryProtoMapper) mapResponseCategoryMonthlyPrice(category *response.CategoryMonthPriceResponse) *pb.CategoryMonthPriceResponse {
	return &pb.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   int32(category.CategoryID),
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: int32(category.TotalRevenue),
	}
}

func (s *categoryProtoMapper) mapResponsesCategoryMonthlyPrices(c []*response.CategoryMonthPriceResponse) []*pb.CategoryMonthPriceResponse {
	var categoryRecords []*pb.CategoryMonthPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseCategoryMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryProtoMapper) mapResponseCategoryYearlyPrice(category *response.CategoryYearPriceResponse) *pb.CategoryYearPriceResponse {
	return &pb.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         int32(category.CategoryID),
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       int32(category.TotalRevenue),
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func (s *categoryProtoMapper) mapResponsesCategoryYearlyPrices(c []*response.CategoryYearPriceResponse) []*pb.CategoryYearPriceResponse {
	var categoryRecords []*pb.CategoryYearPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseCategoryYearlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryProtoMapper) mapResponseCashierMonthlyTotalPrice(c *response.CategoriesMonthlyTotalPriceResponse) *pb.CategoriesMonthlyTotalPriceResponse {
	return &pb.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int32(c.TotalRevenue),
	}
}

func (s *categoryProtoMapper) mapResponseCategoryMonthlyTotalPrices(c []*response.CategoriesMonthlyTotalPriceResponse) []*pb.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pb.CategoriesMonthlyTotalPriceResponse

	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, s.mapResponseCashierMonthlyTotalPrice(Category))
	}

	return CategoryRecords
}

func (s *categoryProtoMapper) mapResponseCategoryYearlyTotalSale(c *response.CategoriesYearlyTotalPriceResponse) *pb.CategoriesYearlyTotalPriceResponse {
	return &pb.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: int32(c.TotalRevenue),
	}
}

func (s *categoryProtoMapper) mapResponseCategoryYearlyTotalPrices(c []*response.CategoriesYearlyTotalPriceResponse) []*pb.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pb.CategoriesYearlyTotalPriceResponse

	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, s.mapResponseCategoryYearlyTotalSale(Category))
	}

	return CategoryRecords
}
