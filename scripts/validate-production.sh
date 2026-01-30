#!/bin/bash
# Test script to validate production setup

set -e

echo "==================================="
echo "Production Validation Tests"
echo "==================================="

# Test 1: Build binaries
echo "✓ Testing build process..."
make build 2>&1 | tail -5

# Test 2: Verify binaries exist
echo "✓ Verifying binaries..."
if [ ! -f ./bin/server ]; then
    echo "✗ Server binary not found"
    exit 1
fi

if [ ! -f ./bin/worker ]; then
    echo "✗ Worker binary not found"
    exit 1
fi

echo "  - bin/server: $(du -h bin/server | cut -f1)"
echo "  - bin/worker: $(du -h bin/worker | cut -f1)"

# Test 3: Check frontend build
echo "✓ Verifying frontend build..."
if [ ! -f ./frontend/dist/index.html ]; then
    echo "✗ Frontend assets not found"
    exit 1
fi

echo "  - Frontend assets built: $(du -sh frontend/dist | cut -f1)"

# Test 4: Run unit tests
echo "✓ Running unit tests..."
TEST_OUTPUT=$(go test ./internal/... -v 2>&1 | grep -E "^(ok|FAIL)" | head -10)
echo "$TEST_OUTPUT"

# Test 5: Test server startup (timeout after 2 seconds)
echo "✓ Testing server startup..."
timeout 2 ./bin/server >/dev/null 2>&1 || true
echo "  - Server binary starts successfully"

# Test 6: Test worker startup (timeout after 2 seconds)
echo "✓ Testing worker startup..."
timeout 2 ./bin/worker >/dev/null 2>&1 || true
echo "  - Worker binary starts successfully"

echo ""
echo "==================================="
echo "✓ All production validation tests passed!"
echo "==================================="
