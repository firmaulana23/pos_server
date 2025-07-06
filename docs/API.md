# POS System API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

The API uses JWT Bearer tokens for authentication. Include the token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "email": "admin@pos.com",
    "password": "admin123"
}
```

**Response:**
```json
{
    "success": true,
    "message": "Login successful",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "username": "admin",
            "email": "admin@pos.com",
            "role": "admin",
            "is_active": true
        }
    }
}
```

## User Management

### Get Profile
```http
GET /api/v1/profile
Authorization: Bearer <token>
```

### Register User (Admin only)
```http
POST /api/v1/auth/register
Authorization: Bearer <admin_token>
Content-Type: application/json

{
    "username": "newuser",
    "email": "newuser@pos.com",
    "password": "password123",
    "role": "cashier"
}
```

## Menu Management

### Get Categories
```http
GET /api/v1/menu/categories
```

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "name": "Coffee",
            "description": "Hot and cold coffee beverages",
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
        }
    ]
}
```

### Create Category (Admin/Manager)
```http
POST /api/v1/menu/categories
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "New Category",
    "description": "Category description"
}
```

### Get Menu Items
```http
GET /api/v1/menu/items?category_id=1&available=true
```

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "category_id": 1,
            "name": "Espresso",
            "description": "Rich and strong coffee shot",
            "price": 15000,
            "cogs": 8000,
            "margin": 46.67,
            "is_available": true,
            "image_url": "",
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z",
            "category": {
                "id": 1,
                "name": "Coffee"
            }
        }
    ]
}
```

### Create Menu Item (Admin/Manager)
```http
POST /api/v1/menu/items
Authorization: Bearer <token>
Content-Type: application/json

{
    "category_id": 1,
    "name": "New Coffee",
    "description": "Delicious new coffee",
    "price": 25000,
    "cogs": 12000,
    "is_available": true
}
```

## Add-ons Management

### Get Add-ons
```http
GET /api/v1/add-ons?available=true
```

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "name": "Extra Shot",
            "description": "Additional espresso shot",
            "price": 8000,
            "cogs": 4000,
            "margin": 50.0,
            "is_available": true,
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
        }
    ]
}
```

### Create Add-on (Admin/Manager)
```http
POST /api/v1/add-ons
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "Oat Milk",
    "description": "Premium oat milk substitute",
    "price": 7000,
    "cogs": 4000,
    "is_available": true
}
```

## Transactions

### Create Transaction
```http
POST /api/v1/transactions
Authorization: Bearer <token>
Content-Type: application/json

{
    "items": [
        {
            "menu_item_id": 1,
            "quantity": 2,
            "add_ons": [
                {
                    "add_on_id": 1,
                    "quantity": 1
                }
            ]
        }
    ],
    "payment_method": "cash",
    "tax": 2500,
    "discount": 0
}
```

**Response:**
```json
{
    "success": true,
    "message": "Transaction created successfully",
    "data": {
        "id": 1,
        "transaction_no": "TXN-20240101-0001",
        "status": "pending",
        "payment_method": "cash",
        "sub_total": 46000,
        "tax": 2500,
        "discount": 0,
        "total": 48500,
        "created_at": "2024-01-01T12:00:00Z",
        "items": [
            {
                "id": 1,
                "menu_item_id": 1,
                "quantity": 2,
                "unit_price": 15000,
                "total_price": 30000,
                "menu_item": {
                    "id": 1,
                    "name": "Espresso",
                    "price": 15000
                },
                "add_ons": [
                    {
                        "id": 1,
                        "add_on_id": 1,
                        "quantity": 1,
                        "unit_price": 8000,
                        "total_price": 8000,
                        "add_on": {
                            "id": 1,
                            "name": "Extra Shot",
                            "price": 8000
                        }
                    }
                ]
            }
        ]
    }
}
```

### Process Payment
```http
PUT /api/v1/transactions/1/pay
Authorization: Bearer <token>
Content-Type: application/json

{
    "payment_method": "cash"
}
```

### Get Transactions
```http
GET /api/v1/transactions?status=paid&limit=10&offset=0
Authorization: Bearer <token>
```

## Expenses

### Get Expenses
```http
GET /api/v1/expenses?type=raw_material&start_date=2024-01-01&end_date=2024-01-31
Authorization: Bearer <token>
```

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "type": "raw_material",
            "category": "Coffee Beans",
            "description": "Premium Arabica coffee beans - 5kg",
            "amount": 500000,
            "date": "2024-01-01T00:00:00Z",
            "user_id": 1,
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z",
            "user": {
                "id": 1,
                "username": "admin",
                "role": "admin"
            }
        }
    ]
}
```

### Create Expense (Admin/Manager)
```http
POST /api/v1/expenses
Authorization: Bearer <token>
Content-Type: application/json

{
    "type": "operational",
    "category": "Utilities",
    "description": "Monthly electricity bill",
    "amount": 800000,
    "date": "2024-01-01T00:00:00Z"
}
```

## Dashboard Analytics

### Get Dashboard Stats
```http
GET /api/v1/dashboard/stats?period=monthly
Authorization: Bearer <token>
```

**Response:**
```json
{
    "success": true,
    "data": {
        "total_sales": 5250000,
        "total_transactions": 142,
        "total_expenses": 2100000,
        "profit": 3150000,
        "profit_margin": 60.0,
        "top_selling_items": [
            {
                "menu_item": {
                    "id": 1,
                    "name": "Espresso",
                    "price": 15000
                },
                "total_quantity": 45,
                "total_revenue": 675000
            }
        ],
        "sales_chart": [
            {
                "date": "2024-01-01",
                "sales": 125000,
                "transactions": 8
            }
        ],
        "expense_breakdown": [
            {
                "category": "Raw Materials",
                "amount": 1200000,
                "percentage": 57.14
            }
        ]
    }
}
```

## Payment Methods

### Get Payment Methods
```http
GET /api/v1/payment-methods
```

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "name": "Cash",
            "code": "cash",
            "is_active": true
        },
        {
            "id": 2,
            "name": "Credit Card",
            "code": "card",
            "is_active": true
        },
        {
            "id": 3,
            "name": "Digital Wallet",
            "code": "digital_wallet",
            "is_active": true
        }
    ]
}
```

## Error Responses

All endpoints return errors in the following format:

```json
{
    "success": false,
    "message": "Error description",
    "error": "Detailed error information"
}
```

### Common HTTP Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `422` - Validation Error
- `500` - Internal Server Error

## Rate Limiting

The API implements basic rate limiting to prevent abuse. If you exceed the rate limit, you'll receive a `429 Too Many Requests` response.

## CORS

The API supports Cross-Origin Resource Sharing (CORS) for web applications. All origins are allowed in development mode.

## Testing

You can test the API using tools like:
- **Postman**: Import the collection from `docs/postman/`
- **cURL**: Use the examples above
- **Insomnia**: REST client alternative to Postman
