# ==============================================================================
# Migrate

.PHONY: migrateup
migrateup:
	migrate -path ./migration -database "postgresql://postgres:admin@localhost:5432/order?sslmode=disable" -verbose up $(or $(n))

.PHONY: migratedown
migratedown:
	migrate -path ./migration -database "postgresql://postgres:admin@localhost:5432/order?sslmode=disable" -verbose down $(or $(n))

.PHONY: new_migration
new_migration:
	migrate create -ext sql -dir migration -seq $(name)

# ==============================================================================
# Test

.PHONY: benchmark
benchmark:
	go test -bench=. -benchmem -parallel=1 | prettybenchmarks ms