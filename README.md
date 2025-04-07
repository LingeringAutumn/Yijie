依赖安装
```bash
# 安装kitex和hertz
go get github.com/cloudwego/kitex
go get github.com/cloudwego/hertz
go get github.com/apache/thrift
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
go install github.com/cloudwego/thriftgo@latest

# 安装kafka
go get github.com/segmentio/kafka-go

# 安装mysql
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/go-sql-driver/mysql

# 安装redis
go get github.com/redis/go-redis/v9

# 清理依赖
go mod tidy


```

make clean-all权限不足
```bash
sudo chown -R $(whoami):$(whoami) ./docker/data

```