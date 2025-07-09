# POS System API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Key Features

### Dashboard Analytics
The POS system provides comprehensive analytics including:
- **Financial Metrics**: Total sales, COGS, gross profit, gross margin percentage, net profit
- **Order Statistics**: Total orders, pending orders, paid orders
- **Performance Analytics**: Top-selling menu items and add-ons
- **Visual Charts**: Sales trends and expense breakdowns by date and type

### Menu-Dependent Add-ons
The POS system supports two types of add-ons:

1. **Global Add-ons** (`menu_item_id: null`): Available for all menu items
   - Example: "Whipped Cream", "Extra Hot", "Decaf"
   
2. **Menu-Specific Add-ons** (`menu_item_id: 4`): Only available for specific menu items
   - Example: "Latte Art" (only for Lattes), "Extra Foam" (only for Cappuccinos)

### Cost Management
- **COGS Tracking**: Track Cost of Goods Sold for menu items and add-ons
- **Margin Calculation**: Automatic margin calculation: `((Price - COGS) / Price) * 100`
- **Expense Categories**: Support for raw materials and operational expenses
- **Profit Analysis**: Gross profit (Sales - COGS) and net profit (Gross profit - Expenses)

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
    "username": "admin",
    "password": "admin123"
}
```

**Response:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": 1,
        "username": "admin",
        "email": "admin@pos.com",
        "full_name": "",
        "role": "admin",
        "is_active": true,
        "created_at": "2025-07-07T00:34:37.207182+07:00",
        "updated_at": "2025-07-07T00:34:37.207182+07:00"
    }
}
```

## Public Endpoints (No Authentication Required)

The POS system provides public endpoints for Point of Sale operations that don't require authentication. These are designed to be used by POS terminals.

### Get Categories (Public)
```http
GET /api/v1/public/menu/categories
```

### Get Menu Items (Public)
```http
GET /api/v1/public/menu/items?category_id=1
```

### Get Menu Item (Public)
```http
GET /api/v1/public/menu/items/{id}
```

### Get Add-ons (Public)
```http
GET /api/v1/public/add-ons?available=true
```

### Get Add-ons for Menu Item (Public)
```http
GET /api/v1/public/menu-item-add-ons/{menu_item_id}
```

### Get Payment Methods (Public)
```http
GET /api/v1/public/payment-methods
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
Authorization: Bearer <token>
```

**Response:**
```json
[
    {
        "id": 1,
        "name": "Coffee",
        "description": "Hot and cold coffee beverages",
        "created_at": "2025-07-09T15:07:53.253915+07:00",
        "updated_at": "2025-07-09T15:07:53.253915+07:00",
        "menu_items": [
            {
                "id": 4,
                "category_id": 1,
                "name": "Latte",
                "price": 28000,
                "cogs": 14000,
                "is_available": true
            }
        ]
    }
]
```

**Note:** Categories are also available at the public endpoint `/api/v1/public/menu/categories` for POS display without authentication.

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
GET /api/v1/menu/items?category_id=1&page=1&limit=10
Authorization: Bearer <token>
```

**Query Parameters:**
- `category_id` (optional): Filter by category
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Response:**
```json
{
    "data": [
        {
            "id": 4,
            "category_id": 1,
            "name": "Latte",
            "description": "Espresso with steamed milk",
            "price": 28000,
            "cogs": 14000,
            "margin": 50.0,
            "is_available": true,
            "image_url": "",
            "created_at": "2025-07-09T15:07:53.253915+07:00",
            "updated_at": "2025-07-09T15:07:53.253915+07:00",
            "category": {
                "id": 1,
                "name": "Coffee",
                "description": "Hot and cold coffee beverages"
            },
            "add_ons": [
                {
                    "id": 17,
                    "menu_item_id": 4,
                    "name": "Double Shot for Latte",
                    "description": "Double espresso shot specifically for lattes",
                    "price": 8000,
                    "cogs": 3000,
                    "margin": 62.5,
                    "is_available": true
                },
                {
                    "id": 2,
                    "menu_item_id": null,
                    "name": "Whipped Cream",
                    "description": "Fresh whipped cream",
                    "price": 3000,
                    "cogs": 1500,
                    "margin": 50.0,
                    "is_available": true
                }
            ]
        }
    ]
}
```

**Note:** Menu items now include their associated add-ons (both menu-specific and global add-ons).

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

The system supports both **global add-ons** (available for all menu items) and **menu-specific add-ons** (only available for specific menu items).

### Get All Add-ons
```http
GET /api/v1/add-ons?available=true&menu_item_id=4
```

**Query Parameters:**
- `available` (boolean): Filter by availability status
- `menu_item_id` (integer): Filter by specific menu item ID, or use "global" for global add-ons only

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "menu_item_id": null,
            "name": "Extra Shot",
            "description": "Additional espresso shot",
            "price": 8000,
            "cogs": 4000,
            "margin": 50.0,
            "is_available": true,
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
        },
        {
            "id": 17,
            "menu_item_id": 4,
            "name": "Double Shot for Latte",
            "description": "Double espresso shot specifically for lattes",
            "price": 8000,
            "cogs": 3000,
            "margin": 62.5,
            "is_available": true,
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z",
            "menu_item": {
                "id": 4,
                "name": "Latte",
                "price": 28000
            }
        }
    ]
}
```

