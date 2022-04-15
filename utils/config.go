package utils

import (
	"github.com/spf13/viper"
)

// Config - Application Config
type Config struct {
	ServerPort    string `mapstructure:"SERVER_PORT"`
	DBHost        string `mapstructure:"MONGO_HOST"`
	DBPort        string `mapstructure:"MONGO_PORT"`
	DBName        string `mapstructure:"MONGO_DATABASE_NAME"`
	AWSProfile    string `mapstructure:"AWS_PROFILE"`
	AWSRegion     string `mapstructure:"AWS_REGION"`
	AWSBucketName string `mapstructure:"AWS_BUCKET_NAME"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
