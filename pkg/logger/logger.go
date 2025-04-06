package logger

import (
	"fmt"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/cloudwego/kitex/pkg/klog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// controlLogger 结构体用于管理日志记录器的配置和状态。
// 它包含一个读写锁、一个日志记录器实例、日志钩子、一个原子布尔值用于标记是否已开始定时更新，以及服务名称。
type controlLogger struct {
	mu      sync.RWMutex
	logger  *logger
	hooks   []func(zapcore.Entry) error
	done    atomic.Bool
	svcName string
}

// logger 结构体是对 zap.Logger 的封装，方便后续操作。
type logger struct {
	*zap.Logger
}

// 全局变量部分
var (
	// control 是 controlLogger 的实例，用于管理日志记录器。
	control controlLogger
	// logLevel 表示日志的级别，默认为 Info 级别。
	logLevel = zapcore.InfoLevel
	// callerSkip 表示调用栈的跳过层数，用于定位日志调用位置。
	callerSkip = 2
	// logFileHandler 是一个原子值，用于存储日志文件的句柄。
	logFileHandler atomic.Value
	// stdErrFileHandler 是一个原子值，用于存储错误日志文件的句柄，全局变量避免被 GC 回收。
	stdErrFileHandler atomic.Value
	// defaultService 是默认的服务名称。
	defaultService = "_default"
)

// init 函数主要用于在 logger.Init 之前输出日志。
// 它会构建一个默认的日志配置，并初始化控制日志记录器和日志钩子。
// 同时，将 Klog 的日志记录器设置为当前实现的日志记录器。
func init() {
	cfg := buildConfig(nil)
	control.logger = &logger{BuildLogger(cfg, control.addZapOptions(defaultService)...)}
	control.hooks = make([]func(zapcore.Entry) error, 0)

	klog.SetLogger(GetKlogLogger())
}

// Init 函数用于初始化日志记录器。
// 它接收服务名称和日志级别作为参数。
// 如果服务名称为空，会触发 panic。
// 它会设置服务名称、解析日志级别、更新日志记录器，并启动定时更新日志记录器的任务。
func Init(service string, level string) {
	if service == "" {
		panic("server should not be empty")
	}

	control.svcName = service
	logLevel = parseLevel(level)
	control.updateLogger(service)
	control.scheduleUpdateLogger(service)
}

// Ignore 函数将日志记录器设置为无操作的日志记录器，即忽略所有日志输出。
func Ignore() {
	control.logger = &logger{zap.NewNop()}
}

// AddLoggerHook 函数用于添加日志钩子。
// 日志钩子会在每次日志输出后执行。
func AddLoggerHook(fns ...func(zapcore.Entry) error) {
	control.hooks = append(control.hooks, fns...)
}

// scheduleUpdateLogger 函数用于启动定时更新日志记录器的任务。
// 确保只开启一次定时更新，每天凌晨更新日志记录器。
func (l *controlLogger) scheduleUpdateLogger(service string) {
	// 确保只开启一次定时更新
	if !l.done.Load() {
		l.done.Store(true)
		go func() {
			for {
				now := time.Now()
				// 计算下一个零点的时间
				next := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
				// 等待到下一个零点
				time.Sleep(time.Until(next))
				// 更新日志记录器
				l.updateLogger(service)
			}
		}()
	}
}

// updateLogger 函数用于更新日志记录器。
// 它会获取当前目录，设置日志文件和错误日志文件的路径，打开文件，
// 并根据日志级别将日志输出到不同的文件，最后重新构建日志记录器。
func (l *controlLogger) updateLogger(service string) {
	// 避免 logger 更新时引发竞态
	l.mu.Lock()
	defer l.mu.Unlock()
	var err error
	var pwd string

	// 获取当前目录
	if pwd, err = getCurrentDirectory(l.svcName); err != nil {
		panic(err)
	}

	// 设置文件输出的位置
	date := time.Now().Format("2006-01-02")
	logPath := fmt.Sprintf(constants.LogFilePathTemplate, pwd, constants.LogFilePath, service, date)
	stderrPath := fmt.Sprintf(constants.ErrorLogFilePathTemplate, pwd, constants.LogFilePath, service, date)

	// 打开文件,并设置无引用时关闭文件
	logFileHandler.Store(checkAndOpenFile(logPath))
	stdErrFileHandler.Store(checkAndOpenFile(stderrPath))

	// 让日志输出到不同的位置
	logLevelFn := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= logLevel
	})
	errLevelFn := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > logLevel
	})

	logCore := zapcore.NewCore(defaultEnc(), zapcore.Lock(logFileHandler.Load().(*os.File)), logLevelFn)    //nolint
	errCore := zapcore.NewCore(defaultEnc(), zapcore.Lock(stdErrFileHandler.Load().(*os.File)), errLevelFn) //nolint

	cfg := buildConfig(zapcore.NewTee(logCore, errCore))

	l.logger.Logger = BuildLogger(cfg, l.addZapOptions(service)...)
}

