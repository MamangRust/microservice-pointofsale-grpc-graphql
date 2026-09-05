package seeder

import (
	productdb "github.com/MamangRust/microservice-point-of-sale-product/database/schema"
	categorydb "github.com/MamangRust/microservice-point-of-sale-category/database/schema"
	merchantdb "github.com/MamangRust/microservice-point-of-sale-merchant/database/schema"
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type productSeeder struct {
	productdb *productdb.Queries
	categorydb *categorydb.Queries
	merchantdb *merchantdb.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewProductSeeder(productdb *productdb.Queries, categorydb *categorydb.Queries, merchantdb *merchantdb.Queries, ctx context.Context, logger logger.LoggerInterface) *productSeeder {
	return &productSeeder{
		productdb: productdb,
		categorydb: categorydb,
		merchantdb: merchantdb,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *productSeeder) Seed() error {
	merchants, err := r.merchantdb.GetMerchants(r.ctx, merchantdb.GetMerchantsParams{
		Column1: "",
		Limit:   20,
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to get merchants:", zap.Any("error", err))
		return err
	}

	categories, err := r.categorydb.GetCategories(r.ctx, categorydb.GetCategoriesParams{
		Column1: "",
		Limit:   20,
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to get categories:", zap.Any("error", err))
		return err
	}

	if len(merchants) == 0 || len(categories) == 0 {
		r.logger.Error("No merchants or categories found, skipping seeding")
		return nil
	}

	productNames := []string{
		"Smartphone", "Laptop", "Wireless Earbuds", "Gaming Mouse", "Mechanical Keyboard",
		"Smartwatch", "Power Bank", "Bluetooth Speaker", "External Hard Drive", "Monitor",
	}
	brands := []string{"Samsung", "Apple", "Sony", "Logitech", "Razer", "Xiaomi", "HP", "Dell", "Acer", "Asus"}
	images := []string{
		"image1.jpg", "image2.jpg", "image3.jpg", "image4.jpg", "image5.jpg",
		"image6.jpg", "image7.jpg", "image8.jpg", "image9.jpg", "image10.jpg",
	}

	for i := 0; i < 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		category := categories[rand.Intn(len(categories))]
		name := productNames[rand.Intn(len(productNames))]
		brand := brands[rand.Intn(len(brands))]
		price := int32(rand.Intn(5000000) + 50000)
		countInStock := int32(rand.Intn(100) + 1)
		weight := int32(rand.Intn(5000) + 100)
		slug := fmt.Sprintf("%s-%d", name, rand.Intn(1000))
		image := images[rand.Intn(len(images))]
		barcode := fmt.Sprintf("BC-%d", rand.Intn(9999999))

		_, err := r.productdb.CreateProduct(r.ctx, productdb.CreateProductParams{
			MerchantID:   merchant.MerchantID,
			CategoryID:   category.CategoryID,
			Name:         name,
			Description:  ptrString(fmt.Sprintf("Description for %s", name)),
			Price:        price,
			CountInStock: countInStock,
			Brand:        ptrString(brand),
			Weight:       ptrInt32(weight),
			SlugProduct:  ptrString(slug),
			ImageProduct: ptrString(image),
			Barcode:      ptrString(barcode),
		})

		if err != nil {
			r.logger.Error("Failed to create product:", zap.Error(err))
			return err
		}

	}

	r.logger.Info("Product seeding completed successfully.", zap.Int("count", 10))
	return nil
}
