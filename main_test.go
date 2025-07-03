package main

import (
	"database/sql"
	"testing"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lucasHSantiago/go-select-benchmark/.gen/order/public/model"
	. "github.com/lucasHSantiago/go-select-benchmark/.gen/order/public/table"
)

var db *sql.DB

func setup(b *testing.B) {
	connStr := "user=postgres password=admin dbname=order sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		b.Fatalf("failed to open db: %v", err)
	}
}

func shutdown() {
	db.Close()
}

func BenchmarkJet(b *testing.B) {
	setup(b)
	defer shutdown()

	for range b.N {
		stmt := SELECT(
			Orders.AllColumns,
			OrderItems.AllColumns,
		).FROM(
			Orders.
				INNER_JOIN(OrderItems, Orders.ID.EQ(OrderItems.OrderID)),
		).ORDER_BY(Orders.ID.ASC())

		var dest []struct {
			model.Orders
			Itens []struct {
				model.OrderItems
			}
		}

		err := stmt.Query(db, &dest)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}

		if len(dest) != 50000 {
			b.Fatalf("expected 50000 results, got %d", len(dest))
		}

		if len(dest[0].Itens) != 5 {
			b.Fatalf("expected 5 itens, got %d", len(dest[0].Itens))
		}
	}
}

func BenchmarkSqlx(b *testing.B) {
	setup(b)
	defer shutdown()

	dbx := sqlx.NewDb(db, "postgres")

	type OrderWithItems struct {
		model.Orders
		Itens []model.OrderItems
	}

	for i := 0; i < b.N; i++ {
		var results []struct {
			ID           int32      `db:"orders.id"`
			CustomerName string     `db:"orders.customer_name"`
			CreatedAt    *time.Time `db:"orders.created_at"`
			OrderItemID  int32      `db:"order_items.id"`
			OrderID      *int32     `db:"order_items.order_id"`
			ProductName  string     `db:"order_items.product_name"`
			Price        float64    `db:"order_items.price"`
			Quantity     *int32     `db:"order_items.quantity"`
		}

		query := `
		SELECT orders.id AS "orders.id",
			orders.customer_name AS "orders.customer_name",
			orders.created_at AS "orders.created_at",
			order_items.id AS "order_items.id",
			order_items.order_id AS "order_items.order_id",
			order_items.product_name AS "order_items.product_name",
			order_items.price AS "order_items.price",
			order_items.quantity AS "order_items.quantity"
		FROM public.orders
			INNER JOIN public.order_items ON (orders.id = order_items.order_id)
		ORDER BY orders.id ASC;
		`
		err := dbx.Select(&results, query)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}

		var orders []OrderWithItems
		orderIdx := make(map[int32]int)
		for _, row := range results {
			idx, ok := orderIdx[row.ID]
			if !ok {
				orders = append(orders, OrderWithItems{
					Orders: model.Orders{
						ID:           row.ID,
						CustomerName: row.CustomerName,
						CreatedAt:    row.CreatedAt,
					},
				})
				idx = len(orders) - 1
				orderIdx[row.ID] = idx
			}
			orders[idx].Itens = append(orders[idx].Itens, model.OrderItems{
				ID:          row.OrderItemID,
				OrderID:     row.OrderID,
				ProductName: row.ProductName,
				Price:       row.Price,
				Quantity:    row.Quantity,
			})
		}

		if len(orders) != 50000 {
			b.Fatalf("expected 50000 results, got %d", len(orders))
		}

		if len(orders[0].Itens) != 5 {
			b.Fatalf("expected 5 itens, got %d", len(orders[0].Itens))
		}
	}
}
