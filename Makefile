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
TMUX_EXISTS := $(shell command -v tmux)
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
SERVICES := gateway user community interface chat team
service = $(word 1, $@)

EnvironmentStartEnv=YIJIE_ENVIRONMENT_STARTED
EnvironmentStartFlag=true
EtcdAddrEnv=ETCD_ADDR
EtcdAddr=127.0.0.1:2379

## --------------------------------------
## 构建与调试
## --------------------------------------
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


# 构建指定对象，构建后在没有给 BUILD_ONLY 参的情况下会自动运行，需要熟悉 tmux 环境
# 用于本地调试
.PHONY: $(SERVICES)
$(SERVICES):
	@if [ -z "$(TMUX_EXISTS)" ]; then \
		echo "$(PREFIX) tmux is not installed. Please install tmux first."; \
		exit 1; \
	fi
	@if [ -z "$$TMUX" ]; then \
		echo "$(PREFIX) you are not in tmux, press ENTER to start tmux environment."; \
		read -r; \
		if tmux has-session -t fzuhelp 2>/dev/null; then \
			echo "$(PREFIX) Tmux session 'fzuhelp' already exists. Attaching to session and running command."; \
			tmux attach-session -t fzuhelp; \
			tmux send-keys -t fzuhelp "make $(service)" C-m; \
		else \
			echo "$(PREFIX) No tmux session found. Creating a new session."; \
			tmux new-session -s fzuhelp "make $(service); $$SHELL"; \
		fi; \
	else \
		echo "$(PREFIX) Build $(service) target..."; \
		mkdir -p output; \
		bash $(DIR)/docker/script/build.sh $(service); \
		echo "$(PREFIX) Build $(service) target completed"; \
	fi
ifndef BUILD_ONLY
	@echo "$(PREFIX) Automatic run server"
	@if tmux list-windows -F '#{window_name}' | grep -q "^yijie-$(service)$$"; then \
		echo "$(PREFIX) Window 'yijie-$(service)' already exists. Reusing the window."; \
		tmux select-window -t "yijie-$(service)"; \
	else \
		echo "$(PREFIX) Window 'yijie-$(service)' does not exist. Creating a new window."; \
		tmux new-window -n "yijie-$(service)"; \
		tmux split-window -h ; \
		tmux select-layout -t "yijie-$(service)" even-horizontal; \
	fi
	@echo "$(PREFIX) Running $(service) service in tmux..."
	@tmux send-keys -t yijie-$(service).0 'export SERVICE=$(service) && bash ./docker/script/entrypoint.sh' C-m
	@tmux select-pane -t yijie-$(service).1
endif


# 格式化代码，我们使用 gofumpt，是 fmt 的严格超集
.PHONY: fmt
fmt:
	gofumpt -l -w .

# 优化 import 顺序结构
.PHONY: import
import:
	goimports -w -local github.com/LingeringAutumn .

# 检查可能的错误
.PHONY: vet
vet:
	go vet ./...

# 代码格式校验
.PHONY: lint
lint:
	golangci-lint run --config=./.golangci.yml --tests --allow-parallel-runners --sort-results --show-stats --print-resources-usage

# 检查依赖漏洞
.PHONY: vulncheck
vulncheck:
	govulncheck ./...

.PHONY: tidy
tidy:
	go mod tidy

# 一键修正规范并执行代码检查，同时运行 license 检查
# TODO lint要加进去
.PHONY: verify
verify:  vet fmt import vulncheck tidy
