package main

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

func CreateAndRunHTTPServer() {
	router := gin.Default()

	notifGroupRoute := router.Group("notifications")
	{
		notifGroupRoute.GET("", getNotificationHandler)
		notifGroupRoute.POST("", sendNotificationHandler)
	}

	if err := router.Run(":80"); err != nil {
		panic(err)
	}
}

func sendNotificationHandler(c *gin.Context) {
	var payload NewNotificationChanParam
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	eventManager := GetEventManager()
	eventManager.Message <- payload

	c.JSON(200, gin.H{"status": "ok"})
}

func getNotificationHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	userID := c.GetHeader("X-User-ID")

	clientChannel := make(chan NewNotificationChanParam)

	eventManager := GetEventManager()
	eventManager.NewClients <- ChannelLifecycleEvent{
		UserID:  userID,
		Channel: clientChannel,
	}

	defer func() {
		eventManager.ClosedClients <- ChannelLifecycleEvent{
			UserID:  userID,
			Channel: clientChannel,
		}
	}()

	log.Println("New client connected: ", userID)

	var alreadyPing bool
	c.Stream(func(w io.Writer) bool {
		if !alreadyPing {
			c.SSEvent("ping", "pong")
			alreadyPing = true
			return true
		}

		if msg, ok := <-clientChannel; ok {
			c.SSEvent("new-notification", msg.Message)
			return true
		}

		return false
	})
}
