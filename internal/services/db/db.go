package db

import (
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/go-redis/redis/v7"
)

// Функция инициализации Redis БД
func InitRedis(config *config.RedisConfig, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		DB:   db,
	})

	// Тест соединения
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}
