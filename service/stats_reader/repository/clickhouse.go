package repository

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseReaderRepository struct {
	conn clickhouse.Conn
}

func NewClickHouseReaderRepository(conn clickhouse.Conn) *ClickHouseReaderRepository {
	return &ClickHouseReaderRepository{conn: conn}
}

// --- Order stats ---

func (r *ClickHouseReaderRepository) GetMonthlyTotalRevenue(ctx context.Context, year, month int) ([]MonthlyRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       sum(total_price) AS total_revenue,
		       count() AS order_count
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY year, month, toMonth(event_time)
		ORDER BY year, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query monthly total revenue: %w", err)
	}
	defer rows.Close()

	var results []MonthlyRevenue
	for rows.Next() {
		var m MonthlyRevenue
		if err := rows.Scan(&m.Year, &m.Month, &m.TotalRevenue, &m.OrderCount); err != nil {
			return nil, fmt.Errorf("scan monthly total revenue: %w", err)
		}
		results = append(results, m)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetYearlyTotalRevenue(ctx context.Context, year int) ([]YearlyRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       sum(total_price) AS total_revenue,
		       count() AS order_count
		FROM order_daily FINAL
		WHERE toYear(event_time) IN (?, ?)
		GROUP BY year
		ORDER BY year DESC
	`
	rows, err := r.conn.Query(ctx, query, year, year-1)
	if err != nil {
		return nil, fmt.Errorf("query yearly total revenue: %w", err)
	}
	defer rows.Close()

	var results []YearlyRevenue
	for rows.Next() {
		var y YearlyRevenue
		if err := rows.Scan(&y.Year, &y.TotalRevenue, &y.OrderCount); err != nil {
			return nil, fmt.Errorf("scan yearly total revenue: %w", err)
		}
		results = append(results, y)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierMonthlyRevenue(ctx context.Context, cashierID int) ([]CashierMonthlyRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       cashier_id,
		       sum(total_price) AS total_revenue,
		       count() AS order_count
		FROM order_daily FINAL
		WHERE cashier_id = ?
		GROUP BY year, month, toMonth(event_time), cashier_id
		ORDER BY toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, cashierID)
	if err != nil {
		return nil, fmt.Errorf("query cashier monthly revenue: %w", err)
	}
	defer rows.Close()

	var results []CashierMonthlyRevenue
	for rows.Next() {
		var c CashierMonthlyRevenue
		if err := rows.Scan(&c.Year, &c.Month, &c.CashierID, &c.TotalRevenue, &c.OrderCount); err != nil {
			return nil, fmt.Errorf("scan cashier monthly revenue: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

// --- Product / category stats ---

func (r *ClickHouseReaderRepository) GetProductMonthlySold(ctx context.Context, year, month int) ([]ProductMonthlySold, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       product_id,
		       sum(quantity) AS quantity,
		       sum(subtotal) AS subtotal
		FROM order_item_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY month, toMonth(event_time), product_id
		ORDER BY toMonth(event_time), subtotal DESC
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query product monthly sold: %w", err)
	}
	defer rows.Close()

	var results []ProductMonthlySold
	for rows.Next() {
		var p ProductMonthlySold
		if err := rows.Scan(&p.Month, &p.ProductID, &p.Quantity, &p.Subtotal); err != nil {
			return nil, fmt.Errorf("scan product monthly sold: %w", err)
		}
		results = append(results, p)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryMonthlySold(ctx context.Context, year, month int) ([]CategoryMonthlySold, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       category_id,
		       sum(quantity) AS quantity,
		       sum(subtotal) AS subtotal
		FROM order_item_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY month, toMonth(event_time), category_id
		ORDER BY toMonth(event_time), subtotal DESC
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query category monthly sold: %w", err)
	}
	defer rows.Close()

	var results []CategoryMonthlySold
	for rows.Next() {
		var c CategoryMonthlySold
		if err := rows.Scan(&c.Month, &c.CategoryID, &c.Quantity, &c.Subtotal); err != nil {
			return nil, fmt.Errorf("scan category monthly sold: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

// --- Transaction stats ---

func (r *ClickHouseReaderRepository) GetTransactionMonthlySuccess(ctx context.Context, year, month int) ([]TransactionMonthlySuccess, error) {
	// POS stores both "success" (API creates) and "Completed"/"completed"
	// (legacy seeder data) for a successful payment; match case-insensitively so
	// the reader cross-checks with the OLTP payment_status values.
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       countIf(lower(status) IN ('success', 'completed')) AS total_count,
		       sumIf(amount, lower(status) IN ('success', 'completed')) AS total_amount
		FROM transaction_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY month, toMonth(event_time)
		ORDER BY toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query transaction monthly success: %w", err)
	}
	defer rows.Close()

	var results []TransactionMonthlySuccess
	for rows.Next() {
		var t TransactionMonthlySuccess
		if err := rows.Scan(&t.Month, &t.TotalCount, &t.TotalAmount); err != nil {
			return nil, fmt.Errorf("scan transaction monthly success: %w", err)
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

// --- Cashier stats ---

func (r *ClickHouseReaderRepository) GetCashierMonthlyOrders(ctx context.Context, cashierID int) ([]CashierMonthlyOrders, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       cashier_id,
		       count() AS order_count,
		       sum(total_price) AS total_amount
		FROM order_daily FINAL
		WHERE cashier_id = ?
		GROUP BY month, toMonth(event_time), cashier_id
		ORDER BY toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, cashierID)
	if err != nil {
		return nil, fmt.Errorf("query cashier monthly orders: %w", err)
	}
	defer rows.Close()

	var results []CashierMonthlyOrders
	for rows.Next() {
		var c CashierMonthlyOrders
		if err := rows.Scan(&c.Month, &c.CashierID, &c.OrderCount, &c.TotalAmount); err != nil {
			return nil, fmt.Errorf("scan cashier monthly orders: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}
