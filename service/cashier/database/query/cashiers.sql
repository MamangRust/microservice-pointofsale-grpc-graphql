-- GetCashiers: Retrieves paginated list of active cashiers with search capability
-- Purpose: List all active cashiers for management UI
-- Parameters:
--   $1: search_term - Optional text to filter cashiers by name (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All cashier fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted cashiers (deleted_at IS NULL)
--   - Supports partial text matching on name field (case-insensitive)
--   - Returns newest cashiers first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetCashiers :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM cashiers
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- GetCashiersActive: Retrieves paginated list of active cashiers with search capability
-- Purpose: List all active cashiers for management UI
-- Parameters:
--   $1: search_term - Optional text to filter cashiers by name (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All cashier fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted cashiers (deleted_at IS NULL)
--   - Supports partial text matching on name field (case-insensitive)
--   - Returns newest cashiers first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetCashiersActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM cashiers
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- GetCashiersTrashed: Retrieves paginated list of soft-deleted cashiers with search capability
-- Purpose: List all trashed (soft-deleted) cashiers for recovery or audit purposes
-- Parameters:
--   $1: search_term - Optional text to filter cashiers by name (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All cashier fields plus total_count of matching records
-- Business Logic:
--   - Includes only soft-deleted cashiers (deleted_at IS NOT NULL)
--   - Supports partial text matching on name field (case-insensitive)
--   - Returns newest deleted cashiers first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetCashiersTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM cashiers
WHERE deleted_at IS NOT NULL
  AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;



-- GetCashiersByMerchant: Retrieves active cashiers filtered by merchant_id
-- Parameters:
--   $1: Merchant ID (required)
--   $2: Search term for cashier name (optional)
--   $3: Limit
--   $4: Offset
-- Returns:
--   Cashier records belonging to specified merchant with total_count
-- name: GetCashiersByMerchant :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM cashiers
WHERE merchant_id = $1
  AND deleted_at IS NULL
  AND ($2::TEXT IS NULL OR name ILIKE '%' || $2 || '%')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;



-- GetMonthlyTotalSalesCashier: Retrieves monthly sales totals for cashiers across two date ranges
-- Purpose: Compare sales performance between two time periods (typically current vs previous period)
-- Parameters:
--   $1: Start date of first period
--   $2: End date of first period
--   $3: Start date of second period
--   $4: End date of second period
-- Returns:
--   year: The year of the sales data
--   month_name: The full month name (e.g., "January")
--   total_sales: Sum of order totals for that month (0 if no sales)
-- Business Logic:
--   - Compares sales between two customizable time windows (e.g. this month vs last month)
--   - Ensures all months appear in results even with no sales (gap filling)
--   - Only includes active/non-deleted orders and cashiers
--   - Formats output for easy display in reports/dashboards
-- name: GetMonthlyTotalSalesCashier :many
WITH monthly_totals AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::TEXT AS year,
        EXTRACT(MONTH FROM o.created_at)::integer AS month,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND (
            (o.created_at >= $1 AND o.created_at <= $2)  
            OR (o.created_at >= $3 AND o.created_at <= $4)  
        )
    GROUP BY
        EXTRACT(YEAR FROM o.created_at),
        EXTRACT(MONTH FROM o.created_at)
),
all_months AS (
    SELECT 
        EXTRACT(YEAR FROM $1)::TEXT AS year,
        EXTRACT(MONTH FROM $1)::integer AS month,
        TO_CHAR($1, 'FMMonth') AS month_name
    
    UNION
    
    SELECT 
        EXTRACT(YEAR FROM $3)::TEXT AS year,
        EXTRACT(MONTH FROM $3)::integer AS month,
        TO_CHAR($3, 'FMMonth') AS month_name
)
SELECT 
    COALESCE(am.year, EXTRACT(YEAR FROM $1)::TEXT) AS year,
    COALESCE(am.month_name, TO_CHAR($1, 'FMMonth')) AS month,
    COALESCE(mt.total_sales, 0) AS total_sales
FROM 
    all_months am
LEFT JOIN 
    monthly_totals mt ON am.year = mt.year AND am.month = mt.month
ORDER BY 
    am.year::INT DESC,
    am.month DESC;

-- GetYearlyTotalSalesCashier: Retrieves yearly sales totals for cashiers across current and previous year
-- Purpose: Year-over-year sales comparison
-- Parameters:
--   $1: The current year (integer)
-- Returns:
--   year: The year as text
--   total_sales: Sum of order totals for that year (0 if no sales)
-- Business Logic:
--   - Automatically compares current year with previous year
--   - Includes zero-value years for complete reporting
--   - Filters by merchant while maintaining data integrity
-- name: GetYearlyTotalSalesCashier :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::integer AS year,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND (
            EXTRACT(YEAR FROM o.created_at) = $1::integer
            OR EXTRACT(YEAR FROM o.created_at) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM o.created_at)
),
all_years AS (
    SELECT $1::integer AS year
    UNION
    SELECT $1::integer - 1 AS year
)
SELECT 
    a.year::text AS year,
    COALESCE(yd.total_sales, 0) AS total_sales
FROM 
    all_years a
LEFT JOIN 
    yearly_data yd ON a.year = yd.year
ORDER BY 
    a.year DESC;


-- GetMonthlyTotalSalesByMerchant: Retrieves monthly sales totals filtered by merchant ID
-- Purpose: Provides monthly sales analytics for a specific merchant across two time periods
-- Parameters:
--   $1: Start date of first comparison period
--   $2: End date of first comparison period  
--   $3: Start date of second comparison period
--   $4: End date of second comparison period
--   $5: Merchant ID to filter by
-- Returns:
--   year: Year of sales data (text format)
--   month_name: Full month name (e.g. "January")
--   total_sales: Sum of order totals for that month (0 if no sales)
-- Business Logic:
--   - Compares sales between two customizable time windows (e.g. this month vs last month)
--   - Ensures all months appear in results even with no sales (gap filling)
--   - Only includes active/non-deleted orders and cashiers
--   - Formats output for easy display in reports/dashboards
-- name: GetMonthlyTotalSalesByMerchant :many
WITH monthly_totals AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::TEXT AS year,
        EXTRACT(MONTH FROM o.created_at)::integer AS month,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND (
            (o.created_at >= $1 AND o.created_at <= $2)  
            OR (o.created_at >= $3 AND o.created_at <= $4)  
        )
        AND o.merchant_id = $5
    GROUP BY
        EXTRACT(YEAR FROM o.created_at),
        EXTRACT(MONTH FROM o.created_at)
),
all_months AS (
    SELECT 
        EXTRACT(YEAR FROM $1)::TEXT AS year,
        EXTRACT(MONTH FROM $1)::integer AS month,
        TO_CHAR($1, 'FMMonth') AS month_name
    
    UNION
    
    SELECT 
        EXTRACT(YEAR FROM $3)::TEXT AS year,
        EXTRACT(MONTH FROM $3)::integer AS month,
        TO_CHAR($3, 'FMMonth') AS month_name
)
SELECT 
    COALESCE(am.year, EXTRACT(YEAR FROM $1)::TEXT) AS year,
    COALESCE(am.month_name, TO_CHAR($1, 'FMMonth')) AS month,
    COALESCE(mt.total_sales, 0) AS total_sales
