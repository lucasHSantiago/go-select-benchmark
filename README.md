# Go Database SELECT Benchmark

This project benchmarks and compares the performance of different Go database libraries for executing and mapping SQL `SELECT` queries. It leverages Go's built-in testing and benchmarking tools, and provides a way to visualize and analyze benchmark results in a human-friendly format.

## Features

- **Benchmarking**: Uses Go's `testing` package to run performance benchmarks on different Go database libraries (Jet, Sqlx, Carta) for executing and mapping SQL `SELECT` queries.
- **Pretty Output**: Integrates with the [`prettybenchmarks`](https://github.com/florianorben/prettybenchmarks) tool to format benchmark results into readable tables, supporting both standard and memory allocation benchmarks (`-benchmem`).
- **Docker Support**: Includes a `docker-compose.yaml` for easy setup and reproducibility.
- **Database Migrations**: Contains a `migration/` directory for managing database schema changes required by the benchmarks.
- **Makefile**: Provides common build and test commands for convenience.

## How to Run

1. **Install Dependencies**
   ```sh
   go mod download
   ```

2. **Run Benchmarks**
   ```sh
   make benchmark
   ```
   - The `pb` command is provided by the `prettybenchmarks` tool. Install it if you haven't:
     ```sh
     go install github.com/florianorben/prettybenchmarks/pb@latest
     ```
   - You can change the time unit (`ms`, `us`, `ns`, `s`) as needed.

3. **Using Docker**
   - To run the project in a Docker environment:
     ```sh
     docker-compose up --build
     ```

## Project Structure

- `main_test.go` — Contains Go benchmark tests for Jet, Sqlx, and Carta database libraries.
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

## References
- [Go Benchmarking Documentation](https://golang.org/pkg/testing/#hdr-Benchmarks)
- [prettybenchmarks](https://github.com/florianorben/prettybenchmarks)

---

Feel free to modify the benchmarks or add new database libraries or query patterns to expand the analysis!
