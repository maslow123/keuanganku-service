package config

import "github.com/spf13/viper"

type Config struct {
	Port                   string `mapstructure:"PORT"`
	DBUrl                  string `mapstructure:"DB_URL"`
	JWTSecretKey           string `mapstructure:"JWT_SECRET_KEY"`
	BalanceServiceUrl      string `mapstructure:"BALANCE_SERVICE_URL"`
	CloudinaryCloudName    string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CloudinaryApiKey       string `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryApiSecretKey string `mapstructure:"CLOUDINARY_API_SECRET_KEY"`
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
