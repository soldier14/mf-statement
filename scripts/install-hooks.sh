#!/bin/bash

# Install Git Hooks for MF Statement
# This script installs pre-commit hooks to ensure code quality

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_info() {
    echo -e "${BLUE}${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

echo "🔧 Installing Git hooks for MF Statement..."

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    print_error "Not in a git repository. Please run this script from the project root."
    exit 1
fi

# Check if Makefile exists
if [ ! -f "Makefile" ]; then
    print_error "Makefile not found. Please run this script from the project root."
    exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Install pre-commit hook
print_info "Installing pre-commit hook..."
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash

# MF Statement Pre-commit Hook
# This hook runs linting and tests before each commit

set -e

echo " Running pre-commit checks..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "Makefile" ]; then
    print_error "Makefile not found. Please run this hook from the project root."
    exit 1
fi

# 1. Run linting
echo "📝 Running linting..."
if make lint; then
    print_status "Linting passed"
else
    print_error "Linting failed. Please fix the issues and try again."
    exit 1
fi

# 2. Run tests
echo " Running tests..."
if make test; then
    print_status "Tests passed"
else
    print_error "Tests failed. Please fix the failing tests and try again."
    exit 1
fi

# 3. Check for uncommitted changes after linting
if ! git diff --quiet; then
    print_warning "Code was reformatted by go fmt. Please review and commit the changes."
    echo "Modified files:"
    git diff --name-only
    echo ""
    echo "You can run 'git add .' to stage the formatting changes and commit again."
    exit 1
fi

print_status "All pre-commit checks passed! "
echo "Ready to commit."
EOF

# Make the hook executable
chmod +x .git/hooks/pre-commit

print_status "Pre-commit hook installed successfully!"

# Create a commit-msg hook for conventional commits (optional)
print_info "Installing commit-msg hook for conventional commits..."
cat > .git/hooks/commit-msg << 'EOF'
#!/bin/bash

# MF Statement Commit Message Hook
# Validates commit messages follow conventional commit format

commit_regex='^(feat|fix|docs|style|refactor|test|chore)(\(.+\))?: .{1,50}'

if ! grep -qE "$commit_regex" "$1"; then
    echo " Invalid commit message format!"
    echo ""
    echo "Please use conventional commit format:"
    echo "  feat: add new feature"
    echo "  fix: fix bug"
    echo "  docs: update documentation"
    echo "  style: formatting changes"
    echo "  refactor: code refactoring"
    echo "  test: add or update tests"
    echo "  chore: maintenance tasks"
    echo ""
    echo "Examples:"
    echo "  feat: add make lint command"
    echo "  fix: resolve test failures"
    echo "  docs: update README"
    exit 1
fi
EOF

chmod +x .git/hooks/commit-msg

print_status "Commit message hook installed successfully!"

# Install pre-push hook
print_info "Installing pre-push hook for coverage enforcement..."
cat > .git/hooks/pre-push << 'EOF'
#!/bin/bash

# MF Statement Pre-push Hook
# This hook enforces test coverage requirements before allowing pushes

set -e

echo " Running pre-push checks..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_info() {
    echo -e "${BLUE}${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "Makefile" ]; then
    print_error "Makefile not found. Please run this hook from the project root."
    exit 1
fi

# Coverage requirements
TOTAL_COVERAGE_MIN=80
TEST_COVERAGE_MIN=90

print_info "Running coverage analysis..."
print_info "Requirements: Total coverage ≥ ${TOTAL_COVERAGE_MIN}%, Test coverage ≥ ${TEST_COVERAGE_MIN}%"

# Run coverage analysis
if ! make coverage > /tmp/coverage_output.txt 2>&1; then
    print_error "Coverage analysis failed. Please check the output:"
    cat /tmp/coverage_output.txt
    exit 1
fi

# Extract coverage percentages from the output
TOTAL_COVERAGE=$(grep -o "total:.*[0-9.]*%" /tmp/coverage_output.txt | grep -o "[0-9.]*" | tail -1)
TEST_COVERAGE=$(grep -o "internal.*[0-9.]*%" /tmp/coverage_output.txt | grep -o "[0-9.]*" | tail -1)

# Clean up temp file
rm -f /tmp/coverage_output.txt

# Check if we got valid coverage numbers
if [ -z "$TOTAL_COVERAGE" ] || [ -z "$TEST_COVERAGE" ]; then
    print_error "Could not extract coverage percentages from output."
    print_info "Please run 'make coverage' manually to check the output format."
    exit 1
fi

print_info "Coverage Results:"
print_info "  Total coverage: ${TOTAL_COVERAGE}%"
print_info "  Test coverage: ${TEST_COVERAGE}%"

# Check total coverage requirement
if (( $(echo "$TOTAL_COVERAGE < $TOTAL_COVERAGE_MIN" | bc -l) )); then
    print_error "Total coverage ${TOTAL_COVERAGE}% is below required ${TOTAL_COVERAGE_MIN}%"
    print_info "Please add more tests to increase coverage."
    exit 1
fi

# Check test coverage requirement
if (( $(echo "$TEST_COVERAGE < $TEST_COVERAGE_MIN" | bc -l) )); then
    print_error "Test coverage ${TEST_COVERAGE}% is below required ${TEST_COVERAGE_MIN}%"
    print_info "Please add more tests to increase test coverage."
    exit 1
fi

print_status "Coverage requirements met!"
print_status "Total coverage: ${TOTAL_COVERAGE}% (≥ ${TOTAL_COVERAGE_MIN}%)"
print_status "Test coverage: ${TEST_COVERAGE}% (≥ ${TEST_COVERAGE_MIN}%)"

print_status "All pre-push checks passed! "
echo "Ready to push."
EOF

chmod +x .git/hooks/pre-push

print_status "Pre-push hook installed successfully!"

echo ""
print_info "Git hooks installed! "
echo ""
echo "The following hooks are now active:"
echo "  • pre-commit: Runs linting and tests before each commit"
echo "  • commit-msg: Validates commit message format"
echo "  • pre-push: Enforces coverage requirements (≥80% total, ≥90% test)"
echo ""
echo "To bypass hooks (not recommended):"
echo "  git commit --no-verify"
echo "  git push --no-verify"
echo ""
echo "To uninstall hooks:"
echo "  rm .git/hooks/pre-commit .git/hooks/commit-msg .git/hooks/pre-push"

