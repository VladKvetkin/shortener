package config

type Config struct {
	host string
	port string
}

func NewConfig(host string, port string) Config {
	return Config{
		host: host,
		port: port,
	}
}

func (c *Config) GetAddress() string {
	return c.host + ":" + c.port
}
