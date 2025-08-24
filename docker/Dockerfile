# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要工具
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./examples/basic/main.go

# 运行阶段
FROM alpine:latest

# 安装ca证书
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/.env.example .env

# 创建日志目录
RUN mkdir -p logs

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]