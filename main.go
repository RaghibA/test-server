package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var randSeed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var req int32

func randomString(length int) string {
	c := make([]byte, length)
	for i := range c {
		c[i] = charset[randSeed.Intn(len(charset))]
	}

	return string(c)
}

func main() {
	serverId := randomString(6)

	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		atomic.AddInt32(&req, 1)

		c.JSON(http.StatusOK, gin.H{
			"message":  "Test server GET",
			"time":     time.Now(),
			"serverId": serverId,
			"request":  req,
		})
	})

	r.POST("/test", func(c *gin.Context) {
		atomic.AddInt32(&req, 1)

		var PostBody struct {
			Data string `json:"data"`
		}

		if c.Bind(&PostBody) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  PostBody.Data,
			"time":     time.Now(),
			"serverId": serverId,
			"request":  req,
		})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server err: %v", err)
		}
	}()

	<-stop
	log.Println("Attempting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("graceful shutdown complete")
}
