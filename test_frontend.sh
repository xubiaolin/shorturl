#!/bin/bash

echo "=========================================="
echo "短链管理系统前端测试"
echo "=========================================="
echo ""

BASE_URL="http://localhost:8080"

echo "1. 测试登录页面访问"
echo "GET /login"
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" "$BASE_URL/login"
echo ""

echo "2. 测试修改密码页面访问"
echo "GET /change-password"
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" "$BASE_URL/change-password"
echo ""

echo "3. 测试管理后台页面访问"
echo "GET /dashboard"
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" "$BASE_URL/dashboard"
echo ""

echo "4. 测试首页重定向到登录页"
echo "GET /"
curl -s -o /dev/null -w "HTTP Status: %{http_code}, Redirect: %{redirect_url}\n" "$BASE_URL/"
echo ""

echo "5. 测试登录 API"
echo "POST /api/v1/login"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}')
echo "$RESPONSE"
echo ""

echo "6. 测试修改密码 API"
echo "POST /api/v1/change-password"
curl -s -X POST "$BASE_URL/api/v1/change-password" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer shorturl-secret-token-2026" \
  -d '{"old_password":"password","new_password":"newpassword123"}'
echo ""
echo ""

echo "7. 使用新密码登录"
echo "POST /api/v1/login"
curl -s -X POST "$BASE_URL/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"newpassword123"}'
echo ""
echo ""

echo "8. 恢复初始密码"
echo "POST /api/v1/change-password"
curl -s -X POST "$BASE_URL/api/v1/change-password" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer shorturl-secret-token-2026" \
  -d '{"old_password":"newpassword123","new_password":"password"}'
echo ""
echo ""

echo "=========================================="
echo "测试完成!"
echo "=========================================="
echo ""
echo "访问地址："
echo "  登录页面：http://localhost:8080/login"
echo "  管理后台：http://localhost:8080/dashboard"
echo ""
echo "默认账户："
echo "  用户名：admin"
echo "  密码：password"
