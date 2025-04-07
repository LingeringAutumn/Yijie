package logger

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/LingeringAutumn/Yijie/pkg/constants"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisLogger 是一个用于记录 Redis 操作日志的结构体。
// 这里不加入自定义的 logger 字段，目的是避免在 logger 更新后出现指针引用问题。
type RedisLogger struct{}

// Printf 方法用于格式化并记录 Redis 操作的日志信息。
// 它接收一个上下文、一个格式化模板字符串和可变数量的参数。
// 该方法会将格式化后的日志信息以 Info 级别记录，并标记日志来源为 RedisSource。
func (l *RedisLogger) Printf(ctx context.Context, template string, args ...interface{}) {
	control.info(fmt.Sprintf(template, args...), zap.String(constants.SourceKey, constants.RedisSource))
}

// DialHook 方法是 Redis 的连接钩子函数。
// 它接收一个下一个钩子函数作为参数，并返回一个新的连接钩子函数。
// 该方法只是简单地调用下一个钩子函数，不做额外的日志记录或处理。
func (l *RedisLogger) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// ProcessHook 方法是 Redis 的命令处理钩子函数。
// 它接收一个下一个钩子函数作为参数，并返回一个新的命令处理钩子函数。
// 该方法会记录 Redis 命令的执行时间，并判断是否为慢查询。
func (l *RedisLogger) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		// 记录命令开始执行的时间（毫秒）
		start := time.Now().UnixMilli()

		// 调用下一个钩子函数执行 Redis 命令
		if err := next(ctx, cmd); err != nil {
			return err
		}

		// 计算命令执行的耗时（毫秒）
		consume := time.Now().UnixMilli() - start
		// 如果耗时超过预设的慢查询时间阈值
		if consume >= constants.RedisSlowQuery {
			// 以 Warn 级别记录慢查询日志，包含耗时和查询命令信息，并标记日志来源为 RedisSource
			Warn(fmt.Sprintf("slowly redis query. consume %d microsecond, query: %s", consume, cmd.String()),
				zap.String(constants.SourceKey, constants.RedisSource))
		}

		return nil
	}
}

// ProcessPipelineHook 方法是 Redis 的管道命令处理钩子函数。
// 它接收一个下一个钩子函数作为参数，并返回一个新的管道命令处理钩子函数。
// 该方法只是简单地调用下一个钩子函数，不做额外的日志记录或处理。
func (l *RedisLogger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

// GetRedisLogger 是一个工厂函数，用于创建并返回一个 RedisLogger 实例的指针。
func GetRedisLogger() *RedisLogger {
	return &RedisLogger{}
}