FROM 
    all_months am
LEFT JOIN 
    monthly_totals mt ON am.year = mt.year AND am.month = mt.month
ORDER BY 
    am.year::INT DESC,
    am.month DESC;



-- GetYearlyTotalSalesByMerchant: Retrieves yearly sales totals filtered by merchant ID  
-- Purpose: Provides year-over-year sales comparison for a specific merchant
-- Parameters:
--   $1: Current year to analyze (integer)
--   $2: Merchant ID to filter by
-- Returns: 
--   year: Year as text
--   total_sales: Annual sales total (0 if no sales)
-- Business Logic:
--   - Automatically compares current year with previous year
--   - Includes zero-value years for complete reporting
--   - Filters by merchant while maintaining data integrity
-- name: GetYearlyTotalSalesByMerchant :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::integer AS year,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND (
            EXTRACT(YEAR FROM o.created_at) = $1::integer
            OR EXTRACT(YEAR FROM o.created_at) = $1::integer - 1
        )
        AND o.merchant_id = $2
    GROUP BY
        EXTRACT(YEAR FROM o.created_at)
),
all_years AS (
    SELECT $1::integer AS year
    UNION
    SELECT $1::integer - 1 AS year
)
SELECT 
    a.year::text AS year,
    COALESCE(yd.total_sales, 0) AS total_sales
