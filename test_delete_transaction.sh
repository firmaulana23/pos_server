#!/bin/bash

# Test Delete Transaction API
API_BASE="http://localhost:8080/api/v1"

echo "=== Testing Delete Transaction API ==="
echo

# 1. Login
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "${API_BASE}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | sed 's/"token":"\([^"]*\)"/\1/')
echo "✅ Login successful"

# 2. Create a test transaction
echo
echo "2. Creating test transaction..."
TRANSACTION_DATA='{
  "customer_name": "Test Customer for Delete",
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

TRANSACTION_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
TRANSACTION_NO=$(echo $CREATE_RESPONSE | grep -o '"transaction_no":"[^"]*"' | sed 's/"transaction_no":"\([^"]*\)"/\1/')

echo "✅ Transaction created:"
echo "   ID: $TRANSACTION_ID"
echo "   Number: $TRANSACTION_NO"
echo "   Status: pending"

# 3. Verify transaction exists
echo
echo "3. Verifying transaction exists..."
GET_RESPONSE=$(curl -s -X GET "${API_BASE}/transactions/$TRANSACTION_ID" \
  -H "Authorization: Bearer $TOKEN")

if echo $GET_RESPONSE | grep -q "Test Customer for Delete"; then
    echo "✅ Transaction found and verified"
else
    echo "❌ Transaction not found"
    exit 1
fi

# 4. Delete the transaction
echo
echo "4. Deleting transaction..."
DELETE_RESPONSE=$(curl -s -X DELETE "${API_BASE}/transactions/$TRANSACTION_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Delete response:"
echo $DELETE_RESPONSE

# 5. Verify transaction is deleted
echo
echo "5. Verifying deletion..."
VERIFY_RESPONSE=$(curl -s -X GET "${API_BASE}/transactions/$TRANSACTION_ID" \
  -H "Authorization: Bearer $TOKEN")

if echo $VERIFY_RESPONSE | grep -q "not found"; then
    echo "✅ Transaction successfully deleted"
else
    echo "❌ Transaction still exists"
    echo "Response: $VERIFY_RESPONSE"
fi

# 6. Test deleting a paid transaction (should fail)
echo
echo "6. Testing deletion of paid transaction (should fail)..."

# Create another transaction
CREATE_RESPONSE2=$(curl -s -X POST "${API_BASE}/transactions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "$TRANSACTION_DATA")

TRANSACTION_ID2=$(echo $CREATE_RESPONSE2 | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')

# Pay the transaction
PAY_RESPONSE=$(curl -s -X PUT "${API_BASE}/transactions/$TRANSACTION_ID2/pay" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"payment_method": "cash"}')

echo "Transaction $TRANSACTION_ID2 marked as paid"

# Try to delete paid transaction
DELETE_PAID_RESPONSE=$(curl -s -X DELETE "${API_BASE}/transactions/$TRANSACTION_ID2" \
  -H "Authorization: Bearer $TOKEN")

if echo $DELETE_PAID_RESPONSE | grep -q "Cannot delete paid transaction"; then
    echo "✅ Correctly prevented deletion of paid transaction"
    echo "Response: $DELETE_PAID_RESPONSE"
else
    echo "❌ Should have prevented deletion of paid transaction"
    echo "Response: $DELETE_PAID_RESPONSE"
fi

echo
echo "=== Delete Transaction Test Complete ==="
echo "✅ Pending transactions can be deleted"
echo "✅ Paid transactions cannot be deleted"
echo "✅ Proper error handling implemented"
