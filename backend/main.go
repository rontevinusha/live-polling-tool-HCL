package main

import (
	"net/http"
	gorillaWebsocket "github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"live-pooling-backend/middleware"
	"live-pooling-backend/controllers"
	"live-pooling-backend/services"
	"live-pooling-backend/websocket"
)
var upgrader = gorillaWebsocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func main() {

	err := services.ConnectMongoDB()
	if err != nil {
		panic(err)
	}

	err = services.ConnectRedis()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	
r.Use(func(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
})
r.POST("/login", controllers.Login)
	hub := websocket.NewHub()
	websocket.CurrentHub = hub

	r.POST("/polls", middleware.AuthMiddleware(),controllers.CreatePoll)
	r.GET("/polls", controllers.GetPolls)
	r.POST("/polls/:id/vote/:optionId", controllers.VotePoll)

	r.GET("/ws", func(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	hub.AddClient(conn)

	defer hub.RemoveClient(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Live Polling Tool Backend is Running!",
		})
	})

	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}