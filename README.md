# 短链管理系统 (ShortURL Management System)

一个基于 Go 和 Gin 框架构建的轻量级短链生成与管理系统，提供完整的短链创建、管理、跳转统计等功能，并配备现代化的 Web 管理界面。

## 功能特性

- 🔗 **短链生成** - 支持自定义短码和自动过期时间设置
- 📊 **访问统计** - 实时追踪短链点击次数
- 🔐 **用户认证** - JWT Token 鉴权，支持修改密码
- 🎨 **现代化 UI** - 基于 Tailwind CSS 的美观管理界面
- 🛡️ **中间件** - 完善的请求鉴权中间件
- 💾 **数据持久化** - 使用 SQLite 数据库，GORM ORM 管理
- 📱 **响应式设计** - 管理界面支持多端访问

## 技术栈

- **后端框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **数据库**: SQLite
- **前端**: Tailwind CSS + 原生 JavaScript
- **鉴权**: JWT Bearer Token

## 项目结构

```
shorturl/
├── main.go                 # 应用入口和路由配置
├── go.mod                  # Go 模块依赖
├── Dockerfile              # Docker 镜像构建文件
├── docker-compose.yml      # Docker Compose 配置文件
├── .dockerignore           # Docker 构建忽略文件
├── handlers/               # HTTP 请求处理器
│   ├── auth.go            # 认证相关处理
│   ├── auth_handler.go    # 认证处理器
│   ├── redirect.go        # 短链跳转处理
│   └── shorturl.go        # 短链 CRUD 处理
├── middleware/             # 中间件
│   └── auth.go            # JWT 鉴权中间件
├── models/                 # 数据模型
│   ├── database.go        # 数据库初始化
│   ├── shorturl.go        # 短链模型
│   └── user.go            # 用户模型
├── service/                # 业务逻辑层
│   └── shorturl.go        # 短链业务服务
├── static/                 # 静态资源
│   ├── js/
│   │   └── app.js         # 前端 JavaScript
│   ├── login.html         # 登录页面
│   ├── dashboard.html     # 管理后台页面
│   └── change-password.html  # 修改密码页面
├── test_api.sh            # API 测试脚本
└── test_api_simple.sh     # 简化版 API 测试脚本
```

## 快速开始

### 方式一：Docker 部署（推荐）

#### 环境要求

- Docker 20.10+
- Docker Compose 2.0+

#### 使用 Docker Compose 启动

```bash
# 构建并启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

服务启动后，数据将自动挂载在当前目录的 `./data` 文件夹中。

#### 使用 Docker 命令启动

```bash
# 构建镜像
docker build -t shorturl:latest .

# 启动容器
docker run -d \
  --name shorturl-service \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e GIN_MODE=release \
  shorturl:latest
```

### 方式二：本地运行

#### 环境要求

- Go 1.25.0 或更高版本
- SQLite3

#### 安装依赖

```bash
go mod download
```

#### 启动服务

```bash
go run main.go
```

服务默认运行在 `http://localhost:8080`

### 默认账户

- 用户名：`admin`
- 密码：`admin123`

## API 接口

### 公开接口（无需鉴权）

#### 健康检查
```bash
GET /health
```

#### 用户登录
```bash
POST /api/v1/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

#### 短链跳转
```bash
GET /:code
```

### 认证接口（需要 Token）

在请求头中携带 Token：
```
Authorization: Bearer <your_token>
```

#### 创建短链
```bash
POST /api/v1/shorturls
Content-Type: application/json

{
  "original_url": "https://example.com",
  "custom_code": "mylink",  // 可选
  "expires_at": 1234567890   // 可选，Unix 时间戳
}
```

#### 获取短链列表
```bash
GET /api/v1/shorturls
```

#### 获取单个短链
```bash
GET /api/v1/shorturls/:id
```

#### 更新短链
```bash
PUT /api/v1/shorturls/:id
Content-Type: application/json

{
  "original_url": "https://new-url.com",
  "is_active": true
}
```

#### 删除短链
```bash
DELETE /api/v1/shorturls/:id
```

#### 查看短链统计
```bash
GET /api/v1/shorturls/:id/stats
```

#### 修改密码
```bash
POST /api/v1/change-password
Content-Type: application/json

{
  "old_password": "admin123",
  "new_password": "newpassword123"
}
```

## 前端页面

- **登录页面**: `http://localhost:8080/login`
- **管理后台**: `http://localhost:8080/dashboard`
- **修改密码**: `http://localhost:8080/change-password`

## 测试

### 完整 API 测试

```bash
chmod +x test_api.sh
./test_api.sh
```

### 简化版测试

```bash
chmod +x test_api_simple.sh
./test_api_simple.sh
```

## 数据模型

### User（用户）
- `id`: 主键
- `username`: 用户名（唯一，50 字符）
- `password`: 密码（哈希存储）
- `is_active`: 是否激活
- `created_at`, `updated_at`, `deleted_at`: 时间戳

### ShortURL（短链）
- `id`: 主键
- `original_url`: 原始 URL
- `short_code`: 短码（唯一，20 字符）
- `short_url`: 完整短链地址
- `clicks`: 点击次数
- `is_active`: 是否激活
- `expires_at`: 过期时间（可选）
- `created_at`, `updated_at`, `deleted_at`: 时间戳

## 数据持久化

使用 Docker 部署时，数据会自动挂载到 `./data` 目录：

- **数据库文件**: `./data/shorturl.db`
- **数据备份**: 只需备份 `data` 目录即可

## 安全特性

- ✅ 密码哈希存储
- ✅ JWT Token 鉴权
- ✅ 中间件级别的权限控制
- ✅ SQL 注入防护（GORM）
- ✅ 软删除支持

## 开发说明

### 添加新的 API 接口

1. 在 `handlers/` 目录下创建对应的处理器方法
2. 在 `service/` 层实现业务逻辑
3. 在 `main.go` 中注册路由

### 数据库迁移

项目使用 GORM 的 AutoMigrate 功能自动管理数据库 schema，启动时会自动同步模型变更。

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
