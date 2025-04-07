package logger

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudwego/kitex/pkg/klog"
	"go.uber.org/zap"

	"github.com/LingeringAutumn/Yijie/pkg/constants"
)

// KlogLogger 结构体实现了 KiteX 日志记录器接口，用于将 KiteX 的日志输出重定向到自定义的日志记录器。
type KlogLogger struct{}

// GetKlogLogger 是一个工厂函数，用于创建并返回一个 KlogLogger 实例的指针。
func GetKlogLogger() *KlogLogger {
	return &KlogLogger{}
}

// Trace 方法实现了 KiteX 日志记录器的 Trace 级别日志记录功能。
// 它将传入的参数转换为字符串，并调用 control 的 debug 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Trace(v ...interface{}) {
	control.debug(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Debug 方法实现了 KiteX 日志记录器的 Debug 级别日志记录功能。
// 它将传入的参数转换为字符串，并调用 control 的 debug 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Debug(v ...interface{}) {
	control.debug(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Info 方法实现了 KiteX 日志记录器的 Info 级别日志记录功能。
// 它将传入的参数转换为字符串，并调用 control 的 info 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Info(v ...interface{}) {
	control.info(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Notice 方法实现了 KiteX 日志记录器的 Notice 级别日志记录功能。
// 由于 Notice 级别在自定义日志记录器中没有特别区分，这里将其视为 Info 级别，调用 control 的 info 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Notice(v ...interface{}) {
	control.info(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Warn 方法实现了 KiteX 日志记录器的 Warn 级别日志记录功能。
// 它将传入的参数转换为字符串，并调用 control 的 warn 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Warn(v ...interface{}) {
	control.warn(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Error 方法实现了 KiteX 日志记录器的 Error 级别日志记录功能。
// 它将传入的参数转换为字符串，并调用 control 的 error 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Error(v ...interface{}) {
	control.error(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Fatal 方法实现了 KiteX 日志记录器的 Fatal 级别日志记录功能。
// 它将传入的参数转换为字符串，并调用 control 的 fatal 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Fatal(v ...interface{}) {
	control.fatal(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Tracef 方法实现了 KiteX 日志记录器的 Trace 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 debug 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Tracef(format string, v ...interface{}) {
	control.debug(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Debugf 方法实现了 KiteX 日志记录器的 Debug 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 debug 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Debugf(format string, v ...interface{}) {
	control.debug(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Infof 方法实现了 KiteX 日志记录器的 Info 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 info 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Infof(format string, v ...interface{}) {
	control.info(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Noticef 方法实现了 KiteX 日志记录器的 Notice 级别格式化日志记录功能。
// 由于 Notice 级别在自定义日志记录器中没有特别区分，这里将其视为 Info 级别，调用 control 的 info 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Noticef(format string, v ...interface{}) {
	control.info(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Warnf 方法实现了 KiteX 日志记录器的 Warn 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 warn 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Warnf(format string, v ...interface{}) {
	control.warn(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Errorf 方法实现了 KiteX 日志记录器的 Error 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 error 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Errorf(format string, v ...interface{}) {
	control.error(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// Fatalf 方法实现了 KiteX 日志记录器的 Fatal 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 fatal 方法记录日志，同时标记日志来源为 KlogSource。
func (l *KlogLogger) Fatalf(format string, v ...interface{}) {
	control.fatal(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// CtxTracef 方法实现了 KiteX 日志记录器在上下文中的 Trace 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 debug 方法记录日志，同时标记日志来源为 KlogSource。
// 这里上下文参数 ctx 未被使用，可根据实际需求扩展。
func (l *KlogLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	control.debug(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// CtxDebugf 方法实现了 KiteX 日志记录器在上下文中的 Debug 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 debug 方法记录日志，同时标记日志来源为 KlogSource。
// 这里上下文参数 ctx 未被使用，可根据实际需求扩展。
func (l *KlogLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	control.debug(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// CtxInfof 方法实现了 KiteX 日志记录器在上下文中的 Info 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 info 方法记录日志，同时标记日志来源为 KlogSource。
// 这里上下文参数 ctx 未被使用，可根据实际需求扩展。
func (l *KlogLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	control.info(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// CtxNoticef 方法实现了 KiteX 日志记录器在上下文中的 Notice 级别格式化日志记录功能。
// 由于 Notice 级别在自定义日志记录器中没有特别区分，这里将其视为 Info 级别，调用 control 的 info 方法记录日志，同时标记日志来源为 KlogSource。
// 这里上下文参数 ctx 未被使用，可根据实际需求扩展。
func (l *KlogLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	control.info(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// CtxWarnf 方法实现了 KiteX 日志记录器在上下文中的 Warn 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 warn 方法记录日志，同时标记日志来源为 KlogSource。
// 这里上下文参数 ctx 未被使用，可根据实际需求扩展。
func (l *KlogLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	control.warn(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// CtxErrorf 方法实现了 KiteX 日志记录器在上下文中的 Error 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 error 方法记录日志，同时标记日志来源为 KlogSource。
// 这里上下文参数 ctx 未被使用，可根据实际需求扩展。
func (l *KlogLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	control.error(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// CtxFatalf 方法实现了 KiteX 日志记录器在上下文中的 Fatal 级别格式化日志记录功能。
// 它将传入的格式化字符串和参数进行格式化，并调用 control 的 fatal 方法记录日志，同时标记日志来源为 KlogSource。
// 这里上下文参数 ctx 未被使用，可根据实际需求扩展。
func (l *KlogLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	control.fatal(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

// SetLevel 方法是 KiteX 日志记录器接口的一部分，但在这里没有实现具体的日志级别设置逻辑。
// 可根据实际需求添加日志级别设置的功能。
func (l *KlogLogger) SetLevel(klog.Level) {
}

// SetOutput 方法是 KiteX 日志记录器接口的一部分，但在这里没有实现具体的输出设置逻辑。
// 可根据实际需求添加输出设置的功能。
func (l *KlogLogger) SetOutput(io.Writer) {
}
