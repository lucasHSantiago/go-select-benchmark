package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jackskj/carta"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lucasHSantiago/go-select-benchmark/.gen/order/public/model"
	. "github.com/lucasHSantiago/go-select-benchmark/.gen/order/public/table"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *sql.DB

func TestMain(m *testing.M) {
	const connStr = "user=postgres password=admin dbname=order sslmode=disable"
	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic("failed to open db: " + err.Error())
	}
	defer db.Close()

	db.SetMaxOpenConns(0)
	db.SetMaxIdleConns(0)

	os.Exit(m.Run())
}

func BenchmarkJet(b *testing.B) {
	stmt := SELECT(
		Orders.AllColumns,
		OrderItems.AllColumns,
	).FROM(
		Orders.
			INNER_JOIN(OrderItems, Orders.ID.EQ(OrderItems.OrderID)),
	).ORDER_BY(Orders.ID.ASC())

	type dest []struct {
		model.Orders
		Itens []struct {
			model.OrderItems
		}
	}

	for range b.N {
		dest := make(dest, 0, 50000)
		err := stmt.QueryContext(b.Context(), db, &dest)
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

func BenchmarkJetOneResult(b *testing.B) {
	stmt := SELECT(
		Orders.AllColumns,
		OrderItems.AllColumns,
	).FROM(
		Orders.
			INNER_JOIN(OrderItems, Orders.ID.EQ(OrderItems.OrderID)),
	).WHERE(
		Orders.ID.EQ(Int32(1)),
	).ORDER_BY(
		Orders.ID.ASC(),
	)

	type dest struct {
		model.Orders
		Itens []struct {
			model.OrderItems
		}
	}

	for range b.N {
		dest := dest{}
		err := stmt.QueryContext(b.Context(), db, &dest)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}

		if len(dest.Itens) != 5 {
			b.Fatalf("expected 5 itens, got %d", len(dest.Itens))
		}
	}
}

func BenchmarkSqlx(b *testing.B) {
	dbx := sqlx.NewDb(db, "postgres")

	type OrderWithItems struct {
		model.Orders
		Itens []model.OrderItems
	}

	type results []struct {
		ID           int32      `db:"orders.id"`
		CustomerName string     `db:"orders.customer_name"`
		CreatedAt    *time.Time `db:"orders.created_at"`
		OrderItemID  int32      `db:"order_items.id"`
		OrderID      *int32     `db:"order_items.order_id"`
		ProductName  string     `db:"order_items.product_name"`
		Price        float64    `db:"order_items.price"`
		Quantity     *int32     `db:"order_items.quantity"`
	}

	const query = `
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

	for range b.N {
		results := make(results, 0, 50000)
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

func BenchmarkSqlxOneResult(b *testing.B) {
	dbx := sqlx.NewDb(db, "postgres")

	type OrderWithItems struct {
		model.Orders
		Itens []model.OrderItems
	}

	type results []struct {
		ID           int32      `db:"orders.id"`
		CustomerName string     `db:"orders.customer_name"`
		CreatedAt    *time.Time `db:"orders.created_at"`
		OrderItemID  int32      `db:"order_items.id"`
		OrderID      *int32     `db:"order_items.order_id"`
		ProductName  string     `db:"order_items.product_name"`
		Price        float64    `db:"order_items.price"`
		Quantity     *int32     `db:"order_items.quantity"`
	}

	const query = `
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
	WHERE orders.id = 1
	ORDER BY orders.id ASC;
	`

	for range b.N {
		results := make(results, 0, 5)
		err := dbx.Select(&results, query)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}

		var order OrderWithItems
		orderIdx := make(map[int32]int)
		for _, row := range results {
			idx, ok := orderIdx[row.ID]
			if !ok {
				order = OrderWithItems{
					Orders: model.Orders{
						ID:           row.ID,
						CustomerName: row.CustomerName,
						CreatedAt:    row.CreatedAt,
					},
				}
				orderIdx[row.ID] = idx
			}
			order.Itens = append(order.Itens, model.OrderItems{
				ID:          row.OrderItemID,
				OrderID:     row.OrderID,
				ProductName: row.ProductName,
				Price:       row.Price,
				Quantity:    row.Quantity,
			})
		}

		if len(order.Itens) != 5 {
			b.Fatalf("expected 5 itens, got %d", len(order.Itens))
		}
	}
}

func BenchmarkCarta(b *testing.B) {
	type orders []struct {
		ID           int32      `db:"orders.id"`
		CustomerName string     `db:"orders.customer_name"`
		CreatedAt    *time.Time `db:"orders.created_at"`
		Itens        []struct {
			OrderItemID int32   `db:"order_items.id"`
			OrderID     *int32  `db:"order_items.order_id"`
			ProductName string  `db:"order_items.product_name"`
			Price       float64 `db:"order_items.price"`
			Quantity    *int32  `db:"order_items.quantity"`
		}
	}

	const query = `
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

	for range b.N {
		rows, err := db.Query(query)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}

		orders := make(orders, 0, 50000)
		err = carta.Map(rows, &orders)
		if err != nil {
			b.Fatalf("mapping failed: %v", err)
		}

		if len(orders) != 50000 {
			b.Fatalf("expected 50000 results, got %d", len(orders))
		}

		if len(orders[0].Itens) != 5 {
			b.Fatalf("expected 5 itens, got %d", len(orders[0].Itens))
		}
	}
}