FROM 
    all_years a
LEFT JOIN 
    yearly_data yd ON a.year = yd.year
ORDER BY 
    a.year DESC;



-- GetMonthlyTotalSalesById: Retrieves monthly sales totals filtered by cashier ID
-- Purpose: Provides monthly sales analytics for a specific merchant across two time periods
-- Parameters:
--   $1: Start date of first comparison period
--   $2: End date of first comparison period  
--   $3: Start date of second comparison period
--   $4: End date of second comparison period
--   $5: Cashier ID to filter by
-- Returns:
--   year: Year of sales data (text format)
--   month_name: Full month name (e.g. "January")
--   total_sales: Sum of order totals for that month (0 if no sales)
-- Business Logic:
--   - Compares sales between two customizable time windows (e.g. this month vs last month)
--   - Ensures all months appear in results even with no sales (gap filling)
--   - Only includes active/non-deleted orders and cashiers
--   - Formats output for easy display in reports/dashboards
-- name: GetMonthlyTotalSalesById :many
WITH monthly_totals AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::TEXT AS year,
        EXTRACT(MONTH FROM o.created_at)::integer AS month,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND (
            (o.created_at >= $1 AND o.created_at <= $2)  
            OR (o.created_at >= $3 AND o.created_at <= $4)  
        )
        AND c.cashier_id = $5
    GROUP BY
        EXTRACT(YEAR FROM o.created_at),
        EXTRACT(MONTH FROM o.created_at)
),
all_months AS (
    SELECT 
        EXTRACT(YEAR FROM $1)::TEXT AS year,
        EXTRACT(MONTH FROM $1)::integer AS month,
        TO_CHAR($1, 'FMMonth') AS month_name
    
    UNION
    
    SELECT 
        EXTRACT(YEAR FROM $3)::TEXT AS year,
        EXTRACT(MONTH FROM $3)::integer AS month,
        TO_CHAR($3, 'FMMonth') AS month_name
)
SELECT 
    COALESCE(am.year, EXTRACT(YEAR FROM $1)::TEXT) AS year,
    COALESCE(am.month_name, TO_CHAR($1, 'FMMonth')) AS month,
    COALESCE(mt.total_sales, 0) AS total_sales
FROM 
    all_months am
LEFT JOIN 
    monthly_totals mt ON am.year = mt.year AND am.month = mt.month
ORDER BY 
    am.year::INT DESC,
    am.month DESC;


-- GetYearlyTotalSalesById: Retrieves yearly sales totals filtered by cashier ID  
-- Purpose: Provides year-over-year sales comparison for a specific cashier
-- Parameters:
--   $1: Current year to analyze (integer)
--   $2: Merchant ID to filter by
-- Returns: 
--   year: Year as text
--   total_sales: Annual sales total (0 if no sales)
-- Business Logic:
--   - Automatically compares current year with previous year
--   - Includes zero-value years for complete reporting
--   - Filters by cashier while maintaining data integrity
-- name: GetYearlyTotalSalesById :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::integer AS year,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND (
            EXTRACT(YEAR FROM o.created_at) = $1::integer
            OR EXTRACT(YEAR FROM o.created_at) = $1::integer - 1
        )
        AND c.cashier_id = $2
    GROUP BY
        EXTRACT(YEAR FROM o.created_at)
),
all_years AS (
    SELECT $1::integer AS year
    UNION
    SELECT $1::integer - 1 AS year
)
SELECT 
    a.year::text AS year,
    COALESCE(yd.total_sales, 0) AS total_sales
FROM 
    all_years a
LEFT JOIN 
    yearly_data yd ON a.year = yd.year
ORDER BY 
    a.year DESC;



