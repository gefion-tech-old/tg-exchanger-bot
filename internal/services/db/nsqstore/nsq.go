package nsqstore

import (
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/nsqio/go-nsq"
)

// Инициализации NSQ потребителя
func Init(config *config.NsqConfig, topic, channel string) (*nsq.Consumer, error) {
	c := nsq.NewConfig()
	c.MaxAttempts = 10
	c.MaxInFlight = 5
	c.MaxRequeueDelay = time.Second * 900
	c.DefaultRequeueDelay = time.Second * 0

	consumer, err := nsq.NewConsumer(topic, channel, c)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
