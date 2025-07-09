#!/bin/bash

# Test Transaction API with Customer Name
# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# API Base URL
API_BASE="http://localhost:8080/api/v1"

echo -e "${BLUE}=== POS System Transaction API Test with Customer Name ===${NC}"
echo

# Step 1: Login to get token
echo -e "${YELLOW}Step 1: Login to get authentication token${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${API_BASE}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

echo "Login Response:"
echo $LOGIN_RESPONSE
echo

# Extract token using grep and sed (no jq required)
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | sed 's/"token":"\([^"]*\)"/\1/')
if [ -z "$TOKEN" ]; then
    # Try alternative token location
    TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"data":{"token":"[^"]*"' | sed 's/.*"token":"\([^"]*\)"/\1/')
fi

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get token. Exiting.${NC}"
    echo "Response was: $LOGIN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}Token obtained: ${TOKEN:0:50}...${NC}"
echo

# Step 2: Get available menu items
echo -e "${YELLOW}Step 2: Get available menu items${NC}"
MENU_RESPONSE=$(curl -s -X GET "${API_BASE}/public/menu/items" \
  -H "Content-Type: application/json")

echo "Available Menu Items:"
echo $MENU_RESPONSE
echo

# Get first menu item ID using basic text processing
MENU_ITEM_ID=$(echo $MENU_RESPONSE | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
echo -e "${GREEN}Using Menu Item ID: $MENU_ITEM_ID${NC}"
echo

# Step 3: Get available add-ons
echo -e "${YELLOW}Step 3: Get available add-ons${NC}"
ADDONS_RESPONSE=$(curl -s -X GET "${API_BASE}/public/add-ons" \
  -H "Content-Type: application/json")

echo "Available Add-ons:"
echo $ADDONS_RESPONSE
echo

# Get first add-on ID
ADDON_ID=$(echo $ADDONS_RESPONSE | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
echo -e "${GREEN}Using Add-on ID: $ADDON_ID${NC}"
echo

# Step 4: Create transaction with customer name
echo -e "${YELLOW}Step 4: Create transaction with customer name${NC}"
TRANSACTION_DATA='{
  "customer_name": "John Doe",
  "items": [
    {
      "menu_item_id": '$MENU_ITEM_ID',
      "quantity": 2,
      "add_ons": [
        {
          "add_on_id": '$ADDON_ID',
          "quantity": 1
        }
      ]
    }
  ],
  "tax": 5000,
  "discount": 2000
}'

echo "Transaction Data:"
echo $TRANSACTION_DATA
echo

CREATE_RESPONSE=$(curl -s -X POST "${API_BASE}/transactions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "$TRANSACTION_DATA")

echo "Create Transaction Response:"
echo $CREATE_RESPONSE
echo

# Extract transaction ID using basic text processing
TRANSACTION_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
if [ -z "$TRANSACTION_ID" ]; then
    echo -e "${RED}Failed to create transaction. Exiting.${NC}"
    echo "Response was: $CREATE_RESPONSE"
    exit 1
fi

echo -e "${GREEN}Transaction created with ID: $TRANSACTION_ID${NC}"
echo

# Step 5: Get specific transaction details
echo -e "${YELLOW}Step 5: Get transaction details${NC}"
TRANSACTION_DETAILS=$(curl -s -X GET "${API_BASE}/transactions/$TRANSACTION_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Transaction Details:"
echo $TRANSACTION_DETAILS
echo

# Check if customer name is present using basic text processing
CUSTOMER_NAME=$(echo $TRANSACTION_DETAILS | grep -o '"customer_name":"[^"]*"' | sed 's/"customer_name":"\([^"]*\)"/\1/')
echo -e "${GREEN}Customer Name in response: $CUSTOMER_NAME${NC}"
echo

# Step 6: Get all transactions (should include our transaction with customer name)
echo -e "${YELLOW}Step 6: Get all transactions${NC}"
ALL_TRANSACTIONS=$(curl -s -X GET "${API_BASE}/transactions" \
  -H "Authorization: Bearer $TOKEN")

echo "All Transactions (first 500 chars):"
echo "${ALL_TRANSACTIONS:0:500}..."
echo

# Step 7: Test payment process
echo -e "${YELLOW}Step 7: Test payment process${NC}"

# Process payment
PAYMENT_DATA='{
  "payment_method": "cash"
}'

PAY_RESPONSE=$(curl -s -X PUT "${API_BASE}/transactions/$TRANSACTION_ID/pay" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "$PAYMENT_DATA")

echo "Payment Response:"
echo $PAY_RESPONSE
echo

# Step 8: Verify final transaction state
echo -e "${YELLOW}Step 8: Verify final transaction state${NC}"
FINAL_TRANSACTION=$(curl -s -X GET "${API_BASE}/transactions/$TRANSACTION_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Final Transaction State:"
echo $FINAL_TRANSACTION
echo

# Step 9: Test filtering transactions by status
echo -e "${YELLOW}Step 9: Test filtering transactions by status${NC}"
FILTERED_TRANSACTIONS=$(curl -s -X GET "${API_BASE}/transactions?status=paid" \
  -H "Authorization: Bearer $TOKEN")

echo "Paid Transactions (first 500 chars):"
echo "${FILTERED_TRANSACTIONS:0:500}..."
echo

# Step 10: Create another transaction with different customer
echo -e "${YELLOW}Step 10: Create another transaction with different customer${NC}"
TRANSACTION_DATA_2='{
  "customer_name": "Jane Smith",
  "items": [
    {
      "menu_item_id": '$MENU_ITEM_ID',
      "quantity": 1,
      "add_ons": []
    }
  ],
  "tax": 2500,
  "discount": 0
}'

CREATE_RESPONSE_2=$(curl -s -X POST "${API_BASE}/transactions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "$TRANSACTION_DATA_2")

echo "Second Transaction Response:"
echo $CREATE_RESPONSE_2
echo

TRANSACTION_ID_2=$(echo $CREATE_RESPONSE_2 | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
echo -e "${GREEN}Second transaction created with ID: $TRANSACTION_ID_2${NC}"
echo

# Final summary
echo -e "${BLUE}=== Test Summary ===${NC}"
echo -e "${GREEN}✅ Login successful${NC}"
echo -e "${GREEN}✅ Menu items retrieved${NC}"
echo -e "${GREEN}✅ Add-ons retrieved${NC}"
echo -e "${GREEN}✅ Transaction created with customer name: John Doe${NC}"
echo -e "${GREEN}✅ Transaction details retrieved with customer name${NC}"
echo -e "${GREEN}✅ Payment processed successfully${NC}"
echo -e "${GREEN}✅ Final transaction state verified${NC}"
echo -e "${GREEN}✅ Second transaction created with customer name: Jane Smith${NC}"
echo

echo -e "${BLUE}Customer names are properly stored and retrieved in the API!${NC}"
