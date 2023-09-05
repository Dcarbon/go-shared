package gutils

// Config :
type Config struct {
	Port       int                  //
	DbURL      string               //
	JwtKey     string               //
	Name       string               // Service name
	Options    map[string]string    //
	AuthConfig map[string]*ARConfig // Authen of service
}

func (config *Config) GetOption(key string) string {
	if nil == config.Options {
		return ""
	}
	return config.Options[key]
}

func (config *Config) GetAMQPUrl() string {
	return config.Options["AMQP"]
}

func (config *Config) GetRedisUrl() string {
	return config.Options["REDIS"]
}

func (config *Config) GetUser() string {
	return config.Options[ISVUser]
}

func (config *Config) GetPassword() string {
	return config.Options[ISVPass]
}
