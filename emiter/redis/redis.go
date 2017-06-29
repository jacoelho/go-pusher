package redis

import (
	"log"
	"time"

	r "github.com/go-redis/redis"
)

type client struct {
	redis  *r.Client
	log    *log.Logger
	events chan []byte
}

func New(url string, db int, l *log.Logger) *client {
	if url == "" {
		url = "127.0.0.1:6379"
	}

	conn := r.NewClient(&r.Options{
		Addr: url,
		DB:   db,
	})

	return &client{
		redis:  conn,
		events: make(chan []byte, 256),
		log:    l,
	}
}

func (c *client) Watch(channels ...string) error {
	sub := c.redis.Subscribe(channels...)
	ticker := time.NewTicker(time.Second * 5)

	// health check routine
	// TODO exponential backoff
	go func() {
		for {
			select {
			case <-ticker.C:
				_, err := c.redis.Ping().Result()
				if err != nil {
					c.log.Println("redis connection issue")
					break
				}
			}
		}

		ticker.Stop()
	}()

	// listen to redis pubsub messages
	go func() {
		defer sub.Close()

		for {
			message, err := sub.ReceiveMessage()

			if err != nil {
				c.log.Println(err)
				break
			}

			c.events <- []byte(message.Payload)
		}

	}()

	return nil
}

func (c *client) Events() <-chan []byte {
	return c.events
}
