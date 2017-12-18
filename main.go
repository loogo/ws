package main

import (
	"flag"
	"net/http"
	"log"
	"github.com/loogo/ws/server"
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

	err := r.Run(*addr)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
