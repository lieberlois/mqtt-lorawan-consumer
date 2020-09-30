package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	MqttBroker Mqtt   `mapstructure:"mqtt_broker"`
	InfluxDB   Influx `mapstructure:"influx_db"`
	Parser     Parser `mapstructure:"parser"`
}

type Mqtt struct {
	Url      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Topic    string `mapstructure:"topic"`
}

type Influx struct {
	Url      string `mapstructure:"url"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Parser struct {
	MeasurementKey string `mapstructure:"measurement_key"`
	TagsetKey      string `mapstructure:"tagset_key"`
	ValuesKey      string `mapstructure:"values_key"`
}

func LoadConfig(cfg *Config) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()
	err := viper.Unmarshal(cfg)

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
