# Coffee Shop POS System

A comprehensive Point of Sale (POS) system for coffee shops built with Go (Golang) and PostgreSQL. This system includes user authentication, menu management with cost calculation, transaction processing with add-on support, expense tracking, and an admin dashboard with analytics.

## 🚀 Quick Start

### Prerequisites

- Go 1.24+ installed
- PostgreSQL 12+ installed (or Docker)
- Git

### Option 1: Quick Start with Docker (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd pos_system

# Copy environment file
cp .env.example .env

# Start with Docker (includes PostgreSQL)
./scripts/start.sh
```

### Option 2: Manual Setup

```bash
# 1. Clone and setup
git clone <repository-url>
cd pos_system
cp .env.example .env

# 2. Install dependencies
go mod tidy

# 3. Setup PostgreSQL
# On macOS: brew install postgresql
# On Ubuntu: sudo apt install postgresql postgresql-contrib

# 4. Create database
createdb pos_system

# 5. Seed database with test data
go run cmd/seed/main.go

# 6. Run the application
make dev
# or
go run cmd/api/main.go
```

### Option 3: Using Docker Compose

```bash
# Start PostgreSQL and Adminer
docker-compose up -d

# Seed database
go run cmd/seed/main.go

# Run the API
make dev
```

## 📱 Access the Application

- **API Server**: http://localhost:8080
- **Admin Dashboard**: http://localhost:8080/admin
- **POS Interface**: http://localhost:8080/admin/pos
- **Database Admin** (if using Docker): http://localhost:8081

### Test Accounts

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@pos.com | admin123 |
| Manager | manager@pos.com | manager123 |
| Cashier | cashier@pos.com | cashier123 |

### Quick API Testing

Test the new menu-dependent add-ons functionality:
```bash
# Run the API testing script
./test_api.sh

# Or manually test endpoints
curl http://localhost:8080/api/v1/public/menu-item-add-ons/4
```

## Features

- **User Authentication & Authorization**
  - JWT-based authentication
  - Role-based access control (admin, manager, cashier)
  - User management

- **Menu Management**
  - Categories and menu items
  - Price and cost of goods sold (COGS/HPP) tracking
  - Automatic margin calculation
  - **Menu-dependent add-ons** (global and menu-specific)
  - Smart add-on organization for better customer experience

- **Add-ons System**
  - **Global add-ons**: Available for all menu items (e.g., milk alternatives, sweeteners)
  - **Menu-specific add-ons**: Only available for specific items (e.g., latte art for lattes)
  - Advanced filtering and management
  - Contextual add-on selection in POS

- **Transaction Processing**
  - Point of Sale interface
  - Transaction status (pending/paid)
  - Multiple payment methods (cash, card, digital wallet)
  - Add-on support in transactions
  - Save or pay options

- **Expense Tracking**
  - Raw material expenses
  - Operational expenses
  - Expense categorization and reporting

- **Analytics Dashboard**
  - Sales reports and charts
  - Profit analysis
  - Top-selling items and add-ons
  - Expense analysis

- **Cross-Platform API**
  - RESTful API design
  - CORS support for web applications
  - Mobile-ready endpoints

## Tech Stack

- **Backend**: Go (Golang) with Gin framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens
- **Frontend**: HTML, CSS, JavaScript with Chart.js
- **API**: RESTful design with JSON responses

## Project Structure

```
pos_system/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── app/
│   │   └── app.go               # Application setup
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection and setup
│   ├── handlers/
│   │   ├── auth.go              # Authentication handlers
│   │   ├── menu.go              # Menu management handlers
│   │   ├── addon.go             # Add-on management handlers
│   │   ├── transaction.go       # Transaction handlers
│   │   ├── expense.go           # Expense handlers
│   │   └── dashboard.go         # Dashboard handlers
│   ├── middleware/
│   │   └── auth.go              # Authentication middleware
│   ├── models/
│   │   └── models.go            # Database models
│   └── routes/
│       └── routes.go            # Route definitions
├── pkg/
│   └── auth/
│       └── jwt.go               # JWT utilities
├── web/
│   ├── static/
│   │   ├── css/
│   │   │   └── style.css        # Stylesheet
│   │   └── js/
│   │       ├── auth.js          # Authentication utilities
│   │       ├── dashboard.js     # Dashboard functionality
│   │       └── pos.js           # POS functionality
│   └── templates/
│       ├── login.html           # Login page
│       ├── dashboard.html       # Dashboard page
│       └── pos.html             # POS page
├── .env.example                 # Environment variables example
├── go.mod                       # Go module file
└── README.md                    # This file
```

## Installation

### Prerequisites

- Go 1.23.4 or later
- PostgreSQL 12 or later
- Git

### Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd pos_system
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL database**
   ```sql
   CREATE DATABASE pos_system;
   CREATE USER pos_user WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE pos_system TO pos_user;
   ```

4. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env file with your database credentials
   ```

5. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

6. **Access the application**
   - API: http://localhost:8080/api/v1
   - Admin Dashboard: http://localhost:8080/admin

## 🛠️ Development Commands

The project includes a Makefile for common development tasks:

