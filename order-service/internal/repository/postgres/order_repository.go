package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/asliddin03/go-grpc-order-service/order-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		pool: pool,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO orders(user_id, status, total_price, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id
	`

	err = tx.QueryRow(ctx, query,
		order.UserID,
		string(order.Status),
		order.TotalPrice,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID)
	if err != nil {
		return fmt.Errorf("insert into orders: %w", err)
	}

	itemQuery := `INSERT INTO order_items(order_id, product_id, quantity, price)
				  VALUES ($1, $2, $3, $4)
	`

	for _, item := range order.Items {
		_, err = tx.Exec(ctx, itemQuery,
			order.ID,
			item.ProductID,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			return fmt.Errorf("insert into order_items (order_id=%d, product_id=%d): %w", order.ID, item.ProductID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) GetByID(ctx context.Context, orderID int64) (*domain.Order, error) {
	orderQuery := `SELECT id, user_id, status, total_price, created_at, updated_at 
				   FROM orders WHERE id = $1
	`

	var order domain.Order
	var status string

	err := r.pool.QueryRow(ctx, orderQuery, orderID).Scan(
		&order.ID,
		&order.UserID,
		&status,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	order.Status = domain.OrderStatus(status)

	items, err := r.getItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return &order, nil
}

func (r *OrderRepository) ListByUserID(ctx context.Context, userID int64) ([]domain.Order, error) {
	query := `SELECT id, user_id, status, total_price, created_at, updated_at
			  FROM orders WHERE user_id = $1
			  ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)

	for rows.Next() {
		var order domain.Order
		var status string

		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&status,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		order.Status = domain.OrderStatus(status)

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// TODO: оптимизировать загрузку items, чтобы избежать проблемы N+1 запросов
	for i := range orders {
		items, err := r.getItemsByOrderID(ctx, orders[i].ID)
		if err != nil {
			return nil, err
		}

		orders[i].Items = items
	}

	return orders, nil
}

func (r *OrderRepository) getItemsByOrderID(ctx context.Context, orderID int64) ([]domain.OrderItem, error) {
	query := `SELECT product_id, quantity, price
			  FROM order_items WHERE order_id = $1
			  ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.OrderItem, 0)

	for rows.Next() {
		var item domain.OrderItem

		err = rows.Scan(
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
