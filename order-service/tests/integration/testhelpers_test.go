package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/asliddin03/go-grpc-order-service/order-service/internal/domain"
	postgresrepo "github.com/asliddin03/go-grpc-order-service/order-service/internal/repository/postgres"
)

func newTestPostgresPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("POSTGRES_DSN")
	require.NotEmpty(t, dsn, "POSTGRES_DSN must be set for integration tests")

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	err = pool.Ping(ctx)
	require.NoError(t, err)

	ensureSchema(t, ctx, pool)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func newTestOrderRepository(t *testing.T) (context.Context, *pgxpool.Pool, *postgresrepo.OrderRepository) {
	t.Helper()

	ctx := context.Background()
	pool := newTestPostgresPool(t, ctx)
	repository := postgresrepo.NewOrderRepository(pool)

	cleanupTables(t, ctx, pool)

	return ctx, pool, repository
}

func newTestOrder(
	userID int64,
	status domain.OrderStatus,
	totalPrice int64,
	items []domain.OrderItem,
) *domain.Order {
	now := time.Now().UTC()

	return &domain.Order{
		UserID:     userID,
		Status:     status,
		TotalPrice: totalPrice,
		CreatedAt:  now,
		UpdatedAt:  now,
		Items:      items,
	}
}

func ensureSchema(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	schema := `
		CREATE TABLE IF NOT EXISTS orders (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			status TEXT NOT NULL,
			total_price BIGINT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		);

		CREATE TABLE IF NOT EXISTS order_items (
			id BIGSERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			product_id BIGINT NOT NULL,
			quantity INTEGER NOT NULL,
			price BIGINT NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
	`

	_, err := pool.Exec(ctx, schema)
	require.NoError(t, err)
}

func cleanupTables(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE order_items, orders RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}
