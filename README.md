# MF Statement CLI

A command-line tool for generating monthly financial statements from CSV transaction data.

> **Assignment Documentation**: This project includes comprehensive technical documentation in [SOLUTION.md](SOLUTION.md) that addresses assignment requirements including thought process, technology choices, design decisions, requirement fulfillment, and future work considerations.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Examples](#examples)
- [Performance](#performance)
- [Development](#development)
- [Architecture](#architecture)
- [CI/CD Pipeline](#cicd-pipeline)
- [Git Hooks](#git-hooks)

## Features

- **CSV Processing**: Parse transaction data from CSV files
- **Monthly Statements**: Generate statements grouped by year and month
- **File Output**: Save statements to files or print to console
- **Memory Optimized**: Streaming parser with early filtering for large datasets
- **High Performance**: Process 1M transactions in under 0.3 seconds

## Installation

### Prerequisites

- Go 1.20 or later
- A CSV file with transaction data

### Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/mf-statement.git
   cd mf-statement
   ```

2. Build the application:
   ```bash
   # Using Makefile (recommended)
   make build
   
   # Or manually
   go mod tidy
   go build -o bin/mf-statement ./cmd/statement
   ```

3. Make it executable:
   ```bash
   chmod +x bin/mf-statement
   ```

## Quick Start

1. **Prepare your CSV file** with the following format:
   ```csv
   date,amount,content
   2025/01/05,2000,Salary
   2025/01/09,-300,Grocery
   2025/01/01,100,Gift
   ```

2. **Generate a statement**:
   ```bash
   ./bin/mf-statement generate --period 202501 --csv transactions.csv
   ```

3. **Save to file**:
   ```bash
   ./bin/mf-statement generate --period 202501 --csv transactions.csv --out statement.json
   ```

## Usage

### Basic Commands

```bash
# Show help
./bin/mf-statement --help

# Show version
./bin/mf-statement version

# Generate statement (standard)
./bin/mf-statement generate --period 202501 --csv transactions.csv

# Generate statement (optimized for large files)
./bin/mf-statement generate-optimized --period 202501 --csv transactions.csv
```

### Generate Command Options

| Flag | Short | Description | Required |
|------|-------|-------------|----------|
| `--period` | `-p` | Month in YYYYMM format (e.g., 202501) | Yes |
| `--csv` | `-c` | Path to CSV file | Yes |
| `--out` | `-o` | Output JSON file path (default: stdout) | No |
| `--verbose` | `-v` | Enable verbose logging | No |
| `--timeout` | `-t` | Timeout in seconds (default: 30) | No |

### Command Variants

| Command | Use Case | Memory Usage | Performance |
|---------|----------|--------------|-------------|
| `generate` | Standard processing | Higher memory | Good for small-medium files |
| `generate-optimized` | Large datasets | 90% less memory | 4.6x faster for large files |



## Git Hooks

The project includes pre-commit hooks to ensure code quality:

```bash
# Install git hooks
./scripts/install-hooks.sh

# Manual installation
chmod +x .git/hooks/pre-commit
```

**Pre-commit Hook Features:**
- Runs `make lint` (formatting + vetting)
- Runs `make test` (all tests)
- Checks for formatting changes after linting
- Prevents commits with failing tests or linting issues

**Pre-push Hook Features:**
- Enforces coverage requirements before push
- Requires ≥80% total coverage
- Requires ≥90% test coverage
- Prevents pushes with insufficient coverage

**Commit Message Hook:**
- Validates conventional commit format
- Examples: `feat: add new feature`, `fix: resolve bug`, `docs: update README`

## CI/CD Pipeline

The project includes automated CI/CD pipelines using GitHub Actions:

### **Continuous Integration (CI)**

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Pipeline Stages:**
1. **Lint Check** 
   - Code formatting (`go fmt`)
   - Static analysis (`go vet`)
   - Ensures code quality standards

2. **Test Suite** 
   - Runs all BDD tests with Ginkgo/Gomega
   - Validates functionality and behavior
   - Ensures no regressions

3. **Coverage Analysis** 
   - Generates detailed coverage reports
   - Enforces coverage thresholds:
     - **Total Coverage**: ≥80%
     - **Test Coverage**: ≥90%
   - Uploads coverage artifacts
   - Comments coverage on PRs

### **Continuous Deployment (CD)**

**Release Pipeline:**
- Triggered on version tags (`v*`)
- Builds multi-platform binaries:
  - Linux (AMD64, ARM64)
  - macOS (Intel, Apple Silicon)
  - Windows (AMD64)
- Creates GitHub releases with checksums
- Automated binary distribution

### **Quality Gates**

The pipeline will **fail** if:
-  Linting issues found
-  Any test fails
-  Coverage below thresholds
-  Build errors

The pipeline will **succeed** if:
-  All linting passes
-  All tests pass
-  Coverage meets requirements
-  Build successful

## Examples

### Example 1: Basic Statement Generation

```bash
# Input CSV (transactions.csv)
date,amount,content
2025/01/05,2000,Salary
2025/01/09,-300,Grocery
2025/01/01,100,Gift
2025/01/15,-150,Transport

# Generate statement
./bin/mf-statement generate --period 202501 --csv transactions.csv

# Output (JSON)
{
  "period": "2025/01",
  "total_income": 2100,
  "total_expenditure": -450,
  "transactions": [
    {
      "date": "2025/01/15",
      "amount": "-150",
      "content": "Transport"
    },
    {
      "date": "2025/01/09",
      "amount": "-300",
      "content": "Grocery"
    },
    {
      "date": "2025/01/05",
      "amount": "2000",
      "content": "Salary"
    }
  ]
}
```

### Example 2: Verbose Logging

```bash
./bin/mf-statement generate --period 202501 --csv transactions.csv --verbose
```

### Example 3: Custom Output File

```bash
./bin/mf-statement generate --period 202501 --csv transactions.csv --out monthly-statement.json
```

### Example 4: Optimized Processing for Large Files

```bash
# For large datasets (1M+ transactions)
./bin/mf-statement generate-optimized --period 202501 --csv large-transactions.csv --out statement.json --verbose

# Performance comparison
echo "Standard processing:"
time ./bin/mf-statement generate --period 202501 --csv large-transactions.csv

echo "Optimized processing:"
time ./bin/mf-statement generate-optimized --period 202501 --csv large-transactions.csv
```

## Performance

The application is optimized for high-performance processing of large transaction datasets:

### **Benchmark Results (1M Transactions)**

| Metric | Value |
|--------|-------|
| **Input File Size** | ~27.4 MB (1,000,001 lines) |
| **Processing Time** | ~0.28-0.32 seconds |
| **Memory Usage** | ~187 MB peak |
| **CPU Usage** | 115-122% (multi-threaded) |
| **Output Size** | ~2.1 MB JSON |

### **Performance by Time Period**

| Period | Processing Time | Transactions | Output Size |
|--------|----------------|--------------|-------------|
| **2025/01** | 0.286s | 22,454 | 2.14 MB |
| **2025/02** | 0.295s | 20,325 | 1.94 MB |
| **2023/12** | 0.304s | 22,511 | 2.15 MB |

### **Optimized vs Standard Performance**

| Version | Processing Time | Memory Usage | CPU Usage | Use Case |
|---------|----------------|--------------|-----------|----------|
| **Standard** | 0.971s | 204 MB | 43% | Small-medium files |
| **Optimized** | 0.212s | 21 MB | 101% | Large datasets (1M+ transactions) |
| **Improvement** | **78% faster** | **90% less memory** | **Better CPU utilization** | **4.6x performance boost** |

### **Performance Features**

**High Throughput**: Process 1M transactions in under 0.3 seconds  
**Memory Efficient**: Reasonable memory usage for large datasets  
**Historical Data**: Fast processing of historical transactions  
**Consistent Performance**: Stable performance across different time periods  
**Scalable**: Performance remains consistent with different timeout settings  

### **Benchmark Commands**

```bash
# Test with sample data (both commands)
time ./bin/mf-statement generate --csv testdata/transactions.sample.csv --period 202501 --out output.json
time ./bin/mf-statement generate-optimized --csv testdata/transactions.sample.csv --period 202501 --out output.json

# Test with large dataset (1M transactions)
mkdir -p testdata/output
time ./bin/mf-statement generate --csv testdata/transactions_1M.sample.csv --period 202501 --out testdata/output/statement-1M.json
time ./bin/mf-statement generate-optimized --csv testdata/transactions_1M.sample.csv --period 202501 --out testdata/output/statement-1M-optimized.json

# Performance comparison with large dataset
echo "Standard processing:"
time ./bin/mf-statement generate --csv testdata/transactions_1M.sample.csv --period 202501

echo "Optimized processing:"
time ./bin/mf-statement generate-optimized --csv testdata/transactions_1M.sample.csv --period 202501
```

## Architecture

The application follows clean architecture principles with clear separation of concerns. For a detailed technical explanation, see [SOLUTION.md](SOLUTION.md).

> **Technical Documentation**: For comprehensive details about the solution approach, design decisions, and implementation rationale, please read [SOLUTION.md](SOLUTION.md). This document provides in-depth answers to assignment questions including thought process, technology choices, architecture decisions, and future work considerations.

```
cmd/statement/          # Application entry point
internal/
├── cli/               # CLI interface and commands
├── domain/             # Domain models and business logic
├── usecase/            # Use cases and application logic
├── adapters/           # External adapters
│   ├── in/            # Input adapters (CLI)
│   └── out/           # Output adapters (parsers, writers)
└── util/              # Utility functions
testdata/              # Test fixtures and sample data
├── transactions.sample.csv    # Sample transaction data (small dataset)
└── transactions_1M.sample.csv # Large dataset for performance testing (1M transactions)
```

### Key Components

- **Domain Layer**: Core business logic and entities
- **Use Case Layer**: Application-specific business rules
- **Adapter Layer**: External interfaces (CLI, file I/O, parsing)
- **Configuration**: Environment-based configuration
- **Error Handling**: Comprehensive error types and handling

### Optimization Techniques

The application includes advanced optimization techniques for handling large datasets:

- **Streaming Parser**: Processes CSV data in streams to reduce memory usage
- **Early Filtering**: Filters transactions during parsing, not after loading
- **Memory Pooling**: Reuses memory buffers to reduce garbage collection
- **Buffer Optimization**: Optimized I/O buffer sizes for better performance
- **Lazy Loading**: Only loads relevant data into memory

**When to use optimized version:**
- Files with 100K+ transactions
- Memory-constrained environments
- Processing multiple large files
- Production environments with high throughput requirements

## Development

### Using Makefile

The project includes a simple Makefile for common development tasks:

```bash
# Show available commands
make help

# Build the application
make build

# Run all tests
make test

# Run tests with coverage (excluding logger)
make coverage

# Run linting (fmt + vet)
make lint

# Clean build artifacts
make clean

# Build and run the application
make run
```

### Manual Commands

If you prefer to run commands manually:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/usecase
```

### Code Style

The project follows Go best practices and includes:

- Comprehensive error handling
- Clean architecture principles
- Extensive testing
- Clear documentation
- Consistent naming conventions

### Adding New Features

1. Define domain models in `internal/domain/`
2. Implement use cases in `internal/usecase/`
3. Create adapters in `internal/adapters/`
4. Add CLI commands in `internal/cli/`
5. Write comprehensive tests
6. Update documentation

### Local CI Check

Run all CI checks locally before pushing:

```bash
# Run complete CI pipeline locally
make ci-check

# Individual checks
make lint      # Linting
make test      # Tests
make coverage  # Coverage analysis
make build     # Build verification
```

### Coverage Report

The `make coverage` command generates a detailed coverage report:
- Excludes logger from coverage calculations
- Generates both console output and HTML report
- HTML report saved to `coverage/coverage.html`
