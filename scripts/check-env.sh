#!/bin/bash

# Development Environment Check Script
# This script checks if all required tools and services are available

echo "🔍 Checking POS System Development Environment..."
echo ""

# Check Go installation
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo "✅ Go is installed: $GO_VERSION"
else
    echo "❌ Go is not installed. Please install Go 1.24+ from https://golang.org/"
    exit 1
fi

# Check PostgreSQL
if command -v psql &> /dev/null; then
    PG_VERSION=$(psql --version | awk '{print $3}')
    echo "✅ PostgreSQL is installed: $PG_VERSION"
    
    # Check if PostgreSQL is running
    if pg_isready -q; then
        echo "✅ PostgreSQL service is running"
    else
        echo "⚠️  PostgreSQL is installed but not running"
        echo "   Start it with: brew services start postgresql (macOS)"
        echo "   Or: sudo systemctl start postgresql (Linux)"
    fi
else
    echo "⚠️  PostgreSQL is not installed"
    echo "   Install with: brew install postgresql (macOS)"
    echo "   Or: sudo apt-get install postgresql (Ubuntu)"
fi

# Check Docker (optional)
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    echo "✅ Docker is installed: $DOCKER_VERSION"
    
    if command -v docker-compose &> /dev/null; then
        DOCKER_COMPOSE_VERSION=$(docker-compose --version | awk '{print $3}' | sed 's/,//')
        echo "✅ Docker Compose is installed: $DOCKER_COMPOSE_VERSION"
    else
        echo "⚠️  Docker is installed but Docker Compose is missing"
    fi
else
    echo "⚠️  Docker is not installed (optional)"
    echo "   Install from: https://docs.docker.com/get-docker/"
fi

# Check Git
if command -v git &> /dev/null; then
    GIT_VERSION=$(git --version | awk '{print $3}')
    echo "✅ Git is installed: $GIT_VERSION"
else
    echo "❌ Git is not installed"
    echo "   Install with: brew install git (macOS)"
    echo "   Or: sudo apt-get install git (Ubuntu)"
fi

# Check Make
if command -v make &> /dev/null; then
    echo "✅ Make is available"
else
    echo "⚠️  Make is not available"
    echo "   Install with: xcode-select --install (macOS)"
    echo "   Or: sudo apt-get install build-essential (Ubuntu)"
fi

echo ""
echo "📋 Project Status:"

# Check if .env file exists
if [ -f ".env" ]; then
    echo "✅ .env file exists"
else
    echo "⚠️  .env file is missing"
    echo "   Create it with: cp .env.example .env"
fi

# Check if go.mod exists
if [ -f "go.mod" ]; then
    echo "✅ go.mod file exists"
else
    echo "❌ go.mod file is missing"
    echo "   Initialize with: go mod init pos-system"
fi

# Check if dependencies are installed
if go list -m all &> /dev/null; then
    echo "✅ Go dependencies are installed"
else
    echo "⚠️  Go dependencies need to be installed"
    echo "   Run: go mod tidy"
fi

echo ""
echo "🚀 Quick Start Commands:"
echo "   make dev        # Start development server"
echo "   make test       # Run tests"
echo "   make help       # Show all available commands"
echo ""

# Summary
echo "📊 Summary:"
CHECKS_PASSED=0
CHECKS_TOTAL=7

if command -v go &> /dev/null; then ((CHECKS_PASSED++)); fi
if command -v psql &> /dev/null; then ((CHECKS_PASSED++)); fi
if command -v docker &> /dev/null; then ((CHECKS_PASSED++)); fi
if command -v git &> /dev/null; then ((CHECKS_PASSED++)); fi
if command -v make &> /dev/null; then ((CHECKS_PASSED++)); fi
if [ -f ".env" ]; then ((CHECKS_PASSED++)); fi
if [ -f "go.mod" ]; then ((CHECKS_PASSED++)); fi

echo "   $CHECKS_PASSED/$CHECKS_TOTAL checks passed"

if [ $CHECKS_PASSED -eq $CHECKS_TOTAL ]; then
    echo "🎉 Your environment is ready for development!"
elif [ $CHECKS_PASSED -ge 5 ]; then
    echo "✅ Your environment is mostly ready. Fix the warnings above if needed."
else
    echo "⚠️  Your environment needs some setup. Please install the missing tools."
fi
