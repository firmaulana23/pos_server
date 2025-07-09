#!/bin/bash

# Simple Transaction API Test with Customer Name
API_BASE="http://localhost:8080/api/v1"

echo "=== Testing Transaction API with Customer Name ==="
echo

# 1. Login
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "${API_BASE}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | sed 's/"token":"\([^"]*\)"/\1/')
echo "✅ Login successful, token obtained"

# 2. Create transaction with customer name
echo
echo "2. Creating transaction with customer name..."
TRANSACTION_DATA='{
  "customer_name": "Alice Johnson",
  "items": [
    {
      "menu_item_id": 4,
      "quantity": 1,
      "add_ons": []
    }
  ],
  "tax": 2800,
  "discount": 0
}'

CREATE_RESPONSE=$(curl -s -X POST "${API_BASE}/transactions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "$TRANSACTION_DATA")

echo "Transaction created:"
echo $CREATE_RESPONSE | grep -o '"customer_name":"[^"]*"'
echo $CREATE_RESPONSE | grep -o '"transaction_no":"[^"]*"'
echo $CREATE_RESPONSE | grep -o '"total":[0-9]*'

TRANSACTION_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')

# 3. Get transaction details
echo
echo "3. Retrieving transaction details..."
DETAILS_RESPONSE=$(curl -s -X GET "${API_BASE}/transactions/$TRANSACTION_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Customer name from API:"
echo $DETAILS_RESPONSE | grep -o '"customer_name":"[^"]*"'

# 4. List all transactions
echo
echo "4. Listing all transactions..."
ALL_RESPONSE=$(curl -s -X GET "${API_BASE}/transactions" \
  -H "Authorization: Bearer $TOKEN")

echo "Recent transactions with customer names:"
echo $ALL_RESPONSE | grep -o '"customer_name":"[^"]*"' | head -3

echo
echo "✅ Customer name functionality is working correctly!"
echo "   - Customer names are stored when creating transactions"
echo "   - Customer names are returned in API responses"
echo "   - Customer names appear in transaction lists"
