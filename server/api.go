package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Api(r *gin.Engine, hub *Hub) {
	r.GET("/onlinecount", func(c *gin.Context) {
		c.String(http.StatusOK, "%d", len(hub.GetDistinct()))
	})
	r.GET("/onlineusers", func(c *gin.Context) {
		c.JSON(http.StatusOK, hub.GetUsers())
	})
	r.POST("/sendto", func(c *gin.Context) {
		var message Message
		if err := c.ShouldBindJSON(&message); err == nil {
			clients := hub.FindBy(message.SendTo)
			if message.Type == 0 && len(clients) > 0 {
				clients[0].Send(message.Data)
			} else if message.Type == 1 {
				for _, client := range clients {
					client.Send(message.Data)
				}
			} else if message.Type == 2 {
				hub.Broadcast(message.Data)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	r.GET("/getuserstatus", func(c *gin.Context) {
		users := c.QueryArray("code")
		clients := hub.FindBy(users)
		var onlineusers []string
		for _, u := range clients {
			onlineusers = append(onlineusers, u.user.Code)
		}
		c.JSON(http.StatusOK, onlineusers)
	})

	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, "1.0.2")
	})

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}