### Get Add-ons for Specific Menu Item
```http
GET /api/v1/public/menu-item-add-ons/{menu_item_id}
```

Returns both global add-ons and menu-specific add-ons for the given menu item.

**Response:**
```json
{
    "add_ons": [
        {
            "id": 1,
            "menu_item_id": null,
            "name": "Extra Shot",
            "description": "Additional espresso shot",
            "price": 8000,
            "cogs": 4000,
            "margin": 50.0,
            "is_available": true
        },
        {
            "id": 17,
            "menu_item_id": 4,
            "name": "Double Shot for Latte",
            "description": "Double espresso shot specifically for lattes",
            "price": 8000,
            "cogs": 3000,
            "margin": 62.5,
            "is_available": true
        }
    ],
    "menu_item": {
        "id": 4,
        "name": "Latte"
    }
}
```

### Create Add-on (Admin/Manager)
```http
POST /api/v1/add-ons
Authorization: Bearer <token>
Content-Type: application/json
```

**Global Add-on Example:**
```json
{
    "name": "Oat Milk",
    "description": "Premium oat milk substitute",
    "price": 7000,
    "cogs": 4000,
    "is_available": true
}
```

**Menu-Specific Add-on Example:**
```json
{
    "menu_item_id": 4,
    "name": "Latte Art",
    "description": "Beautiful latte art design (only for lattes)",
    "price": 5000,
    "cogs": 0,
    "is_available": true
}
```

### Update Add-on (Admin/Manager)
```http
PUT /api/v1/add-ons/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
    "menu_item_id": 4,
    "name": "Updated Add-on Name",
    "description": "Updated description",
    "price": 6000,
    "cogs": 3000,
    "is_available": true
}
```

### Delete Add-on (Admin/Manager)
```http
DELETE /api/v1/add-ons/{id}
Authorization: Bearer <token>
```

## Transactions

