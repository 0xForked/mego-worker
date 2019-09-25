package helper

import "strconv"

type DataConfig struct {
	App      *AppConfig
	MySQL    *MySQlConfig
	RabbitMQ *RabbitMQConfig
}

type AppConfig struct {
	Name  string
	Debug bool
}

type MySQlConfig struct {
	Driver string
	URL    string
}

type RabbitMQConfig struct {
	URL string
}

func GetDataConfig() *DataConfig {
	appDebug, _ := strconv.ParseBool(GetEnv("APP_DEBUG", "false"))
	return &DataConfig{
		MySQL: &MySQlConfig{
			Driver: GetEnv("DB_DRIVER", "mysql"),
			URL:    GetEnv("DB_CONNECTION_URL", ""),
		},
		RabbitMQ: &RabbitMQConfig{
			URL: GetEnv("MQ_CONNECTION_URL", ""),
		},
		App: &AppConfig{
			Name:  GetEnv("APP_NAME", "mego_worker"),
			Debug: appDebug,
		},
	}
}
