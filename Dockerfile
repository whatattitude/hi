# 使用一个基础的Go镜像
FROM golang:1.20 AS builder

# 设置工作目录
WORKDIR /app

# 将本地的Go代码复制到镜像中
COPY ./* .

# 编译Go程序
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# 使用一个轻量级的基础Alpine Linux镜像作为运行时镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段的镜像中复制编译好的可执行文件
COPY --from=builder /app/ .

# 添加读写执行权限到 /app/control 文件
RUN chmod +rwx /app/control.sh

# 暴露8081端口
EXPOSE 8081

# 启动应用程序
CMD ["./app"]