### Customer Name Support
The POS system now supports storing an optional customer name with each transaction. This field:
- Is optional and can be left empty/null
- Accepts any string value (customer's name)
- Is stored with the transaction for future reference
- Can be used for customer service, receipts, or analytics
- Is included in both transaction creation and retrieval endpoints

### Create Transaction
```http
POST /api/v1/transactions
Authorization: Bearer <token>
Content-Type: application/json

{
    "customer_name": "John Doe",
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

**Request Fields:**
- `customer_name` (string, optional): Customer's name for this transaction
- `items` (array, required): Array of menu items to purchase
- `payment_method` (string, required): Payment method (cash, card, etc.)
- `tax` (number, required): Tax amount in smallest currency unit
- `discount` (number, required): Discount amount in smallest currency unit

**Response:**
```json
{
    "success": true,
    "message": "Transaction created successfully",
    "data": {
        "id": 1,
        "transaction_no": "TXN-20240101-0001",
        "customer_name": "John Doe",
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

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "transaction_no": "TXN-20240101-0001",
            "customer_name": "John Doe",
            "status": "paid",
            "payment_method": "cash",
            "sub_total": 46000,
            "tax": 2500,
            "discount": 0,
            "total": 48500,
            "created_at": "2024-01-01T12:00:00Z",
            "updated_at": "2024-01-01T12:30:00Z",
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
    ],
    "pagination": {
        "current_page": 1,
        "per_page": 10,
        "total": 42,
        "total_pages": 5
    }
}
```

### Delete Transaction (Admin/Manager)
```http
DELETE /api/v1/transactions/:id
Authorization: Bearer <token>
```

**Notes:**
- Pending transactions can be deleted by any authenticated user
- Paid transactions can only be deleted by admin users
- Deletes all related transaction items and add-ons

**Response (Success):**
```json
{
    "message": "Transaction deleted successfully"
}
```

**Response (Error - Insufficient Permissions):**
```json
{
    "error": "Only admin can delete paid transactions"
}
```

**Response (Error - Not Found):**
```json
{
    "error": "Transaction not found"
}
```

### Update Transaction
Update basic transaction information (customer name, tax, discount). Only works on pending transactions.

```http
PUT /api/v1/transactions/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
    "customer_name": "Jane Smith",
    "tax": 3000,
    "discount": 500
}
```

**Response:**
```json
{
    "id": 1,
    "transaction_no": "TXN-20240101-0001",
    "customer_name": "Jane Smith",
    "status": "pending",
    "sub_total": 46000,
    "tax": 3000,
    "discount": 500,
    "total": 48500,
    "updated_at": "2024-01-01T13:00:00Z"
}
```

### Add Item to Transaction
Add a new menu item to a pending transaction.

```http
POST /api/v1/transactions/{id}/items
Authorization: Bearer <token>
Content-Type: application/json

{
    "menu_item_id": 3,
    "quantity": 1,
    "add_ons": [
        {
            "add_on_id": 2,
            "quantity": 1
        }
    ]
}
```

**Response:**
```json
{
    "id": 5,
    "transaction_id": 1,
    "menu_item_id": 3,
    "quantity": 1,
    "price": 18000,
    "created_at": "2024-01-01T13:15:00Z"
}
```

### Update Transaction Item
Update an existing transaction item's quantity and add-ons. Only works on pending transactions.

```http
PUT /api/v1/transactions/{id}/items/{item_id}
Authorization: Bearer <token>
Content-Type: application/json

{
    "quantity": 3,
    "add_ons": [
        {
            "add_on_id": 1,
            "quantity": 2
        },
        {
            "add_on_id": 3,
            "quantity": 1
        }
    ]
}
```

**Response:**
```json
{
    "id": 5,
    "transaction_id": 1,
    "menu_item_id": 3,
    "quantity": 3,
    "price": 18000,
    "updated_at": "2024-01-01T13:20:00Z"
}
```

### Delete Transaction Item
Remove an item from a pending transaction. Automatically recalculates transaction totals.

```http
DELETE /api/v1/transactions/{id}/items/{item_id}
Authorization: Bearer <token>
```

**Response:**
```json
{
    "message": "Transaction item deleted successfully"
}
```

**Common Error Responses for Transaction Item Operations:**
```json
{
    "error": "Cannot modify paid transaction"
}
```

```json
{
    "error": "Transaction not found"
}
```

```json
{
    "error": "Transaction item not found"
}
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

**Access Requirements:** Admin or Manager role required for all dashboard endpoints.

### Get Dashboard Stats
```http
GET /api/v1/dashboard/stats?start_date=2025-01-01&end_date=2025-01-31
Authorization: Bearer <admin_or_manager_token>
```

**Query Parameters:**
- `start_date` (optional): Start date for filtering (YYYY-MM-DD format)
- `end_date` (optional): End date for filtering (YYYY-MM-DD format)

**Note:** If no date parameters are provided, all data will be included in the statistics.

**Response:**
```json
{
    "total_sales": 57700,
    "total_cogs": 28850,
    "gross_profit": 28850,
    "gross_margin_percent": 50.0,
    "total_expenses": 15000,
    "net_profit": 13850,
    "total_orders": 3,
    "pending_orders": 1,
    "paid_orders": 2,
    "top_menu_items": [
        {
            "name": "Latte",
            "total_sold": 2,
            "total_revenue": 56000
        }
    ],
    "top_add_ons": [
        {
            "name": "Double Shot for Latte",
            "total_sold": 1,
            "total_revenue": 8000
        }
    ],
    "sales_chart": [
        {
            "date": "2025-07-09",
            "amount": 57700,
            "orders": 3
        }
    ],
    "expense_chart": [
        {
            "date": "2025-07-09",
            "amount": 15000,
            "type": "Raw Materials"
        }
    ]
}
```

