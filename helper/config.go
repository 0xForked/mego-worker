package helper

import "strconv"

type Config struct {
	App      *AppConfig
	Service  *ServiceConfig
	MySQL    *MySQlConfig
	RabbitMQ *RabbitMQConfig
}

type AppConfig struct {
	Name  string
	Debug bool
}

type ServiceConfig struct {
	Delivery string
}

type MySQlConfig struct {
	Driver string
	URL    string
}

type RabbitMQConfig struct {
	URL string
}

func GetConfig() *Config {
	appDebug, _ := strconv.ParseBool(GetEnv("APP_DEBUG", "false"))
	return &Config{
		App: &AppConfig{
			Name:  GetEnv("APP_NAME", "mego_worker"),
			Debug: appDebug,
		},
		Service: &ServiceConfig{
			Delivery: GetEnv("DELIVERY_MODE", "schema"),
		},
		MySQL: &MySQlConfig{
			Driver: GetEnv("DB_DRIVER", "mysql"),
			URL:    GetEnv("DB_CONNECTION_URL", ""),
		},
		RabbitMQ: &RabbitMQConfig{
			URL: GetEnv("MQ_CONNECTION_URL", ""),
		},
	}
}
