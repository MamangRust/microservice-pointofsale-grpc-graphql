package seeder

import (
	categorydb "github.com/MamangRust/microservice-point-of-sale-category/database/schema"
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"

	"go.uber.org/zap"
)

type categorySeeder struct {
	categorydb *categorydb.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCategorySeeder(categorydb *categorydb.Queries, ctx context.Context, logger logger.LoggerInterface) *categorySeeder {
	return &categorySeeder{
		categorydb: categorydb,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *categorySeeder) Seed() error {
	categoryNames := []string{
		"Electronics", "Clothing", "Groceries", "Toys", "Home & Kitchen",
		"Books", "Beauty & Health", "Sports & Outdoors", "Automotive", "Furniture",
	}

	categoryDescriptions := []string{
		"Best electronics products", "Latest fashion trends", "Fresh groceries",
		"Fun toys for kids", "Essentials for home & kitchen",
		"Books for all ages", "Beauty and health products",
		"Outdoor sports equipment", "Automotive accessories", "Stylish furniture",
	}

	for i := 0; i < 10; i++ {
		name := categoryNames[i%len(categoryNames)]
		description := categoryDescriptions[i%len(categoryDescriptions)]
		slugCategory := fmt.Sprintf("%s-%d", name, i+1)

		_, err := r.categorydb.CreateCategory(r.ctx, categorydb.CreateCategoryParams{
			Name:         name,
			Description:  ptrString(description),
			SlugCategory: ptrString(slugCategory),
		})
		if err != nil {
			r.logger.Error("Failed to create category:", zap.Error(err))
			return err
		}
		r.logger.Debug("Category created:", zap.String("slug", slugCategory))
	}

	r.logger.Info("Category seeding completed successfully.", zap.Int("count", 10))
	return nil
}
