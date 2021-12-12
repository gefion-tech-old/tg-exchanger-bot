package config

type Config struct {
	Bot   BotConfig   `toml:"bot"`
	API   ApiConfig   `toml:"api"`
	NSQ   NsqConfig   `toml:"nsq"`
	Redis RedisConfig `toml:"redis"`
}

type BotConfig struct {
	Token      string   `toml:"token"`
	Developers []string `toml:"developers"`
	Debug      bool     `toml:"debug"`
}

type ApiConfig struct {
	Url string `toml:"url"`
}

type NsqConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type RedisConfig struct {
	Host string `toml:"host"`
	Port int64  `toml:"port"`
}

func Init() *Config {
	return &Config{}
}