```bash
# Show all available commands
make help

# Build the application
make build

# Run the application (build first)
make run

# Run in development mode (hot reload)
make dev

# Run tests
make test

# Clean build artifacts
make clean

# Install/update dependencies
make deps

# Format code
make fmt

# Setup database (requires PostgreSQL)
make setup-db

# Seed database with test data
go run cmd/seed/main.go
```

### Docker Commands

```bash
# Start PostgreSQL with Docker
docker-compose up -d postgres

# Start PostgreSQL and Adminer (database admin tool)
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f postgres
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test ./internal/models -run TestUserModel
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `GET /api/v1/profile` - Get user profile
- `PUT /api/v1/profile` - Update user profile

### Menu Management
- `GET /api/v1/menu/categories` - Get all categories
- `POST /api/v1/menu/categories` - Create category (admin/manager)
- `PUT /api/v1/menu/categories/:id` - Update category (admin/manager)
- `DELETE /api/v1/menu/categories/:id` - Delete category (admin/manager)
- `GET /api/v1/menu/items` - Get menu items
- `POST /api/v1/menu/items` - Create menu item (admin/manager)
- `PUT /api/v1/menu/items/:id` - Update menu item (admin/manager)
- `DELETE /api/v1/menu/items/:id` - Delete menu item (admin/manager)

### Add-ons Management
- `GET /api/v1/add-ons` - Get all add-ons
- `POST /api/v1/add-ons` - Create add-on (admin/manager)
- `PUT /api/v1/add-ons/:id` - Update add-on (admin/manager)
- `DELETE /api/v1/add-ons/:id` - Delete add-on (admin/manager)

### Transactions
- `GET /api/v1/transactions` - Get transactions
- `POST /api/v1/transactions` - Create transaction
- `GET /api/v1/transactions/:id` - Get transaction details
- `PUT /api/v1/transactions/:id/pay` - Process payment
- `DELETE /api/v1/transactions/:id` - Delete transaction (admin/manager)

### Expenses
- `GET /api/v1/expenses` - Get expenses
- `POST /api/v1/expenses` - Create expense (admin/manager)
- `PUT /api/v1/expenses/:id` - Update expense (admin/manager)
- `DELETE /api/v1/expenses/:id` - Delete expense (admin/manager)
- `GET /api/v1/expenses/summary` - Get expense summary

### Dashboard
- `GET /api/v1/dashboard/stats` - Get dashboard statistics
- `GET /api/v1/dashboard/sales-report` - Get sales report
- `GET /api/v1/dashboard/profit-analysis` - Get profit analysis

### Public Endpoints (No Authentication Required)
- `GET /api/v1/public/menu/categories` - Get categories
- `GET /api/v1/public/menu/items` - Get menu items
- `GET /api/v1/public/add-ons` - Get add-ons
- `GET /api/v1/public/payment-methods` - Get payment methods

## Usage

### Creating a Transaction with Add-ons

1. **Select menu items** and add them to cart
2. **Choose add-ons** for each item (optional)
3. **Review the cart** with calculated totals
4. **Save transaction** (status: pending) or **Process payment** (status: paid)

### Menu Item with Margin Calculation

When creating/updating menu items, provide:
- `price`: Selling price
- `cogs`: Cost of Goods Sold (HPP)
- System automatically calculates: `margin = (price - cogs) / price * 100`

### Expense Tracking

Track two types of expenses:
- `raw_material`: Coffee beans, milk, etc.
- `operational`: Rent, utilities, salaries, etc.

## Default Data

The system automatically creates:

### Payment Methods
- Cash
- Credit Card
- Digital Wallet

### Categories
- Coffee
- Tea
- Snacks
- Beverages

### Add-ons
- Extra Shot (Rp 5,000)
- Whipped Cream (Rp 3,000)
- Extra Milk (Rp 2,000)
- Vanilla Syrup (Rp 3,000)
- Caramel Syrup (Rp 3,000)
- Decaf (Free)

## User Roles

- **Admin**: Full system access
- **Manager**: Cannot manage users, but can access all other features
- **Cashier**: Can only process transactions and view limited data

## Development

### Running in Development Mode

```bash
# Enable hot reload with air (install first: go install github.com/cosmtrek/air@latest)
air

# Or run directly
go run cmd/api/main.go
```

### Database Migrations

The application uses GORM's auto-migration feature. Database schema is automatically updated when the application starts.

### Adding New Features

1. Add models to `internal/models/models.go`
2. Create handlers in `internal/handlers/`
3. Add routes in `internal/routes/routes.go`
4. Update frontend in `web/` directory

## Production Deployment

1. **Build the application**
   ```bash
   go build -o pos-system cmd/api/main.go
   ```

2. **Set production environment variables**
   ```bash
   export SERVER_HOST=0.0.0.0
   export SERVER_PORT=8080
   export DB_HOST=your-db-host
   export DB_PASSWORD=your-secure-password
   export JWT_SECRET=your-super-secure-secret
   ```

3. **Run the application**
   ```bash
   ./pos-system
   ```

## Security Features

- JWT-based authentication
- Password hashing with bcrypt
- Role-based authorization
- Input validation
- SQL injection prevention (GORM)
- CORS protection

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support or questions, please create an issue in the repository.
