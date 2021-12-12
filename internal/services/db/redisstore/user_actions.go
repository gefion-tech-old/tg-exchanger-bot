package redisstore

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/go-redis/redis/v7"
)

type UserActions struct {
	client *redis.Client
}

type UserActionsI interface {
	New(chatID int64, payload *models.UserAction) error
	Fetch(chatID int64) (map[string]interface{}, error)
	Delete(chatID int64) error
	Close() error
	Clear()
}

func InitUserActions(c *redis.Client) UserActionsI {
	return &UserActions{
		client: c,
	}
}

// Создать новое пользовательское действие
func (c *UserActions) New(chatID int64, payload *models.UserAction) error {
	// Генерирую объект для записи в Redis
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return c.client.Set(fmt.Sprintf("%d", chatID), b, 60*time.Minute).Err()
}

// Проверка наличия у пользователя начатого действия
func (c *UserActions) Fetch(chatID int64) (map[string]interface{}, error) {
	d, err := c.client.Get(fmt.Sprintf("%d", chatID)).Result()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(d), &data); err != nil {
		return nil, err
	}

	return data, nil
}

// Удалить пользовательское действие
func (c *UserActions) Delete(chatID int64) error {
	return c.client.Del(fmt.Sprintf("%d", chatID)).Err()
}

func (c *UserActions) Close() error {
	return c.client.Close()
}

func (c *UserActions) Clear() {
	c.client.FlushAllAsync()
}
