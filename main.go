package main

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/notification-counters", handler)

	if err := router.Run(":80"); err != nil {
		panic(err)
	}
}

func handler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		c.SSEvent("message", fmt.Sprintf("Hello, %s", time.Now().Format(time.RFC3339)))
		return true
	})
}
