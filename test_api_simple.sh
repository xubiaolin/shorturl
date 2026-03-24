#!/bin/bash

# 测试短链服务 API

BASE_URL="http://localhost:8080"
TOKEN="shorturl-secret-token-2026"

echo "=========================================="
echo "短链服务 API 测试"
echo "=========================================="
echo ""

# 1. 测试健康检查接口（无需鉴权）
echo "1. 测试健康检查接口（无需鉴权）"
echo "GET /health"
curl -s -X GET "$BASE_URL/health"
echo ""
echo ""

# 2. 测试无 token 访问（应该失败）
echo "2. 测试无 token 访问（应该返回 401）"
echo "POST /api/v1/shorturls (without token)"
curl -s -X POST "$BASE_URL/api/v1/shorturls" \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://www.google.com"}'
echo ""
echo ""

# 3. 创建短链（带鉴权）
echo "3. 创建短链（带鉴权）"
echo "POST /api/v1/shorturls (with token)"
curl -s -X POST "$BASE_URL/api/v1/shorturls" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url":"https://www.google.com"}'
echo ""
echo ""

# 4. 创建带自定义短码的短链
echo "4. 创建带自定义短码的短链"
echo "POST /api/v1/shorturls (with custom code)"
curl -s -X POST "$BASE_URL/api/v1/shorturls" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url":"https://www.github.com","custom_code":"github"}'
echo ""
echo ""

# 5. 获取短链列表
echo "5. 获取短链列表"
echo "GET /api/v1/shorturls"
curl -s -X GET "$BASE_URL/api/v1/shorturls" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

# 6. 获取短链详情 (ID=1)
echo "6. 获取短链详情 (ID=1)"
echo "GET /api/v1/shorturls/1"
curl -s -X GET "$BASE_URL/api/v1/shorturls/1" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

# 7. 更新短链 (ID=1)
echo "7. 更新短链 (ID=1)"
echo "PUT /api/v1/shorturls/1"
curl -s -X PUT "$BASE_URL/api/v1/shorturls/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url":"https://www.google.com","is_active":true}'
echo ""
echo ""

# 8. 获取短链统计 (ID=1)
echo "8. 获取短链统计 (ID=1)"
echo "GET /api/v1/shorturls/1/stats"
curl -s -X GET "$BASE_URL/api/v1/shorturls/1/stats" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

# 9. 测试短链跳转
echo "9. 测试短链跳转（无需鉴权）"
echo "GET /:code"
curl -s -L -I -o /dev/null -w "HTTP Status: %{http_code}\n" "$BASE_URL/github"
echo ""
echo ""

# 10. 删除短链 (ID=1)
echo "10. 删除短链 (ID=1)"
echo "DELETE /api/v1/shorturls/1"
curl -s -X DELETE "$BASE_URL/api/v1/shorturls/1" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

# 11. 验证删除后的列表
echo "11. 验证删除后的列表"
echo "GET /api/v1/shorturls"
curl -s -X GET "$BASE_URL/api/v1/shorturls" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

echo "=========================================="
echo "测试完成!"
echo "=========================================="
