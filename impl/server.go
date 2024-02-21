package impl

import (
	"chsir-zy/msg-push-center/config"
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewServer() {
	config.LoadConfig()
	config.GORM_DB = config.GormMysql()

	hub := NewHub()
	engine := initRouter(hub)

	// 启动自动注册和注销的协程
	go hub.run()

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
