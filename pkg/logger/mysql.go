package logger

import (
	"fmt"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"go.uber.org/zap"
)

// MysqlLogger 是一个自定义的日志记录器结构体，用于记录 MySQL 相关的日志信息。
// 该结构体实现了特定的日志记录接口，以便在与 MySQL 交互过程中可以按照自定义的方式输出日志。
type MysqlLogger struct{}

// Printf 是 MysqlLogger 结构体实现的一个方法，用于格式化并记录日志信息。
// 此方法遵循了某些日志记录库中常见的 Printf 风格，接收一个格式化模板字符串和可变数量的参数。
// 参数:
// - template: 格式化模板字符串，用于指定日志信息的输出格式。
// - args: 可变数量的参数，这些参数将按照 template 中的格式进行替换。
func (l *MysqlLogger) Printf(template string, args ...interface{}) {
	// 使用 control.info 函数记录日志信息。
	// 首先，通过 fmt.Sprintf 函数将 template 和 args 进行格式化，生成完整的日志消息。
	// 然后，使用 zap.String 函数添加一个键值对到日志中，键为 constants.SourceKey，值为 constants.MysqlSource，
	// 用于标识日志的来源为 MySQL。
	control.info(fmt.Sprintf(template, args...), zap.String(constants.SourceKey, constants.MysqlSource))
}

// GetMysqlLogger 是一个工厂函数，用于创建并返回一个 MysqlLogger 实例的指针。
// 当需要使用 MysqlLogger 进行日志记录时，可以调用此函数获取实例。
// 返回值:
// - 一个指向 MysqlLogger 实例的指针。
func GetMysqlLogger() *MysqlLogger {
	// 创建一个新的 MysqlLogger 实例，并返回其指针。
	return &MysqlLogger{}
}
