#!/bin/bash

# 设置脚本出错时立即退出
set -e

# 替换或添加 aliyun-odps-go-sdk 依赖为 latest 版本
echo "Setting github.com/aliyun/aliyun-odps-go-sdk to master branch..."

if grep -q "github.com/aliyun/aliyun-odps-go-sdk" go.mod; then
    # 如果存在该依赖，则替换为 latest
    sed -i.bak '/github.com\/aliyun\/aliyun-odps-go-sdk/s@.*@\tgithub.com/aliyun/aliyun-odps-go-sdk latest@' go.mod
else
    # 如果不存在该依赖，则添加
    echo 'require github.com/aliyun/aliyun-odps-go-sdk latest' >> go.mod
fi

# 清理备份文件
rm -f go.mod.bak

# 整理依赖
echo "Running GOPROXY=direct go mod tidy..."
GOPROXY=direct go mod tidy

# 构建程序
echo "Building the binary..."
go build -o mc

# 添加可执行权限
echo "Adding execute permission to mc..."
chmod u+x mc

echo "Build complete! The executable is: ./mc"
