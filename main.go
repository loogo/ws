package main

import (
	"flag"
	"net/http"
	"log"
	"./server"
	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(c *gin.Context) {
	w, r := c.Writer, c.Request
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := server.NewHub()
	go hub.Run()
	r := gin.Default()
	r.GET("/", serveHome)
	r.GET("/ws", func(c *gin.Context) {
		w, r := c.Writer, c.Request
		server.ServeWs(hub, w, r)
	})

	r.GET("/onlinecount", func(c *gin.Context) {
		c.String(http.StatusOK, "%d", len(hub.GetDistinct()))
	})

	r.GET("/onlineusers", func(c *gin.Context) {
		c.JSON(http.StatusOK, hub.GetUsers())
	})

	r.POST("/sendto", func(c *gin.Context) {
		var message server.Message
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

	err := r.Run(*addr)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