-- GetMonthlyCashier: Retrieves monthly sales activity for all cashiers within a 1-year period
-- Purpose: Provides cashier performance metrics by month for operational analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis period
-- Returns:
--   cashier_id: Unique identifier for the cashier
--   cashier_name: Full name of the cashier
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   order_count: Number of orders processed
--   total_sales: Gross sales amount generated
-- Business Logic:
--   - Analyzes a rolling 12-month period from the reference date
--   - Excludes deleted records to maintain data integrity
--   - Groups results by cashier and month for granular performance tracking
--   - Uses abbreviated month names for compact visual reporting
--   - Orders chronologically for trend analysis
-- name: GetMonthlyCashier :many
WITH date_range AS (
    SELECT 
        date_trunc('month', $1::timestamp) AS start_date,
        date_trunc('month', $1::timestamp) + interval '1 year' - interval '1 day' AS end_date
),
cashier_activity AS (
    SELECT
        c.cashier_id,
        c.name AS cashier_name,
        date_trunc('month', o.created_at) AS activity_month,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price)::NUMERIC AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND o.created_at BETWEEN (SELECT start_date FROM date_range) 
                             AND (SELECT end_date FROM date_range)
    GROUP BY
        c.cashier_id, c.name, activity_month
)
SELECT
    ca.cashier_id,
    ca.cashier_name,
    TO_CHAR(ca.activity_month, 'Mon') AS month,
    ca.order_count,
    ca.total_sales
FROM
    cashier_activity ca
ORDER BY
    ca.activity_month, ca.cashier_id;


-- GetYearlyCashier: Retrieves annual sales performance for cashiers over 5-year span
-- Purpose: Enables long-term cashier productivity trend analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 5-year analysis window
-- Returns:
--   year: 4-digit year as text
--   cashier_id: Unique cashier identifier  
--   cashier_name: Full name of the cashier
--   order_count: Annual transaction volume
--   total_sales: Yearly revenue generated
-- Business Logic:
--   - Covers current year plus previous 4 years (5-year total window)
--   - Maintains data quality by excluding soft-deleted records
--   - Provides both quantitative (order count) and financial (sales) metrics
--   - Orders results chronologically then by cashier for consistent reporting
--   - Designed for workforce planning and incentive calculations
-- name: GetYearlyCashier :many
WITH last_five_years AS (
    SELECT
        c.cashier_id,
        c.name AS cashier_name,
        EXTRACT(YEAR FROM o.created_at)::text AS year,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price) AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND EXTRACT(YEAR FROM o.created_at) BETWEEN (EXTRACT(YEAR FROM $1::timestamp) - 4) AND EXTRACT(YEAR FROM $1::timestamp)
    GROUP BY
        c.cashier_id, c.name, EXTRACT(YEAR FROM o.created_at)
)
SELECT
    year,
    cashier_id,
    cashier_name,
    order_count,
    total_sales
FROM
    last_five_years
ORDER BY
    year, cashier_id;



-- GetMonthlyCashierByCashierId: Retrieves monthly sales activity for all cashiers within a 1-year period by cashier id
-- Purpose: Provides cashier performance metrics by month for operational analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis period
--   $2: Reference cashier_id
-- Returns:
--   cashier_id: Unique identifier for the cashier
--   cashier_name: Full name of the cashier
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   order_count: Number of orders processed
--   total_sales: Gross sales amount generated
-- Business Logic:
--   - Analyzes a rolling 12-month period from the reference date
--   - Excludes deleted records to maintain data integrity
--   - Groups results by cashier and month for granular performance tracking
--   - Uses abbreviated month names for compact visual reporting
--   - Orders chronologically for trend analysis
-- name: GetMonthlyCashierByCashierId :many
WITH date_range AS (
    SELECT 
        date_trunc('month', $1::timestamp) AS start_date,
        date_trunc('month', $1::timestamp) + interval '1 year' - interval '1 day' AS end_date
),
cashier_activity AS (
    SELECT
        c.cashier_id,
        c.name AS cashier_name,
        date_trunc('month', o.created_at) AS activity_month,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price) AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND c.cashier_id = $2
        AND o.created_at BETWEEN (SELECT start_date FROM date_range) 
                             AND (SELECT end_date FROM date_range)
    GROUP BY
        c.cashier_id, c.name, activity_month
)
SELECT
    ca.cashier_id,
    ca.cashier_name,
    TO_CHAR(ca.activity_month, 'Mon') AS month,
    ca.order_count,
    ca.total_sales
FROM
    cashier_activity ca
ORDER BY
    ca.activity_month, ca.cashier_id;



