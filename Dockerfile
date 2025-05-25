# 构建阶段
FROM golang:1.24-alpine AS builder

# 安装 git（某些依赖需要）
RUN apk add --no-cache git

WORKDIR /app

# 拷贝项目文件
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 构建 Go 可执行文件
RUN go build -o server .

# 运行阶段
FROM alpine:latest

WORKDIR /root/

# 拷贝编译好的二进制
COPY --from=builder /app/server .

# 暴露端口（默认 Gin 在 8080）
EXPOSE 8080

# 设置时区（可选）
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 启动程序
CMD ["./server"]
