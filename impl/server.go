package impl

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewServer() {
	hub := NewHub()
	engine := initRouter(hub)
	err := engine.Run() // 默认运行在本机的8080端口
	if err != nil {
		fmt.Println(err)
	}
}

func initRouter(hub *Hub) *gin.Engine {
	router := gin.Default()
	ws := router.Group("/ws")
	{
		ws.GET("/connect", func(ctx *gin.Context) {
			ServeWs(ctx, hub)
		})

		ws.POST("/send", func(ctx *gin.Context) {
			Send(ctx, hub)
		})
	}

	return router
}
