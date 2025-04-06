#! /usr/bin/env bash

# 获取当前工作目录
CURDIR=$(pwd)

# 此处只涉及 Kitex，但是 Hertz 使用这个没有影响，保留即可
export KITEX_RUNTIME_ROOT=$CURDIR

# 参数替换，检查 ETCD_ADDR 是否已经设置，没有将会设置默认值
: ${ETCD_ADDR:="localhost:2379"}
export ETCD_ADDR

# 这个 SERVICE 环境变量会自动地由 Dockerfile/Makefile 设置
#exec "$CURDIR/output/$SERVICE/yijie-$SERVICE"

# 检查 SERVICE 环境变量是否已设置
if [ -z "$SERVICE" ]; then
    echo "Error: SERVICE environment variable is not set."
    exit 1
fi

# 构建目标路径和可执行文件路径
SERVICE_PATH="$CURDIR/output/$SERVICE"
EXECUTABLE="$SERVICE_PATH/yijie-$SERVICE"

# 检查目标路径是否存在，如果不存在则创建
if [ ! -d "$SERVICE_PATH" ]; then
    echo "Directory $SERVICE_PATH does not exist. Creating it..."
    mkdir -p "$SERVICE_PATH"
fi

# 检查可执行文件是否存在，如果不存在则创建一个空的可执行文件
if [ ! -f "$EXECUTABLE" ]; then
    echo "Executable $EXECUTABLE does not exist. Creating an empty executable file..."
    touch "$EXECUTABLE"
    chmod +x "$EXECUTABLE"
fi

# 检查文件是否可执行
if [ ! -x "$EXECUTABLE" ]; then
    echo "Error: $EXECUTABLE is not executable."
    exit 1
fi

# 执行目标文件
exec "$EXECUTABLE"