func BenchmarkCartaOneResult(b *testing.B) {
	type order struct {
		ID           int32      `db:"orders.id"`
		CustomerName string     `db:"orders.customer_name"`
		CreatedAt    *time.Time `db:"orders.created_at"`
		Itens        []struct {
			OrderItemID int32   `db:"order_items.id"`
			OrderID     *int32  `db:"order_items.order_id"`
			ProductName string  `db:"order_items.product_name"`
			Price       float64 `db:"order_items.price"`
			Quantity    *int32  `db:"order_items.quantity"`
		}
	}

	const query = `
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
	WHERE orders.id = 1
	ORDER BY orders.id ASC;
	`

	for range b.N {
		rows, err := db.QueryContext(b.Context(), query)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}

		order := order{}
		err = carta.Map(rows, &order)
		if err != nil {
			b.Fatalf("mapping failed: %v", err)
		}

		if len(order.Itens) != 5 {
			b.Fatalf("expected 5 itens, got %d", len(order.Itens))
		}
	}
}

func BenchmarkGorm(b *testing.B) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost port=5432 user=postgres password=admin dbname=order sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Fatalf("failed to initialize GORM: %v", err)
	}

	b.ResetTimer()

	for range b.N {
		var orders []OrderWithItems
		err := gormDB.Preload("Itens").Find(&orders).Error
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}

		if len(orders) != 50000 {
			b.Fatalf("expected 50000 results, got %d", len(orders))
		}

		if len(orders[0].Itens) != 5 {
			b.Fatalf("expected 5 itens, got %d", len(orders[0].Itens))
		}
	}
}

func BenchmarkGormOneResult(b *testing.B) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost port=5432 user=postgres password=admin dbname=order sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Fatalf("failed to initialize GORM: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var order OrderWithItems
		err := gormDB.Preload("Itens").First(&order, "id = ?", 1).Error
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}

		if len(order.Itens) != 5 {
			b.Fatalf("expected 5 itens, got %d", len(order.Itens))
		}
	}
}

