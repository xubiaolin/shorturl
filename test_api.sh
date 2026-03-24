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
curl -s -X GET "$BASE_URL/health" | jq .
echo ""
echo ""

# 2. 测试无 token 访问（应该失败）
echo "2. 测试无 token 访问（应该返回 401）"
echo "POST /api/v1/shorturls (without token)"
curl -s -X POST "$BASE_URL/api/v1/shorturls" \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://www.google.com"}' | jq .
echo ""
echo ""

# 3. 创建短链（带鉴权）
echo "3. 创建短链（带鉴权）"
echo "POST /api/v1/shorturls (with token)"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/shorturls" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url":"https://www.google.com"}')
echo "$RESPONSE" | jq .
SHORT_URL_ID=$(echo "$RESPONSE" | jq -r '.data.id')
SHORT_CODE=$(echo "$RESPONSE" | jq -r '.data.short_code')
echo "Created ID: $SHORT_URL_ID, Short Code: $SHORT_CODE"
echo ""
echo ""

# 4. 创建带自定义短码的短链
echo "4. 创建带自定义短码的短链"
echo "POST /api/v1/shorturls (with custom code)"
curl -s -X POST "$BASE_URL/api/v1/shorturls" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url":"https://www.github.com","custom_code":"github"}' | jq .
echo ""
echo ""

# 5. 获取短链列表
echo "5. 获取短链列表"
echo "GET /api/v1/shorturls"
curl -s -X GET "$BASE_URL/api/v1/shorturls" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""
echo ""

# 6. 获取短链详情
echo "6. 获取短链详情 (ID: $SHORT_URL_ID)"
echo "GET /api/v1/shorturls/$SHORT_URL_ID"
curl -s -X GET "$BASE_URL/api/v1/shorturls/$SHORT_URL_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""
echo ""

# 7. 更新短链
echo "7. 更新短链 (ID: $SHORT_URL_ID)"
echo "PUT /api/v1/shorturls/$SHORT_URL_ID"
curl -s -X PUT "$BASE_URL/api/v1/shorturls/$SHORT_URL_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url":"https://www.google.com","is_active":true}' | jq .
echo ""
echo ""

# 8. 获取短链统计
echo "8. 获取短链统计 (ID: $SHORT_URL_ID)"
echo "GET /api/v1/shorturls/$SHORT_URL_ID/stats"
curl -s -X GET "$BASE_URL/api/v1/shorturls/$SHORT_URL_ID/stats" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""
echo ""

# 9. 测试短链跳转
echo "9. 测试短链跳转（无需鉴权）"
echo "GET /$SHORT_CODE (follow redirect)"
curl -s -L -I -o /dev/null -w "HTTP Status: %{http_code}\nRedirect URL: %{redirect_url}\n" "$BASE_URL/$SHORT_CODE"
echo ""
echo ""

# 10. 删除短链
echo "10. 删除短链 (ID: $SHORT_URL_ID)"
echo "DELETE /api/v1/shorturls/$SHORT_URL_ID"
curl -s -X DELETE "$BASE_URL/api/v1/shorturls/$SHORT_URL_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""
echo ""

# 11. 验证删除后的列表
echo "11. 验证删除后的列表"
echo "GET /api/v1/shorturls"
curl -s -X GET "$BASE_URL/api/v1/shorturls" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""
echo ""

echo "=========================================="
echo "测试完成!"
echo "=========================================="
