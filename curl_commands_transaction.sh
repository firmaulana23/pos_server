#!/bin/bash

# Individual Curl Commands for Transaction Management
# Copy and paste these commands in your terminal after setting the variables

# Set these variables first:
export BASE_URL="http://localhost:8080/api/v1"
export TOKEN="your_jwt_token_here"

echo "=== Transaction Management Curl Commands ==="
echo ""

echo "1. UPDATE TRANSACTION (Basic Info)"
echo "curl -X PUT \\"
echo "  \${BASE_URL}/transactions/{transaction_id} \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer \${TOKEN}\" \\"
echo "  -d '{"
echo "    \"customer_name\": \"Updated Customer\","
echo "    \"tax\": 3000,"
echo "    \"discount\": 500"
echo "  }'"
echo ""

echo "2. ADD ITEM TO TRANSACTION"
echo "curl -X POST \\"
echo "  \${BASE_URL}/transactions/{transaction_id}/items \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer \${TOKEN}\" \\"
echo "  -d '{"
echo "    \"menu_item_id\": 3,"
echo "    \"quantity\": 2,"
echo "    \"add_ons\": ["
echo "      {"
echo "        \"add_on_id\": 1,"
echo "        \"quantity\": 1"
echo "      },"
echo "      {"
echo "        \"add_on_id\": 2,"
echo "        \"quantity\": 1"
echo "      }"
echo "    ]"
echo "  }'"
echo ""

echo "3. UPDATE TRANSACTION ITEM"
echo "curl -X PUT \\"
echo "  \${BASE_URL}/transactions/{transaction_id}/items/{item_id} \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer \${TOKEN}\" \\"
echo "  -d '{"
echo "    \"quantity\": 3,"
echo "    \"add_ons\": ["
echo "      {"
echo "        \"add_on_id\": 1,"
echo "        \"quantity\": 2"
echo "      }"
echo "    ]"
echo "  }'"
echo ""

echo "4. DELETE TRANSACTION ITEM"
echo "curl -X DELETE \\"
echo "  \${BASE_URL}/transactions/{transaction_id}/items/{item_id} \\"
echo "  -H \"Authorization: Bearer \${TOKEN}\""
echo ""

echo "5. DELETE ENTIRE TRANSACTION"
echo "curl -X DELETE \\"
echo "  \${BASE_URL}/transactions/{transaction_id} \\"
echo "  -H \"Authorization: Bearer \${TOKEN}\""
echo ""

echo "6. GET TRANSACTION DETAILS"
echo "curl -X GET \\"
echo "  \${BASE_URL}/transactions/{transaction_id} \\"
echo "  -H \"Authorization: Bearer \${TOKEN}\""
echo ""

echo "7. GET ALL TRANSACTIONS"
echo "curl -X GET \\"
echo "  \${BASE_URL}/transactions \\"
echo "  -H \"Authorization: Bearer \${TOKEN}\""
echo ""

echo "8. CREATE NEW TRANSACTION"
echo "curl -X POST \\"
echo "  \${BASE_URL}/transactions \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer \${TOKEN}\" \\"
echo "  -d '{"
echo "    \"customer_name\": \"John Doe\","
echo "    \"items\": ["
echo "      {"
echo "        \"menu_item_id\": 1,"
echo "        \"quantity\": 2,"
echo "        \"add_ons\": ["
echo "          {"
echo "            \"add_on_id\": 1,"
echo "            \"quantity\": 1"
echo "          }"
echo "        ]"
echo "      }"
echo "    ],"
echo "    \"tax\": 2500,"
echo "    \"discount\": 0"
echo "  }'"
echo ""

echo "9. PAY TRANSACTION"
echo "curl -X PUT \\"
echo "  \${BASE_URL}/transactions/{transaction_id}/pay \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer \${TOKEN}\" \\"
echo "  -d '{"
echo "    \"payment_method\": \"cash\""
echo "  }'"
echo ""

echo "=== Authentication (Run this first) ==="
echo "curl -X POST \\"
echo "  \${BASE_URL}/auth/login \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{"
echo "    \"username\": \"admin\","
echo "    \"password\": \"admin123\""
echo "  }'"
echo ""

echo "Extract token from response and set:"
echo "export TOKEN=\"your_actual_token_here\""
