package constants

const (
	// LogFilePath 对应 ${pwd}/{LogFilePath}/log.log 相对于当前运行路径而言
	LogFilePath = "log"

	// wd/log/{ServiceName}/data/*.log
	LogFilePathTemplate      = "%s/%s/%s/%s/std.log"
	ErrorLogFilePathTemplate = "%s/%s/%s/%s/stderr.log"

	// DefaultLogLevel 是默认的日志等级. Supported Level: debug info warn error fatal
	DefaultLogLevel = "INFO"

	StackTraceKey = "stacktrace"
	ServiceKey    = "service"
	SourceKey     = "source"
	ErrorCodeKey  = "error_code"

	KlogSource  = "klog"
	MysqlSource = "mysql"
	RedisSource = "redis"
	KafkaSource = "kafka"
)
