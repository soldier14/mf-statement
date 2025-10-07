# MF Statement CLI

A command-line tool for generating monthly financial statements from CSV transaction data.

- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
- [Configuration](#configuration)
- [Examples](#examples)
- [Development](#development)
- [Architecture](#architecture)

## 🛠 Installation

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

##  Quick Start

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

# Generate statement
./bin/mf-statement generate --period 202501 --csv transactions.csv
```

### Generate Command Options

| Flag | Short | Description | Required |
|------|-------|-------------|----------|
| `--period` | `-p` | Month in YYYYMM format (e.g., 202501) | Yes |
| `--csv` | `-c` | Path to CSV file | Yes |
| `--out` | `-o` | Output JSON file path (default: stdout) | No |
| `--verbose` | `-v` | Enable verbose logging | No |



### Coverage Report

The `make coverage` command generates a detailed coverage report:
- Excludes logger from coverage calculations
- Generates both console output and HTML report
- HTML report saved to `coverage/coverage.html`

### Git Hooks

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

##  CI/CD Pipeline

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
  "net_amount": 1650,
  "transaction_count": 3,
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
  ],
  "generated_at": "2025-01-15T10:30:00Z"
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

## Architecture

The application follows clean architecture principles with clear separation of concerns:

```
cmd/statement/          # Application entry point
internal/
├── cli/               # CLI interface and commands
├── config/             # Configuration management
├── domain/             # Domain models and business logic
├── usecase/            # Use cases and application logic
├── adapters/           # External adapters
│   ├── in/            # Input adapters (CLI)
│   └── out/           # Output adapters (parsers, writers)
└── util/              # Utility functions
```

### Key Components

- **Domain Layer**: Core business logic and entities
- **Use Case Layer**: Application-specific business rules
- **Adapter Layer**: External interfaces (CLI, file I/O, parsing)
- **Configuration**: Environment-based configuration
- **Error Handling**: Comprehensive error types and handling

##  Development

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