### Get Sales Report
```http
GET /api/v1/dashboard/sales-report?start_date=2025-01-01&end_date=2025-01-31
Authorization: Bearer <admin_token>
```

### Get Profit Analysis
```http
GET /api/v1/dashboard/profit-analysis?start_date=2025-01-01&end_date=2025-01-31
Authorization: Bearer <admin_token>
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

## Database Schema Changes

### Add-ons Table Enhancement
The `add_ons` table now includes:
- `menu_item_id` (nullable): Links add-on to specific menu item
- Foreign key constraint to `menu_items` table
- Index on `menu_item_id` for performance

```sql
-- Migration: 005_add_menu_item_id_to_addons.sql
ALTER TABLE add_ons ADD COLUMN menu_item_id INTEGER REFERENCES menu_items(id);
CREATE INDEX idx_add_ons_menu_item_id ON add_ons(menu_item_id);
```

### Data Model Examples

**Global Add-on:**
```json
{
    "id": 2,
    "menu_item_id": null,
    "name": "Whipped Cream",
    "price": 3000,
    "cogs": 1500
}
```

**Menu-Specific Add-on:**
```json
{
    "id": 17,
    "menu_item_id": 4,
    "name": "Double Shot for Latte",
    "price": 8000,
    "cogs": 3000
}
```

### Migration from Old System
Existing add-ons remain unchanged as global add-ons (`menu_item_id: null`). The system is fully backward compatible:

1. **Existing Global Add-ons**: Continue to work for all menu items
2. **New Menu-Specific Add-ons**: Can be created for targeted offerings
3. **API Compatibility**: Old endpoints continue to work as before
4. **UI Enhancement**: Admin interface now shows add-on types with visual indicators

### Best Practices
- Use **global add-ons** for universal options (milk alternatives, sweeteners, temperature preferences)
- Use **menu-specific add-ons** for specialized options (latte art for lattes, extra foam for cappuccinos)
- Consider customer experience when choosing between global vs. specific add-ons

## Dashboard & Analytics

### Get Dashboard Statistics
```http
GET /api/v1/dashboard/stats?start_date=2025-07-01&end_date=2025-07-31
Authorization: Bearer <admin_or_manager_token>
```

**Query Parameters:**
- `start_date` (optional): Start date filter (YYYY-MM-DD)
- `end_date` (optional): End date filter (YYYY-MM-DD)

**Response:**
```json
{
    "total_sales": 111000,
    "total_cogs": 53300,
    "gross_profit": 57700,
    "gross_margin_percent": 51.98,
    "total_expenses": 24000,
    "net_profit": 33700,
    "total_orders": 2,
    "pending_orders": 0,
    "paid_orders": 2,
    "top_menu_items": [
        {
            "name": "Zona12",
            "total_sold": 4,
            "total_revenue": 70000
        }
    ],
    "top_add_ons": [
        {
            "name": "Telur",
            "total_sold": 2,
            "total_revenue": 10000
        }
    ],
    "sales_chart": [
        {
            "date": "2025-07-09T00:00:00Z",
            "amount": 111000,
            "orders": 2
        }
    ],
    "expense_chart": [
        {
            "date": "2025-07-09T00:00:00Z",
            "amount": 24000,
            "type": "raw_material"
        }
    ]
}
```

**Financial Metrics:**
- `total_sales`: Total revenue from paid transactions
- `total_cogs`: Total Cost of Goods Sold (menu items + add-ons)
- `gross_profit`: Sales - COGS
- `gross_margin_percent`: (Gross Profit / Sales) × 100
- `total_expenses`: Sum of all business expenses (excludes deleted expenses)
- `net_profit`: Gross Profit - Total Expenses

### Get Sales Report
```http
GET /api/v1/dashboard/sales-report?start_date=2025-07-01&end_date=2025-07-31
Authorization: Bearer <admin_or_manager_token>
```

### Get Profit Analysis
```http
GET /api/v1/dashboard/profit-analysis?start_date=2025-07-01&end_date=2025-07-31
Authorization: Bearer <admin_or_manager_token>
```

## Expense Management

### Get Expenses
```http
GET /api/v1/expenses?type=raw_material&page=1&limit=10
Authorization: Bearer <token>
```

### Create Expense (Admin/Manager)
```http
POST /api/v1/expenses
Authorization: Bearer <admin_or_manager_token>
Content-Type: application/json

