#!/bin/bash
# 设置脚本在错误时退出
set -e

# 设置可执行权限
chmod +x /www/wwwroot/blog/backend/blog-backend-linux

# 检查并终止可能存在的相同进程
if pgrep -f "blog-backend-linux" > /dev/null; then
    echo "终止现有进程..."
    pkill -f "blog-backend-linux"
    sleep 2
fi

# 设置工作目录
cd /www/wwwroot/blog/backend

# 启动服务并将输出重定向到日志文件
nohup ./blog-backend-linux > app.log 2>&1 &

echo "服务已启动，PID: $!"