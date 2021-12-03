package config

type Config struct {
	Bot     BotConfig     `toml:"bot"`
	API     ApiConfig     `toml:"api"`
	NSQ     NsqConfig     `toml:"nsq"`
	General GeneralConfig `toml:"general"`
}

type BotConfig struct {
	Token string `toml:"token"`
	Debug bool   `toml:"debug"`
}

type ApiConfig struct {
	Url string `toml:"url"`
}

type NsqConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type GeneralConfig struct {
	Cores int `toml:"cores"`
}

func Init() *Config {
	return &Config{}
}
