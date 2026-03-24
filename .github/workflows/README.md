# Docker Hub 配置说明

## 前置准备

### 1. 创建 Docker Hub 账户

如果您还没有 Docker Hub 账户，请先访问 [Docker Hub](https://hub.docker.com/) 注册。

### 2. 获取 Docker Hub Token

1. 登录 Docker Hub
2. 点击右上角头像 -> "Account Settings"
3. 选择左侧 "Security" 标签
4. 点击 "New Access Token"
5. 输入描述（如：GitHub Actions）
6. 选择权限：Read & Write
7. 点击 "Generate" 并保存生成的 Token

### 3. 配置 GitHub Secrets

在您的 GitHub 仓库中配置以下 Secrets：

1. 进入仓库 -> Settings -> Secrets and variables -> Actions
2. 点击 "New repository secret"
3. 添加以下两个 secrets：

```
Name: DOCKERHUB_USERNAME
Value: xubiaolin  # 您的 Docker Hub 用户名

Name: DOCKERHUB_TOKEN
Value: <您的 Docker Hub Token>
```

## 工作流程说明

### 触发条件

- 当 `master` 分支有代码推送时自动触发
- 忽略 Markdown 文件和 LICENSE 文件的变更

### 构建特性

- ✅ **多平台支持**: 
  - `linux/amd64` (x86_64)
  - `linux/arm64` (ARM 64-bit)
  - `linux/arm/v7` (ARM v7)
  
- ✅ **自动标签**:
  - `latest` - 最新镜像（仅 master 分支）
  - `<git-sha>` - Git commit 短哈希
  - `<YYYYMMDD>` - 构建日期

- ✅ **构建缓存**: 使用 GitHub Actions 缓存加速构建

### 镜像推送

- PR 时：只构建，不推送
- Master 分支推送：构建并推送到 Docker Hub

## 使用镜像

```bash
# 拉取最新镜像
docker pull xubiaolin/shorturl:latest

# 拉取特定版本
docker pull xubiaolin/shorturl:<git-sha>

# 运行（x86 平台）
docker run -d -p 8080:8080 -v $(pwd)/data:/app/data xubiaolin/shorturl:latest

# 运行（ARM 平台，如 Raspberry Pi）
docker run -d -p 8080:8080 -v $(pwd)/data:/app/data xubiaolin/shorturl:latest
```

## 手动触发构建（可选）

如需手动触发构建，可以修改工作流文件添加 `workflow_dispatch` 触发器：

```yaml
on:
  push:
    branches:
      - master
  workflow_dispatch:  # 添加此行支持手动触发
```

## 查看构建状态

1. 进入 GitHub 仓库
2. 点击 "Actions" 标签
3. 选择 "Build and Push Docker Image" 工作流
4. 查看构建日志和状态

## 故障排查

### 构建失败

检查以下几点：
1. Dockerfile 是否正确
2. GitHub Secrets 是否配置正确
3. Docker Hub 账户是否有推送权限
4. 查看 Actions 日志获取详细错误信息

### 镜像拉取失败

1. 确认镜像名称和标签正确
2. 检查 Docker Hub 账户权限
3. 确认平台架构是否匹配
