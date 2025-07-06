#!/bin/bash

# Quick start script for POS System

echo "🚀 Starting POS System..."

# Check if Docker is available
if command -v docker &> /dev/null && command -v docker-compose &> /dev/null; then
    echo "🐳 Starting PostgreSQL with Docker..."
    docker-compose up -d postgres
    
    # Wait for PostgreSQL to be ready
    echo "⏳ Waiting for PostgreSQL to be ready..."
    sleep 10
    
    # Check if database is ready
    until docker-compose exec postgres pg_isready -U postgres; do
        echo "⏳ Waiting for PostgreSQL..."
        sleep 2
    done
    
    echo "✅ PostgreSQL is ready!"
else
    echo "⚠️  Docker not found. Please make sure PostgreSQL is running manually."
    echo "   You can install PostgreSQL with: brew install postgresql"
    echo "   Then start it with: brew services start postgresql"
fi

echo ""
echo "🌱 Seeding database with test data..."
go run cmd/seed/main.go

echo ""
echo "🚀 Starting POS System API..."
go run cmd/api/main.go
