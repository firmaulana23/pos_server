# Transaction Management Features

This document describes the enhanced transaction management capabilities in the POS system.

## New Features

### 1. Update Transaction Basic Information
- **Endpoint**: `PUT /api/v1/transactions/{id}`
- **Purpose**: Update customer name, tax, and discount for pending transactions
- **Restrictions**: Only works on pending transactions
- **Auto-recalculation**: Automatically recalculates totals when tax/discount changes

### 2. Add Items to Existing Transactions
- **Endpoint**: `POST /api/v1/transactions/{id}/items`
- **Purpose**: Add new menu items with add-ons to pending transactions
- **Features**: 
  - Supports add-ons
  - Validates menu item existence
  - Automatically updates transaction totals

### 3. Update Transaction Items
- **Endpoint**: `PUT /api/v1/transactions/{id}/items/{item_id}`
- **Purpose**: Modify quantity and add-ons of existing transaction items
- **Features**:
  - Update item quantity
  - Replace all add-ons (removes old, adds new)
  - Automatically recalculates totals

### 4. Remove Items from Transactions
- **Endpoint**: `DELETE /api/v1/transactions/{id}/items/{item_id}`
- **Purpose**: Remove specific items from pending transactions
- **Features**:
  - Removes item and all associated add-ons
  - Automatically updates transaction totals
  - Maintains data integrity

## Permission Changes

### Previous Behavior
- Only `admin` and `manager` roles could delete transactions
- Other operations had various role restrictions

### Current Behavior
- **All authenticated users** can perform all transaction operations
- **Business rule**: Operations only work on pending transactions
- **Paid transactions**: Cannot be modified (regardless of user role)

## Security & Validation

### Transaction Status Validation
```
✅ Pending transactions: All operations allowed
❌ Paid transactions: No modifications allowed
```

### Data Integrity
- All operations use database transactions
- Automatic rollback on errors
- Consistent total calculations
- Foreign key validation

### Input Validation
- Required fields validation
- Positive quantity enforcement
- Menu item existence verification
- Add-on existence verification

## Testing

### Test Scripts Available
1. **`test_transaction_management.sh`** - Complete workflow tests
2. **`test_transaction_items.sh`** - Focus on item management
3. **`curl_commands_transaction.sh`** - Individual command examples

### Running Tests
```bash
# Make scripts executable
chmod +x test_*.sh curl_*.sh

# Run comprehensive tests
./test_transaction_management.sh

# Run item-focused tests
./test_transaction_items.sh

# View individual curl commands
./curl_commands_transaction.sh
```

## API Examples

### Update Transaction
```bash
curl -X PUT "http://localhost:8080/api/v1/transactions/123" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Doe",
    "tax": 2500,
    "discount": 500
  }'
```

### Add Item
```bash
curl -X POST "http://localhost:8080/api/v1/transactions/123/items" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "menu_item_id": 5,
    "quantity": 2,
    "add_ons": [
      {"add_on_id": 1, "quantity": 1}
    ]
  }'
```

### Update Item
```bash
curl -X PUT "http://localhost:8080/api/v1/transactions/123/items/456" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 3,
    "add_ons": [
      {"add_on_id": 1, "quantity": 2}
    ]
  }'
```

### Delete Item
```bash
curl -X DELETE "http://localhost:8080/api/v1/transactions/123/items/456" \
  -H "Authorization: Bearer $TOKEN"
```

## Use Cases

### Typical Workflow
1. Customer places initial order → Create transaction
2. Customer changes mind → Update/add/remove items
3. Apply discount/adjust tax → Update transaction
4. Customer satisfied → Pay transaction
5. ❌ No more modifications allowed after payment

### Error Handling
- Clear error messages for invalid operations
- Proper HTTP status codes
- Detailed validation feedback
- Safe transaction rollbacks

## Notes

- All monetary values should be in the smallest currency unit (e.g., cents)
- Transaction totals are automatically calculated: `total = subtotal + tax - discount`
- Add-ons are tied to specific transaction items
- Deleting an item removes all its add-ons
- Database constraints ensure referential integrity
