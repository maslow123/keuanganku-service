package config

import "github.com/spf13/viper"

type Config struct {
	Port              string `mapstructure:"PORT"`
	DBUrl             string `mapstructure:"DB_URL"`
	PosServiceUrl     string `mapstructure:"POS_SERVICE_URL"`
	BalanceServiceUrl string `mapstructure:"BALANCE_SERVICE_URL"`
}

func LoadConfig(path string, filename string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
