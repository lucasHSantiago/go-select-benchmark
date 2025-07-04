//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Orders = newOrdersTable("public", "orders", "")

type ordersTable struct {
	postgres.Table

	// Columns
	ID           postgres.ColumnInteger
	CustomerName postgres.ColumnString
	CreatedAt    postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
	DefaultColumns postgres.ColumnList
}

type OrdersTable struct {
	ordersTable

	EXCLUDED ordersTable
}

// AS creates new OrdersTable with assigned alias
func (a OrdersTable) AS(alias string) *OrdersTable {
	return newOrdersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new OrdersTable with assigned schema name
func (a OrdersTable) FromSchema(schemaName string) *OrdersTable {
	return newOrdersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new OrdersTable with assigned table prefix
func (a OrdersTable) WithPrefix(prefix string) *OrdersTable {
	return newOrdersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new OrdersTable with assigned table suffix
func (a OrdersTable) WithSuffix(suffix string) *OrdersTable {
	return newOrdersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newOrdersTable(schemaName, tableName, alias string) *OrdersTable {
	return &OrdersTable{
		ordersTable: newOrdersTableImpl(schemaName, tableName, alias),
		EXCLUDED:    newOrdersTableImpl("", "excluded", ""),
	}
}

func newOrdersTableImpl(schemaName, tableName, alias string) ordersTable {
	var (
		IDColumn           = postgres.IntegerColumn("id")
		CustomerNameColumn = postgres.StringColumn("customer_name")
		CreatedAtColumn    = postgres.TimestampColumn("created_at")
		allColumns         = postgres.ColumnList{IDColumn, CustomerNameColumn, CreatedAtColumn}
		mutableColumns     = postgres.ColumnList{CustomerNameColumn, CreatedAtColumn}
		defaultColumns     = postgres.ColumnList{IDColumn, CreatedAtColumn}
	)

	return ordersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:           IDColumn,
		CustomerName: CustomerNameColumn,
		CreatedAt:    CreatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
