package database

import (
	"database/sql"
	"fmt"
	"strings"

	"MyCleanArchitecture/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	if order == nil {
		return fmt.Errorf("order cannot be nil")
	}

	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price, created_at) VALUES (?, ?, ?, ?, NOW())")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return fmt.Errorf("order with ID %s already exists: %w", order.ID, err)
		}
		return fmt.Errorf("failed to execute insert statement: %w", err)
	}
	return nil
}

func (r *OrderRepository) FindAll() ([]*entity.Order, error) {
	rows, err := r.Db.Query("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row into order struct: %w", err)
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("SELECT count(*) FROM orders").Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total count: %w", err)
	}
	return total, nil
}
