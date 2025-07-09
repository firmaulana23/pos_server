#!/bin/bash

# Simple Transaction Management Test (No jq required)
# Tests the new transaction management endpoints

BASE_URL="http://localhost:8080/api/v1"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Simple Transaction Management Test  ${NC}"
echo -e "${BLUE}======================================${NC}"

# Function to extract ID from JSON response
extract_id() {
    echo "$1" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2
}

# Function to check if response contains error
has_error() {
    echo "$1" | grep -q '"error"'
}

# Get token
echo -e "${YELLOW}Step 1: Authentication${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username": "admin", "password": "admin123"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ Failed to authenticate${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Authenticated successfully${NC}"
AUTH_HEADER="Authorization: Bearer $TOKEN"

# Create a test transaction
echo -e "\n${YELLOW}Step 2: Create Transaction${NC}"
CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/transactions" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "customer_name": "Test Customer",
        "items": [
            {
                "menu_item_id": 15,
                "quantity": 1
            }
        ],
        "tax": 1000,
        "discount": 0
    }')

if has_error "$CREATE_RESPONSE"; then
    echo -e "${RED}❌ Failed to create transaction${NC}"
    echo "Response: $CREATE_RESPONSE"
    exit 1
fi

TRANSACTION_ID=$(extract_id "$CREATE_RESPONSE")
echo -e "${GREEN}✅ Transaction created (ID: $TRANSACTION_ID)${NC}"

# Test 1: Update Transaction Basic Info
echo -e "\n${YELLOW}Step 3: Update Transaction Info${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "customer_name": "Updated Customer",
        "tax": 1500,
        "discount": 200
    }')

if has_error "$UPDATE_RESPONSE"; then
    echo -e "${RED}❌ Failed to update transaction${NC}"
    echo "Response: $UPDATE_RESPONSE"
else
    echo -e "${GREEN}✅ Transaction updated successfully${NC}"
fi

# Test 2: Add Item to Transaction
echo -e "\n${YELLOW}Step 4: Add Item to Transaction${NC}"
ADD_RESPONSE=$(curl -s -X POST "${BASE_URL}/transactions/${TRANSACTION_ID}/items" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "menu_item_id": 17,
        "quantity": 2,
        "add_ons": [
            {
                "add_on_id": 72,
                "quantity": 1
            }
        ]
    }')

if has_error "$ADD_RESPONSE"; then
    echo -e "${RED}❌ Failed to add item${NC}"
    echo "Response: $ADD_RESPONSE"
else
    ITEM_ID=$(extract_id "$ADD_RESPONSE")
    echo -e "${GREEN}✅ Item added successfully (ID: $ITEM_ID)${NC}"
fi

# Test 3: Update Transaction Item
if [ ! -z "$ITEM_ID" ]; then
    echo -e "\n${YELLOW}Step 5: Update Transaction Item${NC}"
    UPDATE_ITEM_RESPONSE=$(curl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}/items/${ITEM_ID}" \
        -H "Content-Type: application/json" \
        -H "$AUTH_HEADER" \
        -d '{
            "quantity": 3,
            "add_ons": [
                {
                    "add_on_id": 72,
                    "quantity": 2
                }
            ]
        }')

    if has_error "$UPDATE_ITEM_RESPONSE"; then
        echo -e "${RED}❌ Failed to update item${NC}"
        echo "Response: $UPDATE_ITEM_RESPONSE"
    else
        echo -e "${GREEN}✅ Item updated successfully${NC}"
    fi
fi

