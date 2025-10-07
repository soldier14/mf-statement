#!/bin/bash

# CI/CD Local Check Script
# This script runs the same checks as the GitHub Actions CI pipeline

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[OK]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "Makefile" ]; then
    print_error "Makefile not found. Please run this script from the project root."
    exit 1
fi

print_header "Starting CI/CD Local Check"
echo ""

# Step 1: Lint Check
print_header "Step 1: Lint Check"
echo "Running linting checks..."
if make lint; then
    print_status "Linting completed successfully"
else
    print_error "Linting failed"
    exit 1
fi
echo ""

# Step 2: Test Suite
print_header "Step 2: Test Suite"
echo "Running test suite..."
if make test; then
    print_status "All tests passed"
else
    print_error "Tests failed"
    exit 1
fi
echo ""

# Step 3: Coverage Analysis
print_header "Step 3: Coverage Analysis"
echo "Running coverage analysis..."
if make coverage; then
    print_status "Coverage analysis completed"
else
    print_error "Coverage analysis failed"
    exit 1
fi
echo ""

# Step 4: Coverage Threshold Check
print_header "Step 4: Coverage Threshold Check"
echo "Checking coverage thresholds..."

TOTAL_COVERAGE=$(go tool cover -func=coverage/coverage_filtered.out | grep "total:" | awk '{print $3}' | sed 's/%//')
echo "Total coverage: ${TOTAL_COVERAGE}%"

# Check if coverage meets requirements
if (( $(echo "$TOTAL_COVERAGE >= 80" | bc -l) )); then
    print_status "Total coverage (${TOTAL_COVERAGE}%) meets requirement (≥80%)"
else
    print_error "Total coverage (${TOTAL_COVERAGE}%) below requirement (≥80%)"
    exit 1
fi

# Check test coverage specifically - use the total coverage as test coverage
TEST_COVERAGE=$TOTAL_COVERAGE
if (( $(echo "$TEST_COVERAGE >= 90" | bc -l) )); then
    print_status "Test coverage (${TEST_COVERAGE}%) meets requirement (≥90%)"
else
    print_error "Test coverage (${TEST_COVERAGE}%) below requirement (≥90%)"
    exit 1
fi
echo ""

# Step 5: Build Check
print_header "Step 5: Build Check"
echo "Building application..."
if make build; then
    print_status "Build completed successfully"
else
    print_error "Build failed"
    exit 1
fi
echo ""

# Summary
print_header "CI/CD Local Check Summary"
print_status "All checks passed!"
print_status " Linting: PASSED"
print_status " Tests: PASSED"
print_status " Coverage: PASSED (${TOTAL_COVERAGE}%)"
print_status " Build: PASSED"
echo ""
print_info "Ready for commit and push!"
print_info "The pre-commit and pre-push hooks will run automatically."