{
    "type": "raw_material",
    "category": "Coffee Beans",
    "description": "Premium Arabica beans - 5kg",
    "amount": 450000,
    "date": "2025-07-09T10:00:00Z"
}
```

**Expense Types:**
- `raw_material`: Ingredients, coffee beans, milk, etc.
- `operational`: Rent, utilities, staff wages, etc.

### Delete Expense (Admin/Manager)
```http
DELETE /api/v1/expenses/{id}
Authorization: Bearer <admin_or_manager_token>
```

**Note:** Expenses use soft deletion. Deleted expenses are excluded from all financial calculations and dashboard statistics.

## User Management (Admin Only)

### Get Users
```http
GET /api/v1/users
Authorization: Bearer <admin_token>
```

### Update User Role
```http
PUT /api/v1/users/{id}/role
Authorization: Bearer <admin_token>
Content-Type: application/json

{
    "role": "manager"
}
```

**Available Roles:**
- `admin`: Full system access
- `manager`: Can manage menu, transactions, expenses, view analytics
- `cashier`: Can create transactions, view menu

## Error Handling

**Common HTTP Status Codes:**
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

**Error Response Format:**
```json
{
    "error": "Error message description"
}
```

## Testing

You can test the API using tools like:
- **Postman**: Import the collection from `docs/postman/`
- **cURL**: Use the examples above
- **Insomnia**: REST client alternative to Postman

### Testing Menu-Dependent Add-ons

**Test 1: Get add-ons for a specific menu item**
```bash
curl -X GET "http://localhost:8080/api/v1/public/menu-item-add-ons/4"
```

**Test 2: Create a menu-specific add-on**
```bash
curl -X POST "http://localhost:8080/api/v1/add-ons" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "menu_item_id": 4,
    "name": "Extra Foam for Latte",
    "description": "Additional milk foam specifically for lattes",
    "price": 3000,
    "cogs": 1000,
    "is_available": true
  }'
```

**Test 3: Filter add-ons by menu item**
```bash
curl -X GET "http://localhost:8080/api/v1/add-ons?menu_item_id=4"
```

**Test 4: Get only global add-ons**
```bash
curl -X GET "http://localhost:8080/api/v1/add-ons?menu_item_id=global"
```

**Test 5: Get menu items with their add-ons**
```bash
curl -X GET "http://localhost:8080/api/v1/menu/items" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

This will return menu items with both menu-specific and global add-ons included.

## Web Interface Routes

The POS system includes a web-based admin dashboard accessible through the following routes:

### Admin Dashboard Access
- **Login Page**: `GET /admin/` - Authentication page
- **Dashboard**: `GET /admin/dashboard` - Main analytics dashboard
- **Menu Management**: `GET /admin/menu` - Manage categories and menu items
- **Add-ons Management**: `GET /admin/add-ons` - Manage add-ons
- **Transactions**: `GET /admin/transactions` - View and manage transactions
- **Expenses**: `GET /admin/expenses` - Track expenses and costs
- **Point of Sale**: `GET /admin/pos` - POS terminal interface
- **User Management**: `GET /admin/users` - Manage system users (admin only)

### Static Assets
- **Static Files**: `/static/*` - CSS, JavaScript, and other assets
- **Templates**: Rendered HTML templates from `web/templates/`

### Authentication Requirements
- All admin routes require authentication via JWT tokens
- Token is stored in localStorage and sent with API calls
- Dashboard and analytics require admin or manager role
- User management requires admin role only

### Frontend Architecture
- **JavaScript Modules**: Modular JS files for each feature
- **API Integration**: All frontend calls use `/api/v1` base URL
- **Real-time Updates**: Dashboard auto-refreshes and cross-tab synchronization
- **Error Handling**: Comprehensive error display and logging
- **Responsive Design**: Works on desktop and mobile devices

---

For more detailed development information, see the project README and source code documentation.
