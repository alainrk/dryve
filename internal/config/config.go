package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTPConfig     `mapstructure:"http"`
	Limits   LimitsConfig   `mapstructure:"limits"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Database DatabaseConfig `mapstructure:"database"`
}

type HTTPConfig struct {
	Port int `mapstructure:"port" default:"8666"`
}

type LimitsConfig struct {
	MaxFileSize            int64 `mapstructure:"max_file_size" default:"52428800"`
	FileEndpointsRateLimit int   `mapstructure:"file_endpoints_rate_limit" default:"10"`
}

type StorageConfig struct {
	Path string `mapstructure:"path" default:"/tmp/dryve-file-uploader"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver" default:"postgres"`
	Host     string `mapstructure:"host" default:"localhost"`
	Port     int    `mapstructure:"port" default:"5432"`
	User     string `mapstructure:"user" default:"not_set_user"`
	Password string `mapstructure:"password" default:"not_set_password"`
	Database string `mapstructure:"db_name" default:"not_set_db_name"`
}

// NewConfig creates a new config
// It reads the config file and unmarshals it into a Config struct
func NewConfig(file string) Config {
	var config Config

	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("cannot read from a config, %v", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return config
}
