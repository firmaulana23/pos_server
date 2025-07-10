#!/bin/bash

# Test script for POS system with reverse proxy (/pos/ prefix)
echo "=== Testing POS System with Reverse Proxy (/pos/) ==="

# Configuration
BASE_URL="http://localhost:8080/pos/api/v1"
USERNAME="admin"
PASSWORD="admin123"

echo "Base URL: $BASE_URL"
echo ""

# Test 1: Authentication
echo ">>> 1. Testing Authentication"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

echo "Login Response: $LOGIN_RESPONSE"

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ] || [ "$TOKEN" = "$LOGIN_RESPONSE" ]; then
    echo "❌ Authentication failed"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ Successfully authenticated"
echo "Token: ${TOKEN:0:20}..."
echo ""

# Test 2: Menu Items
echo ">>> 2. Testing Menu Items API"
MENU_RESPONSE=$(curl -s -X GET "$BASE_URL/menu/items" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

if echo "$MENU_RESPONSE" | grep -q "error"; then
    echo "❌ Menu items API failed"
    echo "Response: $MENU_RESPONSE"
else
    echo "✅ Menu items API working"
    echo "Menu items found: $(echo "$MENU_RESPONSE" | grep -o '"id":[0-9]*' | wc -l)"
fi
echo ""

# Test 3: Dashboard Stats
echo ">>> 3. Testing Dashboard Stats API"
DASHBOARD_RESPONSE=$(curl -s -X GET "$BASE_URL/dashboard/stats" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

if echo "$DASHBOARD_RESPONSE" | grep -q "error"; then
    echo "❌ Dashboard stats API failed"
    echo "Response: $DASHBOARD_RESPONSE"
else
    echo "✅ Dashboard stats API working"
    echo "Total Sales: $(echo "$DASHBOARD_RESPONSE" | grep -o '"total_sales":[0-9.]*' | cut -d':' -f2)"
fi
echo ""

# Test 4: Public API endpoints
echo ">>> 4. Testing Public API endpoints"
PUBLIC_MENU_RESPONSE=$(curl -s -X GET "$BASE_URL/public/menu/items" \
  -H "Content-Type: application/json")

if echo "$PUBLIC_MENU_RESPONSE" | grep -q "error"; then
    echo "❌ Public menu API failed"
    echo "Response: $PUBLIC_MENU_RESPONSE"
else
    echo "✅ Public menu API working"
fi
echo ""

# Test 5: Static Files (admin dashboard)
echo ">>> 5. Testing Static Files Access"
STATIC_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8080/pos/static/css/style.css")

if [ "$STATIC_RESPONSE" = "200" ]; then
    echo "✅ Static files accessible"
else
    echo "❌ Static files not accessible (HTTP $STATIC_RESPONSE)"
fi
echo ""

# Test 6: Admin Dashboard Pages
echo ">>> 6. Testing Admin Dashboard Pages"
ADMIN_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8080/pos/admin/")

if [ "$ADMIN_RESPONSE" = "200" ]; then
    echo "✅ Admin dashboard accessible"
else
    echo "❌ Admin dashboard not accessible (HTTP $ADMIN_RESPONSE)"
fi
echo ""

echo "=== Reverse Proxy Test Complete ==="
echo ""
echo "Summary of endpoints:"
echo "- API Base: $BASE_URL"
echo "- Admin Dashboard: http://localhost:8080/pos/admin/"
echo "- Static Files: http://localhost:8080/pos/static/"
echo "- POS Interface: http://localhost:8080/pos/admin/pos"
echo ""
echo "Sample curl commands:"
echo "# Login"
echo "curl -X POST '$BASE_URL/auth/login' -H 'Content-Type: application/json' -d '{\"username\":\"admin\",\"password\":\"admin123\"}'"
echo ""
echo "# Get Dashboard Stats"
echo "curl -X GET '$BASE_URL/dashboard/stats' -H 'Authorization: Bearer YOUR_TOKEN'"
echo ""
echo "# Get Transactions"
echo "curl -X GET '$BASE_URL/transactions' -H 'Authorization: Bearer YOUR_TOKEN'"
