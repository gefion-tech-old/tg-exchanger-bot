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
	prod bool
	cpu  int
)

func init() {
	flag.BoolVar(&prod, "prod", false, "Strat on production server.")
	flag.IntVar(&cpu, "cpu", 2, "Number of processor threads")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Инициализирую конфигурацию
	var cnf *config.Config
	if prod {
		cnf = config.Init()
		if _, err := toml.DecodeFile("config/config.prod.toml", cnf); err != nil {
			panic(err)
		}

	} else {
		cnf = config.Init()
		if _, err := toml.DecodeFile("config/config.local.toml", cnf); err != nil {
			panic(err)
		}
	}

	// Инициализация модуля приложения
	application := app.Init(cnf)
	if err := application.Start(ctx); err != nil {
		fmt.Println(err)
	}

}
