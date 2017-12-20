package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Api(r *gin.Engine, hub *Hub) {
	r.GET("/onlinecount", func(c *gin.Context) {
		allowcors(c)
		c.String(http.StatusOK, "%d", len(hub.GetDistinct()))
	})
	r.GET("/onlineusers", func(c *gin.Context) {
		allowcors(c)
		c.JSON(http.StatusOK, hub.GetUsers())
	})
	r.POST("/sendto", func(c *gin.Context) {
		allowcors(c)
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

	r.GET("getuserstatus", func(c *gin.Context) {
		allowcors(c)
		users := c.QueryArray("code")
		clients := hub.FindBy(users)
		var onlineusers []string
		for _, u := range clients {
			onlineusers = append(onlineusers, u.user.Code)
		}
		c.JSON(http.StatusOK, onlineusers)
	})

	r.GET("version", func(c *gin.Context) {
		allowcors(c)
		c.JSON(http.StatusOK, "1.0.1")
	})
}

func allowcors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
}
