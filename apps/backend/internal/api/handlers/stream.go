package handlers

import (
	"context"

	"github.com/dis70rt/flowback/internal/pubsub"
	"github.com/gin-gonic/gin"
)

func StreamHandler(bus pubsub.Subscriber) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		msgChan := make(chan []byte)

		bus.Subscribe(c.Request.Context(), "dashboard_updates", func(ctx context.Context, payload []byte) {
			msgChan <- payload
		})

		for {
			select {
			case <-c.Request.Context().Done():
				return // Client disconnected
			case msg := <-msgChan:
				c.SSEvent("message", string(msg))
				c.Writer.Flush()
			}
		}
	}
}
