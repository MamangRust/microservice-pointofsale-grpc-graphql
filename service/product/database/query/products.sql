-- GetProducts: Retrieves paginated list of active products with search capability
-- Purpose: List all active (non-deleted) products for display in UI
-- Parameters:
--   $1: search_term - Optional text to filter products (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All product fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted products (deleted_at IS NULL)
--   - Supports partial, case-insensitive search on name, description, brand, slug, and barcode
--   - Orders results by newest first (created_at DESC)
--   - Uses COUNT(*) OVER() to include total matching record count for pagination UI
-- name: GetProducts :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM products as p
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL 
       OR p.name ILIKE '%' || $1 || '%'
       OR p.description ILIKE '%' || $1 || '%'
       OR p.brand ILIKE '%' || $1 || '%'
       OR p.slug_product ILIKE '%' || $1 || '%'
       OR p.barcode ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- GetProductsActive: Retrieves paginated list of active products (duplicate of GetProducts)
-- Purpose: Explicitly return active (non-deleted) products with search capability
-- Parameters:
--   $1: search_term - Optional text to filter products (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All product fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted products (deleted_at IS NULL)
--   - Supports partial, case-insensitive search on name, description, brand, slug, and barcode
--   - Ordered by newest first (created_at DESC)
--   - Useful if frontend/backend wants clearer distinction in naming
-- name: GetProductsActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM products as p
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL 
       OR p.name ILIKE '%' || $1 || '%'
       OR p.description ILIKE '%' || $1 || '%'
       OR p.brand ILIKE '%' || $1 || '%'
       OR p.slug_product ILIKE '%' || $1 || '%'
       OR p.barcode ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- GetProductsTrashed: Retrieves paginated list of trashed (soft-deleted) products
-- Purpose: List deleted products for admin to manage recovery or audit
-- Parameters:
--   $1: search_term - Optional text to filter trashed products (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All product fields plus total_count of matching trashed records
-- Business Logic:
--   - Includes only soft-deleted products (deleted_at IS NOT NULL)
--   - Supports partial, case-insensitive search on name, description, brand, slug, and barcode
--   - Returns by newest first (created_at DESC)
--   - Used for "Trash Bin" UI or soft-delete management
-- name: GetProductsTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM products as p
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL 
       OR p.name ILIKE '%' || $1 || '%'
       OR p.description ILIKE '%' || $1 || '%'
       OR p.brand ILIKE '%' || $1 || '%'
       OR p.slug_product ILIKE '%' || $1 || '%'
       OR p.barcode ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;



-- GetProductsByMerchant: Retrieves paginated and filtered products owned by a specific merchant
-- Purpose: Allow merchants to view and manage their own products with advanced filtering options
-- Parameters:
--   $1: merchant_id - Filter products belonging to this merchant
--   $2: search_term - Optional text to filter by product name or description
--   $3: category_id - Optional category filter (0 or NULL to ignore)
--   $4: min_price - Minimum price filter (0 to ignore)
--   $5: max_price - Maximum price filter (0 to ignore, defaults to very high value)
--   $6: limit - Number of products to return (pagination)
--   $7: offset - Number of products to skip (pagination)
-- Returns:
--   - Filtered list of product fields including category name
--   - total_count of all matching products for pagination UI
-- Business Logic:
--   - Excludes soft-deleted products (deleted_at IS NULL)
--   - Supports case-insensitive partial search on name and description
--   - Filters by category ID only if provided
--   - Filters by price range only if values provided (>= min_price and <= max_price)
--   - Ordered by newest products first (created_at DESC)
-- name: GetProductsByMerchant :many
WITH filtered_products AS (
    SELECT 
        p.product_id,
        p.name,
        p.description,
        p.price,
        p.count_in_stock,
        p.brand,
        p.image_product,
        p.created_at,  
        c.name AS category_name
    FROM 
        products p
    JOIN 
        categories c ON p.category_id = c.category_id
    WHERE 
        p.deleted_at IS NULL
        AND p.merchant_id = $1  
        AND (
            p.name ILIKE '%' || COALESCE($2, '') || '%' 
            OR p.description ILIKE '%' || COALESCE($2, '') || '%'
            OR $2 IS NULL
        )
        AND (
            c.category_id = NULLIF($3, 0) 
            OR NULLIF($3, 0) IS NULL
        )
        AND (
            p.price >= COALESCE(NULLIF($4, 0), 0)
            AND p.price <= COALESCE(NULLIF($5, 0), 999999999)
        )
)
SELECT 
    (SELECT COUNT(*) FROM filtered_products) AS total_count,
    fp.*
FROM 
    filtered_products fp
ORDER BY 
    fp.created_at DESC
LIMIT $6 OFFSET $7;




-- GetProductsByCategoryName: Retrieves paginated and filtered products under a specific category name
-- Purpose: Display products by category for customers or category-focused pages
-- Parameters:
--   $1: category_name - The name of the category to filter by
--   $2: search_term - Optional text to filter by product name or description
--   $3: min_price - Minimum price filter (0 to ignore)
--   $4: max_price - Maximum price filter (0 to ignore, defaults to very high value)
--   $5: limit - Number of products to return (pagination)
--   $6: offset - Number of products to skip (pagination)
-- Returns:
--   - Filtered list of product fields including category name
--   - total_count of all matching products for pagination UI
-- Business Logic:
--   - Excludes soft-deleted products (deleted_at IS NULL)
--   - Matches category name exactly
--   - Supports case-insensitive partial search on name and description
--   - Filters by category ID only if provided
--   - Filters by price range only if values provided
--   - Ordered by newest products first (created_at DESC)
-- name: GetProductsByCategoryName :many
WITH filtered_products AS (
    SELECT 
        p.product_id,
        p.merchant_id,
        p.category_id,
        p.slug_product,
        p.weight,
        p.name,
        p.description,
        p.price,
        p.count_in_stock,
        p.brand,
        p.image_product,
        p.barcode,
        p.created_at,
        p.updated_at,  
        p.deleted_at,
        c.name AS category_name
    FROM 
        products p
    JOIN 
        categories c ON p.category_id = c.category_id
    WHERE 
        p.deleted_at IS NULL
        AND c.name = $1  
        AND (
            $2::text IS NULL 
            OR p.name ILIKE '%' || $2::text || '%' 
            OR p.description ILIKE '%' || $2::text || '%'
        )
        AND (
            ($3::int IS NULL OR p.price >= $3::int)
            AND ($4::int IS NULL OR p.price <= $4::int)
        )
)
SELECT 
    (SELECT COUNT(*) FROM filtered_products) AS total_count,
    fp.*
FROM 
    filtered_products fp
ORDER BY 
    fp.created_at DESC
LIMIT $5 OFFSET $6;



-- CreateProduct: Creates a new product record
-- Purpose: Add a new product to inventory
-- Parameters:
--   $1: merchant_id - Merchant who owns the product
--   $2: category_id - Product category
--   $3: name - Product name
--   $4: description - Detailed description
--   $5: price - Selling price
--   $6: count_in_stock - Inventory quantity
--   $7: brand - Manufacturer brand
--   $8: weight - Product weight
--   $9: slug_product - URL-friendly identifier
--   $10: image_product - Image URL/path
--   $11: barcode - Product barcode
-- Returns: Complete created product record
-- Business Logic:
--   - Sets created_at automatically
--   - Validates required fields
--   - Initializes inventory tracking
-- name: CreateProduct :one
INSERT INTO products (merchant_id, category_id, name, description, price, count_in_stock, brand, weight, slug_product, image_product, barcode)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- GetProductByID: Retrieves active product by ID
-- Purpose: Fetch product details for display/purchase
-- Parameters:
--   $1: product_id - ID of product to retrieve
-- Returns: Full product record if active
-- Business Logic:
--   - Excludes deleted products
--   - Used for product pages and checkout
-- name: GetProductByID :one
SELECT *
FROM products
WHERE product_id = $1
  AND deleted_at IS NULL;

-- GetProductByIdTrashed: Retrieves product including deleted
-- Purpose: View deleted products for restoration
-- Parameters:
--   $1: product_id - ID of product to retrieve
-- Returns: Product record regardless of deletion status
-- Business Logic:
--   - Bypasses deleted_at filter
--   - Used in admin/recovery interfaces
-- name: GetProductByIdTrashed :one
SELECT * FROM products WHERE product_id = $1;

-- UpdateProduct: Modifies product information
-- Purpose: Update product details
-- Parameters:
--   $1: product_id - Target product ID
--   $2: category_id - Updated category
--   $3: name - Updated name
--   $4: description - Updated description
--   $5: price - Updated price
--   $6: count_in_stock - Updated inventory count
--   $7: brand - Updated brand
--   $8: weight - Updated weight
--   $9: image_product - Updated image
--   $10: barcode - Updated barcode
-- Returns: Updated product record
-- Business Logic:
--   - Auto-updates updated_at
--   - Only modifies active products
--   - Validates all fields
-- name: UpdateProduct :one
UPDATE products
SET category_id = $2,
    name = $3,
    description = $4,
    price = $5,
    count_in_stock = $6,
    brand = $7,
    weight = $8,
    image_product = $9,
    barcode = $10,
    updated_at = CURRENT_TIMESTAMP
WHERE product_id = $1
  AND deleted_at IS NULL
  RETURNING *;

-- UpdateProductCountStock: Updates inventory count
-- Purpose: Adjust product stock levels
-- Parameters:
--   $1: product_id - Product to update
--   $2: count_in_stock - New inventory count
-- Returns: Updated product record
-- Business Logic:
--   - Dedicated stock adjustment function
--   - Used when inventory changes
--   - Validates non-negative quantity
-- name: UpdateProductCountStock :one
UPDATE products
SET count_in_stock = $2
WHERE product_id = $1
    AND deleted_at IS NULL
RETURNING *;

-- IncrementProductCountStock: Atomically increments product stock
-- Purpose: Restore inventory after a compensated order operation
-- Parameters:
--   $1: product_id - Product to update
--   $2: quantity - Amount to add to current stock
-- Returns: Updated product record
-- Business Logic:
--   - Atomic increment (no lost update between read and write)
--   - Refuses non-positive quantities
-- name: IncrementProductCountStock :one
UPDATE products
SET count_in_stock = count_in_stock + sqlc.arg(quantity),
    updated_at = CURRENT_TIMESTAMP
WHERE product_id = sqlc.arg(product_id)
    AND deleted_at IS NULL
    AND sqlc.arg(quantity) > 0
RETURNING *;

-- DecrementProductCountStock: Atomically decrements product stock
-- Purpose: Safe inventory decrement with optimistic concurrency guard
-- Parameters:
--   $1: product_id - Product to update
--   $2: quantity - Amount to subtract from current stock
-- Returns: Updated product record
-- Business Logic:
--   - Atomic decrement (no lost update between read and write)
--   - Refuses to go below zero: WHERE count_in_stock >= $2
--   - Returns zero rows when stock is insufficient (caller maps to 400)
-- name: DecrementProductCountStock :one
UPDATE products
SET count_in_stock = count_in_stock - sqlc.arg(quantity),
    updated_at = CURRENT_TIMESTAMP
WHERE product_id = sqlc.arg(product_id)
    AND deleted_at IS NULL
    AND count_in_stock >= sqlc.arg(quantity)
RETURNING *;

-- TrashProduct: Soft-deletes a product
-- Purpose: Remove product from active listings
-- Parameters:
--   $1: product_id - Product to deactivate
-- Returns: The soft-deleted product
-- Business Logic:
--   - Sets deleted_at timestamp
--   - Preserves product data
--   - Excludes from active queries
-- name: TrashProduct :one
UPDATE products
SET
    deleted_at = current_timestamp
WHERE
    product_id = $1
    AND deleted_at IS NULL
    RETURNING *;

-- RestoreProduct: Recovers a soft-deleted product
-- Purpose: Reactivate a removed product
-- Parameters:
--   $1: product_id - Product to restore
-- Returns: The restored product
-- Business Logic:
--   - Nullifies deleted_at
--   - Returns product to active status
-- name: RestoreProduct :one
UPDATE products
SET
    deleted_at = NULL
WHERE
    product_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *;

-- DeleteProductPermanently: Hard-deletes a product
-- Purpose: Completely remove product record
-- Parameters:
--   $1: product_id - Product to delete
-- Business Logic:
--   - Permanent deletion
--   - Only affects already trashed products
--   - Irreversible operation
-- name: DeleteProductPermanently :exec
DELETE FROM products WHERE product_id = $1 AND deleted_at IS NOT NULL;


-- RestoreAllProducts: Mass restoration of deleted products
-- Purpose: Reactivate all trashed products
-- Business Logic:
--   - Bulk restore operation
--   - Used during data recovery
-- name: RestoreAllProducts :exec
UPDATE products
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentProducts: Purges all trashed products
-- Purpose: Clean up deleted products
-- Business Logic:
--   - Bulk permanent deletion
--   - Database maintenance operation
-- name: DeleteAllPermanentProducts :exec
DELETE FROM products
WHERE
    deleted_at IS NOT NULL;