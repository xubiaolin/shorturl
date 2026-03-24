# 构建阶段
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache git gcc musl-dev

# 设置 Go 代理（可选，加速国内构建）
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o shorturl .

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖（SQLite 需要）
RUN apk add --no-cache ca-certificates tzdata

# 设置时区为中国时区
ENV TZ=Asia/Shanghai

# 创建应用目录
WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/shorturl .

# 复制静态文件
COPY --from=builder /app/static ./static

# 暴露端口
EXPOSE 8080

# 创建数据目录
RUN mkdir -p /app/data

# 设置环境变量
ENV DATABASE_PATH=/app/data/shorturl.db

# 启动应用
CMD ["./shorturl"]
