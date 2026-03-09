# 使用Go官方镜像进行编译
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 修改Alpine源为国内镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 设置Go代理
ENV GOPROXY=https://goproxy.cn,direct

# 安装依赖
RUN apk add --no-cache git

# 复制Go模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog-backend-linux .

# 使用Alpine作为运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 修改Alpine源为国内镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 安装必要的依赖
RUN apk add --no-cache ca-certificates tzdata

# 复制编译产物和配置文件
COPY --from=builder /app/blog-backend-linux .
COPY .env .

# 设置执行权限
RUN chmod +x blog-backend-linux

# 创建uploads目录并设置权限
RUN mkdir -p /app/uploads/images && chmod 755 -R /app/uploads

# 暴露端口
EXPOSE 8000

# 启动应用
CMD ["./blog-backend-linux"]