package config

type server struct {
	Secret      string `mapstructure:"private-key"`
	PublicKey   string `mapstructure:"public-key"`
	Version     string
	Name        string
	LogLevel    string `mapstructure:"log-level"`
	IntranetUrl string `mapstructure:"intranet-url"`
}

type snowflake struct {
	DatacenterID int64 `mapstructure:"datacenter-id"`
}

type service struct {
	Name     string
	AddrList []string
	LB       bool `mapstructure:"load-balance"`
}

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type etcd struct {
	Addr string
}

// 先留着，可能要用
type rabbitMQ struct {
	Addr     string
	Username string
	Password string
}

type redis struct {
	Addr     string
	Password string
}

type kafka struct {
	Broker          string `mapstructure:"broker"`
	Topic           string `mapstructure:"topic"`
	ConsumerGroup   string `mapstructure:"consumer_group"`
	MaxConnections  int    `mapstructure:"max_connections"`
	MaxQPS          int    `mapstructure:"max_qps"`
	AutoOffsetReset string `mapstructure:"auto_offset_reset"`
	SASLUser        string `mapstructure:"sasl_user"`
	SASLPassword    string `mapstructure:"sasl_password"`
}

type minio struct {
	Endpoint  string `mapstructure:"endpoint"`   // eg: "127.0.0.1:9000"
	AccessKey string `mapstructure:"access-key"` // MinIO 用户名
	SecretKey string `mapstructure:"secret-key"` // MinIO 密码
	UseSSL    bool   `mapstructure:"use-ssl"`    // 是否使用 HTTPS
}

type config struct {
	Server    server
	Snowflake snowflake
	MySQL     mySQL
	Etcd      etcd
	RabbitMQ  rabbitMQ
	Redis     redis
	Kafka     kafka
	Minio     minio
}
