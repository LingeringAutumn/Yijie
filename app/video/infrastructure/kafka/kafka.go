package kafka

import (
	"github.com/IBM/sarama"

	"github.com/LingeringAutumn/Yijie/config"
)

// NewProducerConfig 构造带 SASL 的 Kafka Producer 配置
func NewProducerConfig() (*sarama.Config, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Return.Successes = true

	cfg.Net.SASL.Enable = true
	cfg.Net.SASL.User = config.Kafka.SASLUser
	cfg.Net.SASL.Password = config.Kafka.SASLPassword
	cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	cfg.Net.TLS.Enable = false
	cfg.Version = sarama.V2_8_0_0

	return cfg, nil
}
