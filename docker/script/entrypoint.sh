#! /usr/bin/env bash
CURDIR=$(pwd)

# 此处只涉及 Kitex，但是 Hertz 使用这个没有影响，保留即可
export KITEX_RUNTIME_ROOT=$CURDIR

# 参数替换，检查 ETCD_ADDR 是否已经设置，没有将会设置默认值
: ${ETCD_ADDR:="localhost:2379"}
export ETCD_ADDR

# 这个 SERVICE 环境变量会自动地由 Dockerfile/Makefile 设置
exec "$CURDIR/output/$SERVICE/yijie-$SERVICE"
