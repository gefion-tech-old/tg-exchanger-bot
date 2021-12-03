package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/gefion-tech/tg-exchanger-bot/internal/app"
	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/config.local.toml", "Path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Инициализирую конфигурацию
	config := config.Init()
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		panic(err)
	}

	// Инициализация модуля приложения
	application := app.Init(config)
	if err := application.Start(ctx); err != nil {
		fmt.Println(err)
	}

}
