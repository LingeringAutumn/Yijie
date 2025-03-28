# 辅助工具安装列表
# 执行 go install github.com/cloudwego/hertz/cmd/hz@latest
# 执行 go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
# 执行 go install golang.org/x/tools/cmd/goimports@latest
# 执行 go install golang.org/x/vuln/cmd/govulncheck@latest
# 执行 go install mvdan.cc/gofumpt@latest
# 访问 https://golangci-lint.run/welcome/install/ 以查看安装 golangci-lint 的方法

# 默认输出帮助信息
.DEFAULT_GOAL := help
# 检查 tmux 是否存在
#TMUX_EXISTS := $(shell command -v tmux)
# 项目 MODULE 名
MODULE = github.com/LingeringAutumn/Yijie
# 当前架构
ARCH := $(shell uname -m)
PREFIX = "[Makefile]"
# 目录相关
DIR = $(shell pwd)
CMD = $(DIR)/cmd
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/idl
OUTPUT_PATH = $(DIR)/output
API_PATH= $(DIR)/cmd/api

# 服务名
SERVICES := gateway user commodity order cart payment assistant
service = $(word 1, $@)

EnvironmentStartEnv=YIJIE_ENVIRONMENT_STARTED
EnvironmentStartFlag=true
EtcdAddrEnv=ETCD_ADDR
EtcdAddr=127.0.0.1:2379

# 启动必要的环境，比如 etcd、mysql
.PHONY: env-up
env-up:
	@ docker compose -f ./docker/docker-compose.yml up -d

# 停止服务
.PHONY: env-down
env-down:
	@docker-compose -f ./docker/docker-compose.yml down

# 查看容器日志
.PHONY: env-logs
env-logs:
	@docker-compose -f ./docker/docker-compose.yml logs -f

# 基于 idl 生成相关的 go 语言描述文件
.PHONY: kitex-gen-%
kitex-gen-%:
	@ kitex -module "${MODULE}" \
		-thrift no_default_serdes \
		${IDL_PATH}/$*.thrift
	@ go mod tidy

# 生成 Hertz 文件
.PHONY:new-hz-%
new-hz-%:
	hz new -idl ${IDL_PATH}/api/$*.thrift

# 生成基于 Hertz 的脚手架
.PHONY: hz-%
hz-%:
	hz update -idl ${IDL_PATH}/api/$*.thrift

# 清除所有的构建产物
.PHONY: clean
clean:
	@find . -type d -name "output" -exec rm -rf {} + -print

# 清除所有构建产物、compose 环境和它的数据
.PHONY: clean-all
clean-all: clean
	@echo "$(PREFIX) Checking if docker-compose services are running..."
	@docker-compose -f ./docker/docker-compose.yml ps -q | grep '.' && docker-compose -f ./docker/docker-compose.yml down || echo "$(PREFIX) No services are running."
	@echo "$(PREFIX) Removing docker data..."
	rm -rf ./docker/data
