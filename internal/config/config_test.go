package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	given := NewConfig("../../test/testconfig.json")
	exp := Config{
		HTTP: HTTPConfig{
			Port: 8666,
		},
		Limits: LimitsConfig{
			MaxFileSize:            52428800,
			FileEndpointsRateLimit: 10,
		},
		Storage: StorageConfig{
			Path: "/tmp/dryve-filestorage",
		},
		Database: DatabaseConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			User:     "server",
			Password: "passwordsuperstrong",
			Database: "main",
		},
	}

	if !reflect.DeepEqual(given, exp) {
		t.Errorf("given config: %v, expected: %v", given, exp)
	}
}

func TestNewConfigDefaults(t *testing.T) {
	given := NewConfig("../../test/testconfig_defaults.json")
	exp := Config{
		HTTP: HTTPConfig{
			Port: 8666,
		},
		Limits: LimitsConfig{
			MaxFileSize:            52428800,
			FileEndpointsRateLimit: 10,
		},
		Storage: StorageConfig{
			Path: "/tmp/dryve-file-uploader",
		},
		Database: DatabaseConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			User:     "not_set_user",
			Password: "not_set_password",
			Database: "not_set_db_name",
		},
	}

	if !reflect.DeepEqual(given, exp) {
		t.Errorf("given config: %v, expected: %v", given, exp)
	}
}
