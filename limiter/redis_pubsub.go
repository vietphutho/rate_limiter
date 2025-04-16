package limiter

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func SubscribeTokenSync(client *redis.Client, channel string, onMessage func(msg string)) {
	ctx := context.Background()
	pubsub := client.Subscribe(ctx, channel)

	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Println("PubSub receive error:", err)
				continue
			}
			onMessage(msg.Payload)
		}
	}()
}

func PublishTokenUpdate(client *redis.Client, channel, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Publish(ctx, channel, message).Err(); err != nil {
		log.Println("PubSub publish error:", err)
	}
}
