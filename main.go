package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var randSeed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(length int) string {
	c := make([]byte, length)
	for i := range c {
		c[i] = charset[rand.Intn(len(charset))]
	}

	return string(c)
}

func main() {
	req := 0
	serverId := randomString(6)

	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		req += 1

		c.JSON(http.StatusOK, gin.H{
			"message":  "Test server GET",
			"time":     time.Now(),
			"serverId": serverId,
			"request":  req,
		})
	})

	r.POST("/test", func(c *gin.Context) {

		var PostBody struct {
			Data string `json:"data"`
		}

		req += 1

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

	log.Fatal(r.Run())
}
