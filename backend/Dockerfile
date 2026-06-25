# ==============================
# 第一阶段：编译阶段 (Builder)
# ==============================
FROM golang:1.21-alpine AS builder

# 设置构建环境变量，开启 Goproxy 国内代理极速下载依赖
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 1. 优先拷贝 go.mod 和 go.sum 缓存依赖层
COPY go.mod go.sum ./
RUN go mod download

# 2. 拷贝所有源代码文件
COPY . .

# 3. 静态编译成不依赖操作系统的二进制可执行文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o catering_api ./cmd/main.go


# ==============================
# 第二阶段：运行阶段 (Production)
# ==============================
# 使用仅 5MB 的极其轻量的 alpine 作为底包
FROM alpine:latest

# 安装 tzdata 以保证日志和数据库保存的时间为正确的中国时间
RUN apk --no-cache add tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建镜像中将编译好的二进制文件复制过来
COPY --from=builder /app/catering_api .
# （可选）将环境配置文件一并带入，尽管在 Docker Compose 中我们推荐用环境变量覆盖
COPY --from=builder /app/.env . 

# 声明内部暴露 8080 端口
EXPOSE 8080

# 启动后端程序
CMD ["./catering_api"]
