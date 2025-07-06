#!/bin/bash

# POS System Database Setup Script

echo "🚀 Setting up POS System Database..."

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL is not installed. Please install it first."
    echo "On macOS with Homebrew: brew install postgresql"
    exit 1
fi

# Database configuration
DB_NAME="pos_system"
DB_USER="postgres"
DB_HOST="localhost"
DB_PORT="5432"

# Check if PostgreSQL service is running
if ! pg_isready -h $DB_HOST -p $DB_PORT; then
    echo "❌ PostgreSQL service is not running."
    echo "Start it with: brew services start postgresql"
    exit 1
fi

echo "✅ PostgreSQL is available"

# Create database if it doesn't exist
echo "📦 Creating database: $DB_NAME"
createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME 2>/dev/null || echo "Database might already exist, continuing..."

echo "✅ Database setup complete!"
echo ""
echo "🔧 Configuration:"
echo "   Database: $DB_NAME"
echo "   Host: $DB_HOST"
echo "   Port: $DB_PORT"
echo "   User: $DB_USER"
echo ""
echo "📝 Next steps:"
echo "   1. Update your .env file with the database credentials"
echo "   2. Run: go run cmd/api/main.go"
echo ""