# Test 4: Get Transaction Details
echo -e "\n${YELLOW}Step 6: Get Transaction Details${NC}"
GET_RESPONSE=$(curl -s -X GET "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "$AUTH_HEADER")

if has_error "$GET_RESPONSE"; then
    echo -e "${RED}❌ Failed to get transaction${NC}"
    echo "Response: $GET_RESPONSE"
else
    echo -e "${GREEN}✅ Transaction retrieved successfully${NC}"
    echo "Transaction Data:"
    echo "$GET_RESPONSE"
fi

# Test 5: Delete Transaction Item
if [ ! -z "$ITEM_ID" ]; then
    echo -e "\n${YELLOW}Step 7: Delete Transaction Item${NC}"
    DELETE_ITEM_RESPONSE=$(curl -s -X DELETE "${BASE_URL}/transactions/${TRANSACTION_ID}/items/${ITEM_ID}" \
        -H "$AUTH_HEADER")

    if has_error "$DELETE_ITEM_RESPONSE"; then
        echo -e "${RED}❌ Failed to delete item${NC}"
        echo "Response: $DELETE_ITEM_RESPONSE"
    else
        echo -e "${GREEN}✅ Item deleted successfully${NC}"
    fi
fi

# Test 6: Error Scenarios
echo -e "\n${YELLOW}Step 8: Error Scenarios${NC}"

echo -e "${BLUE}8.1 Update non-existent transaction:${NC}"
ERROR_RESPONSE=$(curl -s -X PUT "${BASE_URL}/transactions/99999" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"customer_name": "Test"}')

if has_error "$ERROR_RESPONSE"; then
    echo -e "${GREEN}✅ Correctly returned error for non-existent transaction${NC}"
else
    echo -e "${RED}❌ Should have returned error${NC}"
fi

echo -e "${BLUE}8.2 Add item with invalid menu_item_id:${NC}"
ERROR_RESPONSE2=$(curl -s -X POST "${BASE_URL}/transactions/${TRANSACTION_ID}/items" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"menu_item_id": 99999, "quantity": 1}')

if has_error "$ERROR_RESPONSE2"; then
    echo -e "${GREEN}✅ Correctly returned error for invalid menu item${NC}"
else
    echo -e "${RED}❌ Should have returned error${NC}"
fi

# Test 7: Pay Transaction and Test Restrictions
echo -e "\n${YELLOW}Step 9: Pay Transaction and Test Restrictions${NC}"
PAY_RESPONSE=$(curl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}/pay" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"payment_method": "cash"}')

if has_error "$PAY_RESPONSE"; then
    echo -e "${RED}❌ Failed to pay transaction${NC}"
    echo "Response: $PAY_RESPONSE"
else
    echo -e "${GREEN}✅ Transaction paid successfully${NC}"
    
    # Try to modify paid transaction
    echo -e "${BLUE}9.1 Try to modify paid transaction:${NC}"
    MODIFY_PAID_RESPONSE=$(curl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}" \
        -H "Content-Type: application/json" \
        -H "$AUTH_HEADER" \
        -d '{"customer_name": "Should Fail"}')
    
    if has_error "$MODIFY_PAID_RESPONSE"; then
        echo -e "${GREEN}✅ Correctly blocked modification of paid transaction${NC}"
    else
        echo -e "${RED}❌ Should have blocked modification${NC}"
    fi
fi

# Test 8: Delete Transaction
echo -e "\n${YELLOW}Step 10: Delete Transaction${NC}"
DELETE_RESPONSE=$(curl -s -X DELETE "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "$AUTH_HEADER")

if has_error "$DELETE_RESPONSE"; then
    echo -e "${RED}❌ Failed to delete transaction${NC}"
    echo "Response: $DELETE_RESPONSE"
else
    echo -e "${GREEN}✅ Transaction deleted successfully${NC}"
fi

echo -e "\n${GREEN}🎉 All tests completed!${NC}"
echo -e "${BLUE}======================================${NC}"

# Summary
echo -e "\n${YELLOW}Summary of tested endpoints:${NC}"
echo "✅ PUT /transactions/{id} - Update transaction"
echo "✅ POST /transactions/{id}/items - Add item"
echo "✅ PUT /transactions/{id}/items/{item_id} - Update item"
echo "✅ DELETE /transactions/{id}/items/{item_id} - Delete item"
echo "✅ DELETE /transactions/{id} - Delete transaction"
echo "✅ Error handling validation"
echo "✅ Paid transaction restrictions"
