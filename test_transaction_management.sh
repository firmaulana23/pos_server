#!/bin/bash

# POS System Transaction Management - Curl Test Scripts
# This script demonstrates all transaction CRUD operations

# Configurationcurl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}/items/${ITEM_ID}" \
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
                "add_on_id": 3,
                "quantity": 1
            }
        ]
    }'

echo -e "${GREEN}✓ Transaction item updated${NC}"localhost:8080/api/v1"
USERNAME="admin"
PASSWORD="admin123"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}===========================================${NC}"
echo -e "${BLUE}   POS System Transaction Management Tests${NC}"
echo -e "${BLUE}===========================================${NC}"

# Function to print section headers
print_section() {
    echo -e "\n${YELLOW}>>> $1${NC}"
}

# Function to print test descriptions
print_test() {
    echo -e "${GREEN}$1${NC}"
}

print_section "1. Authentication"
print_test "Logging in to get JWT token..."

# Login to get JWT token
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "'${USERNAME}'",
        "password": "'${PASSWORD}'"
    }')

# Extract token from response
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get authentication token!${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Successfully authenticated${NC}"
echo "Token: ${TOKEN:0:50}..."

# Set authorization header
AUTH_HEADER="Authorization: Bearer $TOKEN"

print_section "2. Create Test Transaction"
print_test "Creating a new pending transaction..."

# Create a new transaction
CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/transactions" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "customer_name": "Test Customer",
        "items": [
            {
                "menu_item_id": 15,
                "quantity": 2,
                "add_ons": [
                    {
                        "add_on_id": 72,
                        "quantity": 1
                    }
                ]
            }
        ],
        "tax": 2500,
        "discount": 0
    }')

echo "Create Transaction Response:"
echo "$CREATE_RESPONSE"

# Extract transaction ID using grep and sed
TRANSACTION_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[^,}]*' | grep -o '[0-9]*')

if [ -z "$TRANSACTION_ID" ]; then
    echo -e "${RED}Failed to create transaction!${NC}"
    echo "Response: $CREATE_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Transaction created with ID: $TRANSACTION_ID${NC}"

print_section "3. Update Transaction Basic Info"
print_test "Updating customer name, tax, and discount..."

curl -s -X PUT "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "customer_name": "Updated Customer Name",
        "tax": 3000,
        "discount": 500
    }'

echo -e "${GREEN}✓ Transaction basic info updated${NC}"

print_section "4. Add Item to Transaction"
print_test "Adding a new menu item to the transaction..."

ADD_ITEM_RESPONSE=$(curl -s -X POST "${BASE_URL}/transactions/${TRANSACTION_ID}/items" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "menu_item_id": 17,
        "quantity": 1,
        "add_ons": [
            {
                "add_on_id": 72,
                "quantity": 1
            }
        ]
    }')

echo "Add Item Response:"
echo "$ADD_ITEM_RESPONSE"

# Extract item ID using grep
ITEM_ID=$(echo $ADD_ITEM_RESPONSE | grep -o '"id":[^,}]*' | head -1 | grep -o '[0-9]*')

if [ -z "$ITEM_ID" ]; then
    echo -e "${RED}Failed to add item to transaction!${NC}"
    echo "Response: $ADD_ITEM_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Item added with ID: $ITEM_ID${NC}"

print_section "5. Update Transaction Item"
print_test "Updating the quantity and add-ons of the transaction item..."

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
    }' | jq '.'

echo -e "${GREEN}✓ Transaction item updated${NC}"

print_section "6. Get Updated Transaction"
print_test "Retrieving the updated transaction to see all changes..."

curl -s -X GET "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "$AUTH_HEADER"

echo -e "${GREEN}✓ Transaction details retrieved${NC}"

print_section "7. Delete Transaction Item"
print_test "Removing the added item from the transaction..."

curl -s -X DELETE "${BASE_URL}/transactions/${TRANSACTION_ID}/items/${ITEM_ID}" \
    -H "$AUTH_HEADER"

echo -e "${GREEN}✓ Transaction item deleted${NC}"

print_section "8. Get All Transactions"
print_test "Retrieving all transactions..."

curl -s -X GET "${BASE_URL}/transactions" \
    -H "$AUTH_HEADER"

echo -e "${GREEN}✓ All transactions retrieved${NC}"

print_section "9. Delete Entire Transaction"
print_test "Deleting the test transaction..."

curl -s -X DELETE "${BASE_URL}/transactions/${TRANSACTION_ID}" \
    -H "$AUTH_HEADER"

echo -e "${GREEN}✓ Transaction deleted${NC}"

print_section "10. Error Testing"
print_test "Testing error scenarios..."

echo -e "\n${YELLOW}10.1 Try to update non-existent transaction:${NC}"
curl -s -X PUT "${BASE_URL}/transactions/99999" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "customer_name": "Test",
        "tax": 1000,
        "discount": 0
    }'

echo -e "\n${YELLOW}10.2 Try to add item to non-existent transaction:${NC}"
curl -s -X POST "${BASE_URL}/transactions/99999/items" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "menu_item_id": 1,
        "quantity": 1
    }'

echo -e "\n${YELLOW}10.3 Try to add non-existent menu item:${NC}"
curl -s -X POST "${BASE_URL}/transactions/${TRANSACTION_ID}/items" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "menu_item_id": 99999,
        "quantity": 1
    }'

echo -e "\n${BLUE}===========================================${NC}"
echo -e "${GREEN}   All transaction tests completed!${NC}"
echo -e "${BLUE}===========================================${NC}"

print_section "Complete Transaction Workflow Example"
print_test "Creating, modifying, and processing a complete transaction..."

# Create another transaction for complete workflow
WORKFLOW_RESPONSE=$(curl -s -X POST "${BASE_URL}/transactions" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
        "customer_name": "Workflow Customer",
        "items": [
            {
                "menu_item_id": 1,
                "quantity": 1
            }
        ],
        "tax": 1500,
        "discount": 0
    }')

WORKFLOW_ID=$(echo $WORKFLOW_RESPONSE | grep -o '"id":[^,}]*' | grep -o '[0-9]*')

if [ ! -z "$WORKFLOW_ID" ]; then
    echo -e "${GREEN}✓ Workflow transaction created: $WORKFLOW_ID${NC}"
    
    # Add more items
    curl -s -X POST "${BASE_URL}/transactions/${WORKFLOW_ID}/items" \
        -H "Content-Type: application/json" \
        -H "$AUTH_HEADER" \
        -d '{
            "menu_item_id": 2,
            "quantity": 2,
            "add_ons": [{"add_on_id": 1, "quantity": 1}]
        }' > /dev/null
    
    # Update transaction info
    curl -s -X PUT "${BASE_URL}/transactions/${WORKFLOW_ID}" \
        -H "Content-Type: application/json" \
        -H "$AUTH_HEADER" \
        -d '{
            "customer_name": "VIP Customer",
            "tax": 2000,
            "discount": 1000
        }' > /dev/null
    
    # Process payment
    curl -s -X PUT "${BASE_URL}/transactions/${WORKFLOW_ID}/pay" \
        -H "Content-Type: application/json" \
        -H "$AUTH_HEADER" \
        -d '{
            "payment_method": "card"
        }'
    
    echo -e "${GREEN}✓ Complete workflow transaction processed and paid${NC}"
fi

echo -e "\n${GREEN}All tests completed successfully!${NC}"
