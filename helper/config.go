package helper

type DataConfig struct {
	MySQL *MySQlConfig
	RabbitMQ *RabbitMQConfig
}

type MySQlConfig struct {
	Driver  string
	URL		string
}

type RabbitMQConfig struct {
	URL     string
}

func GetDataConfig() *DataConfig {
	return &DataConfig{
		MySQL: &MySQlConfig{
			Driver:  GetEnv("DB_DRIVER", "mysql"),
			URL:     GetEnv("DB_CONNECTION_URL", ""),
		},
		RabbitMQ: &RabbitMQConfig{
			URL:     GetEnv("MQ_CONNECTION_URL", ""),
		},
	}
}