package config

import (
	"time"

	"github.com/spf13/viper"
)

type ConfigStruct struct {
	DatabaseHost              string        `mapstructure:"DATABASE_HOST"`
	DatabasePort              uint16        `mapstructure:"DATABASE_PORT"`
	DatabaseName              string        `mapstructure:"DATABASE_NAME"`
	DatabaseUser              string        `mapstructure:"DATABASE_USER"`
	DatabasePassword          string        `mapstructure:"DATABASE_PASSWORD"`
	DatabaseMaxConns          int32         `mapstructure:"DATABASE_MAX_CONNS"`
	DatabaseMinConns          int32         `mapstructure:"DATABASE_MIN_CONNS"`
	DatabaseMaxConnLifetime   time.Duration `mapstructure:"DATABASE_MAX_CONN_LIFETIME"`
	DatabaseMaxConnIdleTime   time.Duration `mapstructure:"DATABASE_MAX_CONN_IDLE_TIME"`
	DatabaseHealthCheckPeriod time.Duration `mapstructure:"DATABASE_HEALTH_CHECK_PERIOD"`
	DatabaseConnectTimeout    time.Duration `mapstructure:"DATABASE_CONNECT_TIMEOUT"`

	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	SecretKey string `mapstructure:"SECRET_KEY"`
}

var Config ConfigStruct

func Initialize() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(&Config)
}