// addZapOptions 函数用于添加 zap 日志记录器的选项。
// 它会添加日志钩子、调用者信息和服务相关的字段信息。
func (l *controlLogger) addZapOptions(serviceName string) []zap.Option {
	var opts []zap.Option
	if len(l.hooks) != 0 {
		opts = append(opts, zap.Hooks(l.hooks...))
	}
	opts = append(opts, zap.AddCaller())
	opts = append(opts, zap.AddCallerSkip(callerSkip))
	opts = append(opts, zap.Fields(
		zap.String(constants.ServiceKey, serviceName),
		zap.String(constants.SourceKey, fmt.Sprintf("app-%s", serviceName)),
	))

	return opts
}

// debug 函数用于输出调试级别的日志。
// 它会先获取读锁，确保在读取日志记录器时不会被其他更新操作干扰，然后调用日志记录器的 Debug 方法输出日志。
func (l *controlLogger) debug(msg string, fields ...zap.Field) {
	l.mu.RLock() // 锁的是 logger 的操作权限, 而不是写操作, 写操作在 zap.logger 的内部有锁.
	defer l.mu.RUnlock()
	l.logger.Debug(msg, fields...)
}

// debugf 函数用于以格式化字符串的方式输出调试级别的日志。
// 它会先获取读锁，然后将格式化后的字符串作为信息传递给日志记录器的 Info 方法输出。
func (l *controlLogger) debugf(template string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Info(fmt.Sprintf(template, args...))
}

// info 函数用于输出信息级别的日志。
// 它会先获取读锁，然后调用日志记录器的 Info 方法输出日志。
func (l *controlLogger) info(msg string, fields ...zap.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Info(msg, fields...)
}

// infof 函数用于以格式化字符串的方式输出信息级别的日志。
// 它会先获取读锁，然后将格式化后的字符串作为信息传递给日志记录器的 Info 方法输出。
func (l *controlLogger) infof(template string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Info(fmt.Sprintf(template, args...))
}

// warn 函数用于输出警告级别的日志。
// 它会先获取读锁，然后调用日志记录器的 Warn 方法输出日志。
func (l *controlLogger) warn(msg string, fields ...zap.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Warn(msg, fields...)
}

// warnf 函数用于以格式化字符串的方式输出警告级别的日志。
// 它会先获取读锁，然后将格式化后的字符串作为信息传递给日志记录器的 Warn 方法输出。
func (l *controlLogger) warnf(template string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Warn(fmt.Sprintf(template, args...))
}

// error 函数用于输出错误级别的日志。
// 它会先获取读锁，然后调用日志记录器的 Error 方法输出日志。
func (l *controlLogger) error(msg string, fields ...zap.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Error(msg, fields...)
}

// errorf 函数用于以格式化字符串的方式输出错误级别的日志。
// 它会先获取读锁，然后将格式化后的字符串作为信息传递给日志记录器的 Error 方法输出。
func (l *controlLogger) errorf(template string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Error(fmt.Sprintf(template, args...))
}

// fatal 函数用于输出致命级别的日志。
// 它会先获取读锁，然后调用日志记录器的 Fatal 方法输出日志。
func (l *controlLogger) fatal(msg string, fields ...zap.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Fatal(msg, fields...)
}

// fatalf 函数用于以格式化字符串的方式输出致命级别的日志。
// 它会先获取读锁，然后将格式化后的字符串作为信息传递给日志记录器的 Fatal 方法输出。
func (l *controlLogger) fatalf(template string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Fatal(fmt.Sprintf(template, args...))
}

// parseLevel 函数用于解析日志级别字符串。
// 它会将输入的字符串转换为对应的 zapcore.Level 类型。
// 如果输入的字符串不匹配任何已知级别，则默认返回 Info 级别。
func parseLevel(level string) zapcore.Level {
	var lvl zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = zapcore.DebugLevel
	case "info":
		lvl = zapcore.InfoLevel
	case "warn":
		lvl = zapcore.WarnLevel
	case "error":
		lvl = zapcore.ErrorLevel
	case "fatal":
		lvl = zapcore.FatalLevel
	default:
		lvl = zapcore.InfoLevel
	}
	return lvl
}
