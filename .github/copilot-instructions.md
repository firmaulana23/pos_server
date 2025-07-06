<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# POS System Coffee Shop - Copilot Instructions

This is a Point of Sale (POS) system backend API for a coffee shop built with:
- **Backend**: Go (Golang) with Gin framework
- **Database**: PostgreSQL
- **Authentication**: JWT tokens
- **Frontend**: Admin dashboard with HTML/CSS/JavaScript

## Project Structure
- `cmd/` - Application entry points
- `internal/` - Private application code
- `pkg/` - Public library code
- `web/` - Web assets (admin dashboard)
- `migrations/` - Database migrations
- `configs/` - Configuration files

## Key Features
- User authentication and authorization (JWT)
- Menu management with cost calculation (HPP) and margin analysis
- Transaction processing (POS) with status (pending/paid)
- Add-on support for menu items in transactions
- Payment methods support
- Expense tracking (raw materials and operational costs)
- Admin dashboard with analytics
- Cross-platform API design

## Development Guidelines
- Follow Go best practices and conventions
- Use proper error handling with custom error types
- Implement proper logging with structured logs
- Use dependency injection pattern
- Write comprehensive tests
- Follow RESTful API design principles
- Implement proper validation for all inputs
- Use PostgreSQL transactions for data consistency
