package nsqstore

import (
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/nsqio/go-nsq"
)

type BotConsumers struct {
	Verification *nsq.Consumer
}

func Init(cnf *config.NsqConfig) (*BotConsumers, func(...*nsq.Consumer), error) {
	vConsumer, err := configure(cnf, "messages", "telegram")
	if err != nil {
		return nil, nil, err
	}

	return &BotConsumers{
			Verification: vConsumer,
		}, func(c ...*nsq.Consumer) {
			if len(c) > 0 {
				for _, consumer := range c {
					consumer.Stop()
				}
			}
		}, nil
}

// Инициализации NSQ потребителя
func configure(config *config.NsqConfig, topic, channel string) (*nsq.Consumer, error) {
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
