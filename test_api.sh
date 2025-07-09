#!/bin/bash

# API Testing Script for POS System
# Run this script to test the API endpoints including customer name support

BASE_URL="http://localhost:8080/api/v1"
echo "🧪 Testing POS System API"
echo "========================="

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}📋 Test 1: Getting all public add-ons${NC}"
curl -s "$BASE_URL/public/add-ons" | head -c 200
echo -e "\n"

echo -e "${BLUE}📋 Test 2: Getting add-ons for menu item 4 (Latte)${NC}"
curl -s "$BASE_URL/public/menu-item-add-ons/4" | head -c 300
echo -e "\n"

echo -e "${BLUE}📋 Test 3: Getting menu items with their add-ons${NC}"
curl -s "$BASE_URL/public/menu/items?category_id=1" | head -c 400
echo -e "\n"

echo -e "${BLUE}📋 Test 4: Filtering add-ons by menu item${NC}"
curl -s "$BASE_URL/public/add-ons?menu_item_id=4" | head -c 200
echo -e "\n"

echo -e "${BLUE}📋 Test 5: Getting only global add-ons${NC}"
curl -s "$BASE_URL/public/add-ons?menu_item_id=global" | head -c 200
echo -e "\n"

echo -e "${YELLOW}🔐 Authenticated Tests (requires login)${NC}"
echo -e "${BLUE}💡 To test transactions with customer names, first login:${NC}"
echo -e "${BLUE}   curl -X POST $BASE_URL/auth/login -d '{\"username\":\"admin\",\"password\":\"admin123\"}' -H 'Content-Type: application/json'${NC}"
echo -e "${BLUE}💡 Then create a transaction with customer name:${NC}"
echo -e "${BLUE}   curl -X POST $BASE_URL/transactions -H 'Authorization: Bearer <TOKEN>' -H 'Content-Type: application/json' -d '{${NC}"
echo -e "${BLUE}     \"customer_name\": \"John Doe\",${NC}"
echo -e "${BLUE}     \"items\": [{\"menu_item_id\": 1, \"quantity\": 1}],${NC}"
echo -e "${BLUE}     \"payment_method\": \"cash\", \"tax\": 1000, \"discount\": 0${NC}"
echo -e "${BLUE}   }'${NC}"

echo -e "\n${GREEN}✅ API testing completed!${NC}"
