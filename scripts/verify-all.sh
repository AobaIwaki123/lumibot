#!/usr/bin/env bash
# ==============================================================================
# verify-all.sh - Run all static checks, linters, and unit tests locally
# ==============================================================================

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

echo "========================================================"
echo "1. Checking Go formatting (gofmt)..."
echo "========================================================"
UNFORMATTED=$(gofmt -l .)
if [[ -n "${UNFORMATTED}" ]]; then
  echo "Error: The following files are not formatted properly:"
  echo "${UNFORMATTED}"
  echo "Run 'gofmt -w .' to fix."
  exit 1
fi
echo "OK: All Go files are formatted properly."
echo ""

echo "========================================================"
echo "2. Running Go linters (golangci-lint)..."
echo "========================================================"
if command -v golangci-lint >/dev/null 2>&1; then
  golangci-lint run ./...
  echo "OK: golangci-lint passed."
else
  echo "Warning: golangci-lint is not installed. Running 'go vet' instead."
  go vet ./...
  echo "OK: go vet passed."
fi
echo ""

echo "========================================================"
echo "3. Running Go tests with race detector and coverage..."
echo "========================================================"
go test -v -race -coverprofile=coverage.out ./...
echo "OK: All Go tests passed."
echo ""

echo "========================================================"
echo "4. Checking Go build..."
echo "========================================================"
go build -v -o /dev/null ./cmd/lumibot
echo "OK: Go build succeeded."
echo ""

echo "========================================================"
echo "All local verification checks passed successfully!"
echo "========================================================"
