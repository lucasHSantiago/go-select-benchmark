# Go Database SELECT Benchmark

This project benchmarks and compares the performance of different Go database libraries for executing and mapping SQL `SELECT` queries. It leverages Go's built-in testing and benchmarking tools, and provides a way to visualize and analyze benchmark results in a human-friendly format.

## Features

- **Benchmarking**: Uses Go's `testing` package to run performance benchmarks on different Go database libraries (Jet, Sqlx, Carta, GORM, and pq) for executing and mapping SQL `SELECT` queries.
- **Pretty Output**: Integrates with the [`prettybenchmarks`](https://github.com/florianorben/prettybenchmarks) tool to format benchmark results into readable tables, supporting both standard and memory allocation benchmarks (`-benchmem`).
- **Docker Support**: Includes a `docker-compose.yaml` for easy setup and reproducibility.
- **Database Migrations**: Contains a `migration/` directory for managing database schema changes required by the benchmarks.
- **Makefile**: Provides common build and test commands for convenience.

## How to Run

1. **Install Dependencies**
   ```sh
   go mod download
   ```
2. **Run Docker**
   ```sh
   docker-compose up -d
   ```

3. **Run Migrations**
   ```sh
   make migrateup
   ```

4. **Run Benchmarks**
   ```sh
   make benchmark
   ```
   - The `pb` command is provided by the `prettybenchmarks` tool. Install it if you haven't:
     ```sh
     go install github.com/florianorben/prettybenchmarks/pb@latest
     ```
   - You can change the time unit (`ms`, `us`, `ns`, `s`) as needed.


## Project Structure

- `main_test.go` — Contains Go benchmark tests for Jet, Sqlx, Carta, GORM, and pq database libraries.
- `migration/` — SQL migration scripts for database setup/teardown.
- `docker-compose.yaml` — Docker Compose configuration for running the project in containers.
- `makefile` — Common build, test, and utility commands.
- `go.mod`, `go.sum` — Go module dependencies.

## What Was Done

- Set up a Go project for benchmarking the performance of different Go database libraries for SQL `SELECT` queries.
- Added Docker Compose for containerized development and testing.
- Integrated the `prettybenchmarks` tool for improved benchmark result readability.
- Included database migration scripts for scenarios requiring persistent storage.
- Provided a Makefile for streamlined development workflows.
- Added benchmarks for GORM, including tests for fetching multiple results and one result.
- Updated the README to reflect the inclusion of GORM in the benchmarking suite.

## Results In My Machine

```
go test -bench=. -benchmem -benchtime=10s | prettybenchmarks ms

+--------------------+-------+-----------+-------------+----------------+
| Name               |  Runs |     ms/op |        B/op | allocations/op |
+--------------------+-------+-----------+-------------+----------------+
| Carta              |    12 |   998.322 | 332,259,377 |      9,347,193 |
+--------------------+-------+-----------+-------------+----------------+
| CartaOneResult     | 1,293 |     9.242 |       8,825 |            193 |
+--------------------+-------+-----------+-------------+----------------+
| Gorm               |    15 |   756.696 | 170,669,907 |      5,998,053 |
+--------------------+-------+-----------+-------------+----------------+
| GormOneResult      | 1,282 |     9.399 |      17,055 |            271 |
+--------------------+-------+-----------+-------------+----------------+
| Jet                |     7 | 1,627.808 | 627,821,147 |     14,849,726 |
+--------------------+-------+-----------+-------------+----------------+
| JetOneResult       | 1,269 |     9.462 |      26,522 |            544 |
+--------------------+-------+-----------+-------------+----------------+
| Pq                 |    20 |   547.595 | 128,585,016 |      4,696,459 |
+--------------------+-------+-----------+-------------+----------------+
| PqOneResult        | 1,303 |     9.244 |       4,057 |             86 |
+--------------------+-------+-----------+-------------+----------------+
| Sqlx               |    15 |   737.327 | 248,427,986 |      5,696,475 |
+--------------------+-------+-----------+-------------+----------------+
| SqlxOneResult      | 1,320 |     9.204 |       4,225 |            108 |
+--------------------+-------+-----------+-------------+----------------+

Summary:
+------+
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i5-9400F CPU @ 2.90GHz
```

## References
- [Go Benchmarking Documentation](https://golang.org/pkg/testing/#hdr-Benchmarks)
- [prettybenchmarks](https://github.com/florianorben/prettybenchmarks)

---

Feel free to modify the benchmarks or add new database libraries or query patterns to expand the analysis!
