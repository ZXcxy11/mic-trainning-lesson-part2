# 使用 Go 官方镜像
FROM golang:latest

# 设置 Go 模块代理（如需要）
ENV GOPROXY="https://goproxy.cn"

# 设置容器内的工作目录
WORKDIR /app

# 复制项目的依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制整个项目
COPY . .

# 交叉编译 `product_srv` 服务
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o product_srv ./product_srv

# 为编译后的可执行文件添加执行权限
RUN chmod +x ./product_srv

# 暴露应用的端口（假设应用使用 8000 端口）
EXPOSE 50051

# 运行 `product_srv` 可执行文件
CMD ["./product_srv/product_srv"]
