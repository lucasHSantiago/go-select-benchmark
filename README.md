
# Go Database SELECT Benchmark

This project benchmarks and compares the performance of different Go database libraries for executing and mapping SQL `SELECT` queries. It leverages Go's built-in testing and benchmarking tools, and provides a way to visualize and analyze benchmark results in a human-friendly format.


## Features

- **Benchmarking**: Uses Go's `testing` package to run performance benchmarks on different Go database libraries (Jet, Sqlx, Carta, GORM, pq) for executing and mapping SQL `SELECT` queries, including variants using `json_agg` for grouped results.
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

- `main_test.go` — Contains Go benchmark tests for Jet, Sqlx, Carta, GORM, pq, and `json_agg`/grouped query patterns.
- `migration/` — SQL migration scripts for database setup/teardown.
- `docker-compose.yaml` — Docker Compose configuration for running the project in containers.
- `makefile` — Common build, test, and utility commands.
- `go.mod`, `go.sum` — Go module dependencies.


## Results In My Machine

```
go test -bench=. -benchmem -benchtime=10s | prettybenchmarks ms

+------------------------+------+-----------+-------------+----------------+
| Name                   | Runs |     ms/op |        B/op | allocations/op |
+------------------------+------+-----------+-------------+----------------+
| Carta                  |    2 |   950.586 | 332,277,448 |      9,347,452 |
+------------------------+------+-----------+-------------+----------------+
| CartaOneResult         |   81 |    14.530 |      29,969 |            474 |
+------------------------+------+-----------+-------------+----------------+
| Gorm                   |    2 |   769.458 | 170,717,664 |      5,998,494 |
+------------------------+------+-----------+-------------+----------------+
| GormOneResult          |  129 |     9.187 |      17,798 |            278 |
+------------------------+------+-----------+-------------+----------------+
| Jet                    |    1 | 1,563.279 | 628,004,472 |     14,850,240 |
+------------------------+------+-----------+-------------+----------------+
| JetOneResult           |   81 |    14.518 |      47,495 |            819 |
+------------------------+------+-----------+-------------+----------------+
| Pq                     |    2 |   558.257 | 128,601,128 |      4,696,726 |
+------------------------+------+-----------+-------------+----------------+
| PqJsonAgg              |    2 |   640.142 | 120,461,120 |      1,949,959 |
+------------------------+------+-----------+-------------+----------------+
| PqJsonAggOneResult     |   79 |    14.933 |      25,501 |            328 |
+------------------------+------+-----------+-------------+----------------+
| PqOneResult            |   81 |    14.452 |      24,744 |            360 |
+------------------------+------+-----------+-------------+----------------+
| Sqlx                   |    2 |   729.780 | 248,448,536 |      5,696,735 |
+------------------------+------+-----------+-------------+----------------+
| SqlxOneResult          |   80 |    14.376 |      24,902 |            382 |
+------------------------+------+-----------+-------------+----------------+

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


Feel free to modify the benchmarks or add new database libraries or query patterns (including grouped/aggregated queries) to expand the analysis!
