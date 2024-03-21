package impl

import (
	"chsir-zy/msg-push-center/config"
	"chsir-zy/msg-push-center/impl/util"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewServer() {
	config.LoadConfig()
	// config.GORM_DB = config.GormMysql()
	config.JWT_KEY = config.CONFIG.Jwt.Key

	hub := NewHub()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine = initRouter(hub, engine)

	// 启动自动注册和注销的协程
	go hub.run()

	err := engine.Run(":8089") // 默认运行在本机的8089端口
	if err != nil {
		fmt.Println(err)
	}
}

func initRouter(hub *Hub, router *gin.Engine) *gin.Engine {
	ws := router.Group("/ws")
	{
		ws.GET("/connect", func(ctx *gin.Context) {
			ServeWs(ctx, hub)
		})

		ws.POST("/send", func(ctx *gin.Context) {
			Send(ctx, hub)
		})

		ws.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "pong")
		})
	}

	{
		router.GET("/token", func(ctx *gin.Context) {
			token := util.GenToken(ctx)
			ctx.JSON(200, gin.H{
				"token": token,
			})
		})

		router.GET("/getOnlineUser", func(ctx *gin.Context) {
			uids := GetOnlineUid(*hub)
			ctx.JSON(200, gin.H{
				"count": len(uids),
				"uids":  uids,
			})
		})
	}

	router.LoadHTMLGlob("templates/*")
	router.GET("/test", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "test.tmpl", gin.H{
			"title": "test",
		})
	})

	return router
}

func NewHttpsServer() {
	config.LoadConfig()
	// config.GORM_DB = config.GormMysql()
	config.JWT_KEY = config.CONFIG.Jwt.Key

	hub := NewHub()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine = initRouter(hub, engine)

	cert, _ := tls.LoadX509KeyPair("rootCa.pem", "rootCa.key")
	server := &http.Server{
		Addr:    ":8089",
		Handler: engine,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
			Certificates:       []tls.Certificate{cert},
		},
	}

	err := server.ListenAndServeTLS("rootCa.pem", "rootCa.key")
	if err != nil {
		fmt.Println(err)
	}
}
