package config

import "github.com/spf13/viper"

type Config struct {
	Port                  string `mapstructure:"PORT"`
	UserServiceUrl        string `mapstructure:"USER_SERVICE_URL"`
	PosServiceUrl         string `mapstructure:"POS_SERVICE_URL"`
	TransactionServiceUrl string `mapstructure:"TRANSACTION_SERVICE_URL"`
}

func LoadConfig(path string, fileName string) (c Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(fileName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
