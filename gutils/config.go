package gutils

import (
	"fmt"
	"os"
	"strconv"
)

// Config :
type Config struct {
	Env        string               // prod || dev || stg
	Port       int                  //
	DbURL      string               //
	JwtKey     string               //
	Name       string               // Service name
	Options    map[string]string    //
	AuthConfig map[string]*ARConfig // Authen of service
}

func (config *Config) GetOption(key string) string {
	v := ""
	if nil != config.Options {
		v = config.Options[key]
	}

	if v == "" {
		v = os.Getenv(key)
	}
	return v
}

func (config *Config) GetOptInt(key string) int64 {
	vStr := config.GetOption(key)
	v, err := strconv.ParseInt(vStr, 10, 64)
	if nil != err {
		panic(fmt.Sprintf("Config %s should be string.", key))
	}
	return v
}

func (config *Config) GetAMQPUrl() string {
	return config.GetOption("AMQP_URL")
}

func (config *Config) GetRedisUrl() string {
	return config.GetOption("REDIS_URL")
}

func (config *Config) GetUser() string {
	return config.Options[ISVUser]
}

func (config *Config) GetPassword() string {
	return config.Options[ISVPass]
}
