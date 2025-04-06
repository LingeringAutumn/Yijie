package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// config 结构体用于存储日志记录器的配置信息。
// 它包含了核心组件、编码器、写入器和日志级别等配置。
type config struct {
	core zapcore.Core        // 日志核心组件，负责处理日志的实际写入操作
	enc  zapcore.Encoder     // 日志编码器，用于将日志信息编码为特定格式（如 JSON 或控制台格式）
	ws   zapcore.WriteSyncer // 写入器，负责将编码后的日志信息写入目标位置（如文件、标准输出等）
	lvl  zapcore.Level       // 日志级别，决定了哪些级别的日志会被记录
}

// buildConfig 函数用于构建日志配置。
// 它接收一个 zapcore.Core 类型的参数，如果传入的 core 为 nil，则使用默认的配置创建一个新的 core。
// 最后返回一个指向 config 结构体的指针。
func buildConfig(core zapcore.Core) *config {
	cfg := defaultConfig() // 获取默认的日志配置
	cfg.core = core        // 将传入的 core 赋值给配置中的 core 字段
	if cfg.core == nil {
		// 如果传入的 core 为 nil，则使用默认的编码器、写入器和日志级别创建一个新的 core
		cfg.core = zapcore.NewCore(cfg.enc, cfg.ws, cfg.lvl)
	}

	return cfg
}

// BuildLogger 函数用于根据配置和可选的 zap 选项创建一个 zap.Logger 实例。
// 它接收一个 *config 类型的参数和可变数量的 zap.Option 类型的参数。
// 返回一个指向 zap.Logger 的指针。
func BuildLogger(cfg *config, opts ...zap.Option) *zap.Logger {
	// 使用配置中的 core 和传入的选项创建一个新的 zap.Logger 实例
	return zap.New(cfg.core, opts...)
}

// defaultConfig 函数用于获取默认的日志配置。
// 它调用 defaultEnc、defaultWs 和 defaultLvl 函数分别获取默认的编码器、写入器和日志级别。
// 返回一个指向 config 结构体的指针。
func defaultConfig() *config {
	return &config{
		enc: defaultEnc(), // 默认的编码器
		ws:  defaultWs(),  // 默认的写入器
		lvl: defaultLvl(), // 默认的日志级别
	}
}

// defaultEnc 函数用于获取默认的日志编码器。
// 它创建一个 zapcore.EncoderConfig 结构体，设置了时间键、日志级别键、调用者键等信息。
// 最后返回一个 JSON 编码器实例。在调试时，可以使用 zapcore.NewConsoleEncoder 替换。
func defaultEnc() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "time",                        // 时间键，用于在日志中表示时间信息
		LevelKey:       "level",                       // 日志级别键，用于在日志中表示日志级别
		CallerKey:      "caller",                      // 调用者键，用于在日志中表示调用日志记录的位置
		MessageKey:     "msg",                         // 消息键，用于在日志中表示日志消息
		LineEnding:     zapcore.DefaultLineEnding,     // 行结束符
		EncodeLevel:    zapcore.CapitalLevelEncoder,   // 日志等级大写编码
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // 时间格式采用 ISO8601 编码
		EncodeDuration: zapcore.StringDurationEncoder, // 持续时间编码为字符串
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 调用者位置采用短格式编码
	}

	// return zapcore.NewConsoleEncoder(cfg)
	return zapcore.NewJSONEncoder(cfg) // 返回 JSON 编码器
}

// defaultWs 函数用于获取默认的写入器。
// 它返回 os.Stdout，表示将日志信息写入标准输出。
func defaultWs() zapcore.WriteSyncer {
	return os.Stdout
}

// defaultLvl 函数用于获取默认的日志级别。
// 它返回 zapcore.DebugLevel，表示默认的日志级别为调试级别，即会记录所有级别的日志。
func defaultLvl() zapcore.Level {
	return zapcore.DebugLevel
}
