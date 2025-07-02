package main

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"

	. "github.com/go-jet/jet/v2/postgres"
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
	}
}
