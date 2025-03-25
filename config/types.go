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
	Address  string
	Network  string
	User     string
	Password string
}

type minio struct {
	Addr        string
	AccessKey   string
	AccessKeyID string
	SecretKey   string
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
