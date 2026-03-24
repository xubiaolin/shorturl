# 版本发布指南

## 语义化版本规范

本项目遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN/) 规范。

版本格式：`v主版本号。次版本号。修订号`

### 版本号规则

- **主版本号（Major）**: 不兼容的 API 变更
- **次版本号（Minor）**: 向后兼容的新功能
- **修订号（Patch）**: 向后兼容的问题修正

### 示例

```
v1.0.0    # 初始稳定版本
v1.0.1    # Bug 修复
v1.1.0    # 新功能
v2.0.0    # 不兼容变更
v1.0.0-alpha.1  # 预发布版本
```

## 发布流程

### 1. 更新版本号

在发布新版本前，建议更新项目中的版本信息（如有）。

### 2. 打 Tag 并推送

```bash
# 查看当前版本
git tag -l

# 创建新版本 tag（在本地）
git tag -a v1.0.0 -m "Release version 1.0.0"

# 或者创建轻量级 tag
git tag v1.0.0

# 推送到远程
git push origin v1.0.0

# 或者推送所有 tags
git push origin --tags
```

### 3. 自动化流程

推送 tag 后，GitHub Actions 会自动：

1. ✅ 构建多平台 Docker 镜像（x86, ARM 等）
2. ✅ 推送到 Docker Hub，标签为 `v1.0.0`
3. ✅ 创建 GitHub Release
4. ✅ 自动生成更新日志

### 4. 查看发布状态

- **GitHub Actions**: 查看构建状态
- **GitHub Releases**: 查看发布的版本
- **Docker Hub**: 查看推送的镜像

## Docker 镜像标签说明

| 标签类型 | 示例 | 说明 |
|---------|------|------|
| 版本标签 | `v1.0.0` | 精确版本，推荐生产使用 |
| 主版本 | - | 暂不支持 |
| latest | `latest` | 最新稳定版（master 分支） |
| SHA | `abc1234` | Git commit 短哈希 |

## 常用命令

### 查看现有标签

```bash
# 列出所有标签
git tag -l

# 查看特定标签信息
git show v1.0.0

# 查看最近的标签
git describe --tags --abbrev=0
```

### 创建标签

```bash
# 创建轻量级标签
git tag v1.0.0

# 创建带注释的标签（推荐）
git tag -a v1.0.0 -m "Release version 1.0.0"

# 基于特定 commit 创建标签
git tag -a v1.0.0 abc1234 -m "Release version 1.0.0"
```

### 删除标签

```bash
# 删除本地标签
git tag -d v1.0.0

# 删除远程标签
git push origin :refs/tags/v1.0.0

# 删除并重新推送
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
git tag -a v1.0.0 -m "Re-release version 1.0.0"
git push origin v1.0.0
```

## 发布检查清单

发布新版本前，请确认：

- [ ] 所有测试通过
- [ ] 更新日志已准备
- [ ] 文档已更新
- [ ] 代码已合并到 master 分支
- [ ] Docker Hub 凭证已配置

## 预发布版本

如需发布测试版本：

```bash
# Alpha 版本（内部测试）
git tag -a v1.0.0-alpha.1 -m "Alpha release 1"

# Beta 版本（公开测试）
git tag -a v1.0.0-beta.1 -m "Beta release 1"

# Release Candidate（候选版本）
git tag -a v1.0.0-rc.1 -m "RC 1"
```

## 示例发布流程

```bash
# 1. 确保在 master 分支
git checkout master
git pull origin master

# 2. 运行测试
go test ./...

# 3. 创建新版本
git tag -a v1.0.0 -m "Release version 1.0.0 - Initial stable release"

# 4. 推送 tag
git push origin v1.0.0

# 5. 等待 GitHub Actions 完成构建
# 访问：https://github.com/xubiaolin/shorturl/actions

# 6. 验证 Docker 镜像
docker pull xubiaolin/shorturl:v1.0.0
docker run --rm xubiaolin/shorturl:v1.0.0 --version
```

## 故障排查

### 构建失败

1. 检查 GitHub Actions 日志
2. 确认 Docker Hub 凭证正确
3. 验证 Dockerfile 语法

### Release 未创建

1. 检查 `release.yml` 工作流状态
2. 确认 `GITHUB_TOKEN` 权限

### Docker 镜像未推送

1. 检查 tag 格式（必须是 `v*` 格式）
2. 验证 Docker Hub 凭证
3. 查看构建日志
