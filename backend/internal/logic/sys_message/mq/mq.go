package mq

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// Publish sends a message to the specified Redis topic
func Publish(ctx context.Context, topic string, message interface{}) error {
	jsonStr, err := gjson.EncodeString(message)
	if err != nil {
		return err
	}
	_, err = g.Redis().Publish(ctx, topic, jsonStr)
	return err
}

// Subscribe listens to a Redis topic and processes messages
// This is a blocking call, usually run in a goroutine
func Subscribe(ctx context.Context, topic string, handler func(ctx context.Context, message string) error) {
	conn, err := g.Redis().Conn(ctx)
	if err != nil {
		g.Log().Error(ctx, "Failed to connect to redis for subscription", err)
		return
	}
	defer conn.Close(ctx)

	_, err = conn.Subscribe(ctx, topic)
	if err != nil {
		g.Log().Error(ctx, "Failed to subscribe to topic", topic, err)
		return
	}

	for {
		msg, err := conn.ReceiveMessage(ctx)
		if err != nil {
			g.Log().Error(ctx, "Error receiving message", err)
			// Decide whether to break or retry based on error type
			// For simplicity, we break here to avoid infinite loop on fatal errors
			break
		}
		if err := handler(ctx, msg.Payload); err != nil {
			g.Log().Error(ctx, "Error handling message", err)
		}
	}
}