-- GetYearlyCashierByCashierId: Retrieves annual sales performance for cashiers over 5-year span by cashier id
-- Purpose: Enables long-term cashier productivity trend analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 5-year analysis window
--   $2: Reference cashier_id
-- Returns:
--   year: 4-digit year as text
--   cashier_id: Unique cashier identifier  
--   cashier_name: Full name of the cashier
--   order_count: Annual transaction volume
--   total_sales: Yearly revenue generated
-- Business Logic:
--   - Covers current year plus previous 4 years (5-year total window)
--   - Maintains data quality by excluding soft-deleted records
--   - Provides both quantitative (order count) and financial (sales) metrics
--   - Orders results chronologically then by cashier for consistent reporting
--   - Designed for workforce planning and incentive calculations
-- name: GetYearlyCashierByCashierId :many
WITH last_five_years AS (
    SELECT
        c.cashier_id,
        c.name AS cashier_name,
        EXTRACT(YEAR FROM o.created_at)::text AS year,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price) AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND c.cashier_id = $2
        AND EXTRACT(YEAR FROM o.created_at) BETWEEN (EXTRACT(YEAR FROM $1::timestamp) - 4) AND EXTRACT(YEAR FROM $1::timestamp)
    GROUP BY
        c.cashier_id, c.name, EXTRACT(YEAR FROM o.created_at)
)
SELECT
    year,
    cashier_id,
    cashier_name,
    order_count,
    total_sales
FROM
    last_five_years
ORDER BY
    year, cashier_id;



-- GetMonthlyCashierByMerchant: Retrieves monthly sales activity for all cashiers within a 1-year period by merchant id
-- Purpose: Provides cashier performance metrics by month for operational analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis period
--   $2: Reference merchant_id
-- Returns:
--   cashier_id: Unique identifier for the cashier
--   cashier_name: Full name of the cashier
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   order_count: Number of orders processed
--   total_sales: Gross sales amount generated
-- Business Logic:
--   - Analyzes a rolling 12-month period from the reference date
--   - Excludes deleted records to maintain data integrity
--   - Groups results by cashier and month for granular performance tracking
--   - Uses abbreviated month names for compact visual reporting
--   - Orders chronologically for trend analysis
-- name: GetMonthlyCashierByMerchant :many
WITH date_range AS (
    SELECT 
        date_trunc('month', $1::timestamp) AS start_date,
        date_trunc('month', $1::timestamp) + interval '1 year' - interval '1 day' AS end_date
),
cashier_activity AS (
    SELECT
        c.cashier_id,
        c.name AS cashier_name,
        date_trunc('month', o.created_at) AS activity_month,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price) AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND c.merchant_id = $2
        AND o.created_at BETWEEN (SELECT start_date FROM date_range) 
                             AND (SELECT end_date FROM date_range)
    GROUP BY
        c.cashier_id, c.name, activity_month
)
SELECT
    ca.cashier_id,
    ca.cashier_name,
    TO_CHAR(ca.activity_month, 'Mon') AS month,
    ca.order_count,
    ca.total_sales
FROM
    cashier_activity ca
ORDER BY
    ca.activity_month, ca.cashier_id;



-- GetYearlyCashierByMerchant: Retrieves annual sales performance for cashiers over 5-year span by merchant id
-- Purpose: Enables long-term cashier productivity trend analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 5-year analysis window
--   $2: Reference cashier_id
-- Returns:
--   year: 4-digit year as text
--   cashier_id: Unique cashier identifier  
--   cashier_name: Full name of the cashier
--   order_count: Annual transaction volume
--   total_sales: Yearly revenue generated
-- Business Logic:
--   - Covers current year plus previous 4 years (5-year total window)
--   - Maintains data quality by excluding soft-deleted records
--   - Provides both quantitative (order count) and financial (sales) metrics
--   - Orders results chronologically then by cashier for consistent reporting
--   - Designed for workforce planning and incentive calculations
-- name: GetYearlyCashierByMerchant :many
WITH last_five_years AS (
    SELECT
        c.cashier_id,
        c.name AS cashier_name,
        EXTRACT(YEAR FROM o.created_at)::text AS year,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price) AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND c.merchant_id = $2
        AND EXTRACT(YEAR FROM o.created_at) BETWEEN (EXTRACT(YEAR FROM $1::timestamp) - 4) AND EXTRACT(YEAR FROM $1::timestamp)
    GROUP BY
        c.cashier_id, c.name, EXTRACT(YEAR FROM o.created_at)
)
SELECT
    year,
    cashier_id,
    cashier_name,
    order_count,
    total_sales