func BenchmarkPq(b *testing.B) {
	const query = `
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

	type OrderWithItems struct {
		model.Orders
		Itens []model.OrderItems
	}

	for range b.N {
		rows, err := db.Query(query)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}
		defer rows.Close()

		var orders []OrderWithItems
		orderIdx := make(map[int32]int)
		for rows.Next() {
			var row struct {
				ID           int32
				CustomerName string
				CreatedAt    *time.Time
				OrderItemID  int32
				OrderID      *int32
				ProductName  string
				Price        float64
				Quantity     *int32
			}
			if err := rows.Scan(&row.ID, &row.CustomerName, &row.CreatedAt, &row.OrderItemID, &row.OrderID, &row.ProductName, &row.Price, &row.Quantity); err != nil {
				b.Fatalf("row scan failed: %v", err)
			}

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

func BenchmarkPqOneResult(b *testing.B) {
	const query = `
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
		WHERE orders.id = 1
		ORDER BY orders.id ASC;
	`

	type OrderWithItems struct {
		model.Orders
		Itens []model.OrderItems
	}

	for range b.N {
		rows, err := db.Query(query)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}
		defer rows.Close()

		var order OrderWithItems
		for rows.Next() {
			var row struct {
				ID           int32
				CustomerName string
				CreatedAt    *time.Time
				OrderItemID  int32
				OrderID      *int32
				ProductName  string
				Price        float64
				Quantity     *int32
			}
			if err := rows.Scan(&row.ID, &row.CustomerName, &row.CreatedAt, &row.OrderItemID, &row.OrderID, &row.ProductName, &row.Price, &row.Quantity); err != nil {
				b.Fatalf("row scan failed: %v", err)
			}

			order.Itens = append(order.Itens, model.OrderItems{
				ID:          row.OrderItemID,
				OrderID:     row.OrderID,
				ProductName: row.ProductName,
				Price:       row.Price,
				Quantity:    row.Quantity,
			})
		}

		if len(order.Itens) != 5 {
			b.Fatalf("expected 5 itens, got %d", len(order.Itens))
		}
	}
}

func BenchmarkPqJsonAgg(b *testing.B) {
	const query = `
		SELECT orders.id AS "orders.id",
			   orders.customer_name AS "orders.customer_name",
			   orders.created_at AS "orders.created_at",
			   json_agg(json_build_object(
				   'id', order_items.id,
				   'order_id', order_items.order_id,
				   'product_name', order_items.product_name,
				   'price', order_items.price,
				   'quantity', order_items.quantity
			   )) AS "order_items"
		FROM public.orders
		INNER JOIN public.order_items ON (orders.id = order_items.order_id)
		GROUP BY orders.id, orders.customer_name, orders.created_at
		ORDER BY orders.id ASC;
	`

	type OrderItem struct {
		OrderItemID int32   `json:"id"`
		OrderID     *int32  `json:"order_id"`
		ProductName string  `json:"product_name"`
		Price       float64 `json:"price"`
		Quantity    *int32  `json:"quantity"`
	}

	type OrderWithItems struct {
		ID           int32       `db:"orders.id"`
		CustomerName string      `db:"orders.customer_name"`
		CreatedAt    *time.Time  `db:"orders.created_at"`
		Itens        []OrderItem `json:"order_items"`
	}

	for range b.N {
		orders := make([]OrderWithItems, 0, 50000)
		rows, err := db.QueryContext(b.Context(), query)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var (
				id             int32
				customerName   string
				createdAt      *time.Time
				orderItemsJSON []byte
			)
			if err := rows.Scan(&id, &customerName, &createdAt, &orderItemsJSON); err != nil {
				b.Fatalf("row scan failed: %v", err)
			}
			var itens []OrderItem
			if err := json.Unmarshal(orderItemsJSON, &itens); err != nil {
				b.Fatalf("failed to unmarshal order items: %v", err)
			}
			orders = append(orders, OrderWithItems{
				ID:           id,
				CustomerName: customerName,
				CreatedAt:    createdAt,
				Itens:        itens,
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

func BenchmarkPqJsonAggOneResult(b *testing.B) {
	const query = `
		SELECT orders.id AS "orders.id",
			   orders.customer_name AS "orders.customer_name",
			   orders.created_at AS "orders.created_at",
			   json_agg(json_build_object(
				   'id', order_items.id,
				   'order_id', order_items.order_id,
				   'product_name', order_items.product_name,
				   'price', order_items.price,
				   'quantity', order_items.quantity
			   )) AS "order_items"
		FROM public.orders
		INNER JOIN public.order_items ON (orders.id = order_items.order_id)
		WHERE orders.id = 1
		GROUP BY orders.id, orders.customer_name, orders.created_at
		ORDER BY orders.id ASC;
	`

	type OrderItem struct {
		OrderItemID int32   `json:"id"`
		OrderID     *int32  `json:"order_id"`
		ProductName string  `json:"product_name"`
		Price       float64 `json:"price"`
		Quantity    *int32  `json:"quantity"`
	}

	type OrderWithItems struct {
		ID           int32       `db:"orders.id"`
		CustomerName string      `db:"orders.customer_name"`
		CreatedAt    *time.Time  `db:"orders.created_at"`
		Itens        []OrderItem `json:"order_items"`
	}

	for range b.N {
		var order OrderWithItems
		rows, err := db.QueryContext(b.Context(), query)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}

		if rows.Next() {
			var orderItemsJSON []byte
			if err := rows.Scan(&order.ID, &order.CustomerName, &order.CreatedAt, &orderItemsJSON); err != nil {
				b.Fatalf("row scan failed: %v", err)
			}
			if err := json.Unmarshal(orderItemsJSON, &order.Itens); err != nil {
				b.Fatalf("failed to unmarshal order items: %v", err)
			}
		}

		if len(order.Itens) != 5 {
			b.Fatalf("expected 5 itens, got %d", len(order.Itens))
		}
	}
}
