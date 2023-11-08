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