FROM
    last_five_years
ORDER BY
    year, cashier_id;


-- CreateCashier: Creates a new cashier record
-- Purpose: Add new cashier to the system
-- Parameters:
--   $1: merchant_id - Associated merchant ID
--   $2: user_id - User account ID for the cashier
--   $3: name - Full name of the cashier
-- Returns: Complete created cashier record
-- Business Logic:
--   - Sets created_at timestamp automatically
--   - Requires all mandatory fields
-- name: CreateCashier :one
INSERT INTO cashiers (merchant_id, user_id, name)
VALUES ($1, $2, $3) RETURNING *;


-- GetCashierByID: Retrieves active cashier by ID
-- Purpose: Fetch cashier details for display/editing
-- Parameters:
--   $1: cashier_id - ID of the cashier to retrieve
-- Returns: Full cashier record if found and active
-- Business Logic:
--   - Excludes soft-deleted records
--   - Returns single record or nothing
-- name: GetCashierById :one
SELECT *
FROM cashiers
WHERE cashier_id = $1
  AND deleted_at IS NULL;


-- UpdateCashier: Modifies cashier information
-- Purpose: Update cashier details
-- Parameters:
--   $1: cashier_id - ID of cashier to update
--   $2: name - New name value
-- Returns: Updated cashier record
-- Business Logic:
--   - Automatically updates updated_at timestamp
--   - Only affects active (non-deleted) records
--   - Returns the modified record for confirmation
-- name: UpdateCashier :one
UPDATE cashiers
SET name = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE cashier_id = $1
  AND deleted_at IS NULL
  RETURNING *;

-- TrashCashier: Soft-deletes a cashier record
-- Purpose: Remove cashier from active use without permanent deletion
-- Parameters:
--   $1: cashier_id - ID of cashier to deactivate
-- Returns: The soft-deleted cashier record
-- Business Logic:
--   - Sets deleted_at timestamp to current time
--   - Only works on currently active records
--   - Allows for recovery via restore function
-- name: TrashCashier :one
UPDATE cashiers
SET
    deleted_at = current_timestamp
WHERE
    cashier_id = $1
    AND deleted_at IS NULL
    RETURNING *;

-- RestoreCashier: Recovers a soft-deleted cashier
-- Purpose: Reactivate a previously trashed cashier
-- Parameters:
--   $1: cashier_id - ID of cashier to restore
-- Returns: The restored cashier record
-- Business Logic:
--   - Nullifies the deleted_at field
--   - Only works on previously deleted records
--   - Maintains all original data
-- name: RestoreCashier :one
UPDATE cashiers
SET
    deleted_at = NULL
WHERE
    cashier_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *;

-- DeleteCashierPermanently: Hard-deletes a cashier
-- Purpose: Completely remove cashier from database
-- Parameters:
--   $1: cashier_id - ID of cashier to delete
-- Business Logic:
--   - Permanent deletion of already soft-deleted records
--   - No return value (exec-only)
--   - Use with caution - irreversible operation
-- name: DeleteCashierPermanently :exec
DELETE FROM cashiers WHERE cashier_id = $1 AND deleted_at IS NOT NULL;

-- RestoreAllCashiers: Mass restoration of deleted cashiers
-- Purpose: Recover all trashed cashiers at once
-- Business Logic:
--   - Reactivates all soft-deleted cashiers
--   - No parameters needed
--   - Useful for system recovery scenarios
-- name: RestoreAllCashiers :exec
UPDATE cashiers
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentCashiers: Purges all trashed cashiers
-- Purpose: Clean up all soft-deleted records
-- Business Logic:
--   - Irreversible bulk deletion
--   - Only affects already soft-deleted records
--   - Typically used during database maintenance
-- name: DeleteAllPermanentCashiers :exec
DELETE FROM cashiers
WHERE
    deleted_at IS NOT NULL;
