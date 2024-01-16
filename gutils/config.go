package gutils

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Dcarbon/go-shared/libs/utils"
)

// Config :
type Config struct {
	Env        string               // prod || dev || stg
	Port       int                  //
	DbUrl      string               //
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

func (config *Config) GetDBUrl() string {
	urlParsed, err := url.Parse(config.DbUrl)
	utils.PanicError("Parse db url ", err)
	var query = urlParsed.Query()

	if !strings.Contains(config.DbUrl, "sslmode") {
		query.Add("sslmode", "disable")
	}

	if !strings.Contains(config.DbUrl, "application_name") {
		query.Add("application_name", config.Name)
	}

	urlParsed.RawQuery = query.Encode()
	return urlParsed.String()
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

func (config *Config) GetIotHost() string {
	return config.GetOption(ISVIotInfo)
}

func (config *Config) GetStorageHost() string {
	return config.GetOption(ISVStorage)
}

func (config *Config) GetUser() string {
	return config.Options[ISVUser]
}

func (config *Config) GetPassword() string {
	return config.Options[ISVPass]
}
