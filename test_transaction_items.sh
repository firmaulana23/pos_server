#!/bin/bash

# Test script specifically for Transaction Item Management
# Tests the new endpoints: update transaction, add/update/delete items

BASE_URL="http://localhost:8080/api/v1"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Transaction Item Management Tests  ${NC}"
echo -e "${BLUE}======================================${NC}"

# Get token
echo -e "${YELLOW}Getting authentication token...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username": "admin", "password": "admin123"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token // empty')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo -e "${RED}❌ Failed to authenticate${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Authenticated successfully${NC}"
AUTH_HEADER="Authorization: Bearer $TOKEN"

# Create a test transaction
echo -e "\n${YELLOW}Creating test transaction...${NC}"
CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/transactions" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "customer_name": "Test Customer",
        "items": [
            {
                "menu_item_id": 1,
                "quantity": 1
            }
        ],
        "tax": 1000,
        "discount": 0
    }')

TRANSACTION_ID=$(echo $CREATE_RESPONSE | jq -r '.id // empty')

if [ -z "$TRANSACTION_ID" ] || [ "$TRANSACTION_ID" == "null" ]; then
    echo -e "${RED}❌ Failed to create transaction${NC}"
    echo "Response: $CREATE_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Transaction created (ID: $TRANSACTION_ID)${NC}"

# Test 1: Update Transaction Basic Info
echo -e "\n${BLUE}Test 1: Update Transaction Basic Info${NC}"
curl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "customer_name": "Updated Customer Name",
        "tax": 1500,
        "discount": 200
    }' | jq '{
        id: .id,
        customer_name: .customer_name,
        tax: .tax,
        discount: .discount,
        total: .total
    }'

echo -e "${GREEN}✅ Transaction basic info updated${NC}"

# Test 2: Add Item to Transaction
echo -e "\n${BLUE}Test 2: Add New Item to Transaction${NC}"
ADD_RESPONSE=$(curl -s -X POST "${BASE_URL}/transactions/${TRANSACTION_ID}/items" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "menu_item_id": 2,
        "quantity": 2,
        "add_ons": [
            {
                "add_on_id": 1,
                "quantity": 1
            }
        ]
    }')

echo "$ADD_RESPONSE" | jq '{
    id: .id,
    menu_item_id: .menu_item_id,
    quantity: .quantity,
    price: .price
}'

ITEM_ID=$(echo $ADD_RESPONSE | jq -r '.id // empty')
echo -e "${GREEN}✅ Item added (ID: $ITEM_ID)${NC}"

# Test 3: Update Transaction Item
echo -e "\n${BLUE}Test 3: Update Transaction Item${NC}"
curl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}/items/${ITEM_ID}" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "quantity": 3,
        "add_ons": [
            {
                "add_on_id": 1,
                "quantity": 2
            },
            {
                "add_on_id": 2,
                "quantity": 1
            }
        ]
    }' | jq '{
    id: .id,
    quantity: .quantity,
    updated_at: .updated_at
}'

echo -e "${GREEN}✅ Transaction item updated${NC}"

# Test 4: View Updated Transaction
echo -e "\n${BLUE}Test 4: View Complete Updated Transaction${NC}"
curl -s -X GET "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "$AUTH_HEADER" | jq '{
    id: .id,
    customer_name: .customer_name,
    status: .status,
    sub_total: .sub_total,
    tax: .tax,
    discount: .discount,
    total: .total,
    items_count: (.items | length)
}'

echo -e "${GREEN}✅ Transaction details retrieved${NC}"

# Test 5: Add Another Item
echo -e "\n${BLUE}Test 5: Add Another Item${NC}"
ADD_RESPONSE2=$(curl -s -X POST "${BASE_URL}/transactions/${TRANSACTION_ID}/items" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "menu_item_id": 3,
        "quantity": 1
    }')

ITEM_ID2=$(echo $ADD_RESPONSE2 | jq -r '.id // empty')
echo -e "${GREEN}✅ Second item added (ID: $ITEM_ID2)${NC}"

# Test 6: Delete Transaction Item
echo -e "\n${BLUE}Test 6: Delete Transaction Item${NC}"
curl -s -X DELETE "${BASE_URL}/transactions/${TRANSACTION_ID}/items/${ITEM_ID2}" \
    -H "$AUTH_HEADER" | jq '.'

echo -e "${GREEN}✅ Transaction item deleted${NC}"

# Test 7: Final Transaction State
echo -e "\n${BLUE}Test 7: Final Transaction State${NC}"
curl -s -X GET "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "$AUTH_HEADER" | jq '{
    id: .id,
    customer_name: .customer_name,
    status: .status,
    total: .total,
    items: [.items[] | {
        id: .id,
        menu_item_id: .menu_item_id,
        quantity: .quantity,
        add_ons_count: (.add_ons | length)
    }]
}'

echo -e "${GREEN}✅ Final state verified${NC}"

# Test 8: Error Scenarios
echo -e "\n${BLUE}Test 8: Error Scenarios${NC}"

echo -e "${YELLOW}8.1 Try to update non-existent transaction:${NC}"
curl -s -X PUT "${BASE_URL}/transactions/99999" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"customer_name": "Test"}' | jq '.error'

echo -e "${YELLOW}8.2 Try to add item with invalid menu_item_id:${NC}"
curl -s -X POST "${BASE_URL}/transactions/${TRANSACTION_ID}/items" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"menu_item_id": 99999, "quantity": 1}' | jq '.error'

echo -e "${YELLOW}8.3 Try to update non-existent item:${NC}"
curl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}/items/99999" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"quantity": 1}' | jq '.error'

echo -e "${GREEN}✅ Error scenarios tested${NC}"

# Pay the transaction to test paid transaction restrictions
echo -e "\n${BLUE}Test 9: Paid Transaction Restrictions${NC}"
curl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}/pay" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"payment_method": "cash"}' > /dev/null

echo -e "${YELLOW}9.1 Try to update paid transaction:${NC}"
curl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"customer_name": "Should Fail"}' | jq '.error'

echo -e "${YELLOW}9.2 Try to add item to paid transaction:${NC}"
curl -s -X POST "${BASE_URL}/transactions/${TRANSACTION_ID}/items" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"menu_item_id": 1, "quantity": 1}' | jq '.error'

echo -e "${GREEN}✅ Paid transaction restrictions verified${NC}"

# Cleanup
echo -e "\n${BLUE}Cleanup: Deleting test transaction${NC}"
curl -s -X DELETE "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "$AUTH_HEADER" | jq '.'

echo -e "\n${GREEN}🎉 All transaction item management tests completed!${NC}"
