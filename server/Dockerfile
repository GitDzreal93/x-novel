# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o x-novel-server ./cmd/server

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 复制编译好的二进制文件
COPY --from=builder /app/x-novel-server .

# 复制配置文件
COPY config.yaml .

EXPOSE 8080

CMD ["./x-novel-server"]
