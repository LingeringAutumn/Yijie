package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"go.uber.org/zap"
)

// Debug 函数用于记录调试级别的日志信息。
// 它接收一个字符串消息和可选的 zap.Field 类型的字段参数，
// 并调用 control 的 debug 方法来记录日志。
func Debug(msg string, fields ...zap.Field) {
	control.debug(msg, fields...)
}

// Debugf 函数用于以格式化的方式记录调试级别的日志信息。
// 它接收一个格式化模板字符串和可变数量的参数，
// 并调用 control 的 debugf 方法来记录日志。
func Debugf(template string, args ...interface{}) {
	control.debugf(template, args...)
}

// Info 函数用于记录信息级别的日志信息。
// 它接收一个字符串消息和可选的 zap.Field 类型的字段参数，
// 并调用 control 的 info 方法来记录日志。
func Info(msg string, fields ...zap.Field) {
	control.info(msg, fields...)
}

// Infof 函数用于以格式化的方式记录信息级别的日志信息。
// 它接收一个格式化模板字符串和可变数量的参数，
// 并调用 control 的 infof 方法来记录日志。
func Infof(template string, args ...interface{}) {
	control.infof(template, args...)
}

// Warn 函数用于记录警告级别的日志信息。
// 它接收一个字符串消息和可选的 zap.Field 类型的字段参数，
// 并调用 control 的 warn 方法来记录日志。
func Warn(msg string, fields ...zap.Field) {
	control.warn(msg, fields...)
}

// Warnf 函数用于以格式化的方式记录警告级别的日志信息。
// 它接收一个格式化模板字符串和可变数量的参数，
// 并调用 control 的 warnf 方法来记录日志。
func Warnf(template string, args ...interface{}) {
	control.warnf(template, args...)
}

// Error 函数用于记录错误级别的日志信息。
// 它接收一个字符串消息和可选的 zap.Field 类型的字段参数，
// 并调用 control 的 error 方法来记录日志。
func Error(msg string, fields ...zap.Field) {
	control.error(msg, fields...)
}

// Errorf 函数用于以格式化的方式记录错误级别的日志信息。
// 它接收一个格式化模板字符串和可变数量的参数，
// 并调用 control 的 errorf 方法来记录日志。
func Errorf(template string, args ...interface{}) {
	control.errorf(template, args...)
}

// Fatal 函数用于记录致命级别的日志信息。
// 它接收一个字符串消息和可选的 zap.Field 类型的字段参数，
// 并调用 control 的 fatal 方法来记录日志。
func Fatal(msg string, fields ...zap.Field) {
	control.fatal(msg, fields...)
}

// Fatalf 函数用于以格式化的方式记录致命级别的日志信息。
// 它接收一个格式化模板字符串和可变数量的参数，
// 并调用 control 的 fatalf 方法来记录日志。
func Fatalf(template string, args ...interface{}) {
	control.fatalf(template, args...)
}

// LogError 函数用于记录错误信息。
// 如果传入的错误不为空，则调用 control 的 error 方法记录错误信息。
func LogError(err error) {
	if err != nil {
		control.error(err.Error())
	}
}

// permission 常量定义了文件和目录的权限。
// 用户具有读、写、执行权限，组用户和其他用户具有读、执行权限。
const permission = 0o755

// getCurrentDirectory 函数用于获取当前运行程序的目录。
// 它接收一个服务名称作为参数，用于判断是否在容器环境下运行。
// 返回当前运行目录的路径和可能出现的错误。
func getCurrentDirectory(serviceName string) (string, error) {
	// 获取当前可执行文件的绝对路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	// 将路径中的反斜杠替换为正斜杠，以统一路径分隔符
	path := strings.ReplaceAll(dir, "\\", "/")
	// 将路径按正斜杠分割成多个部分
	paths := strings.Split(path, "/")

	// 在容器下运行时，让所有日志都集中于 output/log/$(ServiceName)
	// 如果路径的最后一部分是服务名称，则将路径缩短到上一级目录
	if paths[len(paths)-1] == serviceName {
		path = path[:len(path)-len(serviceName)-1]
	}

	return path, nil
}

// checkAndOpenFile 函数用于检查并打开指定路径的文件。
// 如果文件所在的目录不存在，则创建该目录。
// 返回一个打开的文件句柄。
func checkAndOpenFile(path string) *os.File {
	var err error
	var handler *os.File
	// 创建文件所在的目录，如果目录已存在则不会报错
	if err = os.MkdirAll(filepath.Dir(path), permission); err != nil {
		// 如果创建目录失败，则触发 panic
		panic(err)
	}

	// 以追加、创建和只写的模式打开文件
	handler, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, permission)
	if err != nil {
		// 如果打开文件失败，则触发 panic
		panic(err)
	}
	// 设置文件句柄的终结器，当文件句柄不再被引用时，自动关闭文件
	runtime.SetFinalizer(handler, func(fd *os.File) {
		if err := fd.Close(); err != nil {
			// 如果关闭文件失败，则记录错误信息
			Infof("close file failed %v", err)
		}
	})
	return handler
}
