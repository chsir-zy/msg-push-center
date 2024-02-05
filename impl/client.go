package impl

import (
	"chsir-zy/msg-push-center/impl/message"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const CLIENT_BUFFER_SIZE = 16 // 客户端发送消息的缓存大小
const WRITER_TIMEOUT = 10     //websocket 写入的超时时间

/*
 *	处理客户端(可以认为是浏览器)的连接信息
 *
 */
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan message.Msg // 发送通道，往客户端消息的通道

	uid string // 用户id
}

// 从客户端读消息
//
// 我们这里只处理服务端往浏览器推送消息 暂时将浏览器发送过来的消息丢弃
func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
		close(c.send)
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
	}

}

// 往浏览器发送消息
func (c *Client) writerPump() {
	defer func() { // 关闭底层连接
		c.conn.Close()
	}()

	for {
		// 当send里面有数据时，就往客户端发送  没有消息的话会阻塞在这里
		msgLog, ok := <-c.send

		// 设置超时时间
		c.conn.SetWriteDeadline(time.Now().Add(time.Second * WRITER_TIMEOUT))

		if !ok { // 如果send通道关闭了 例如在注销client的时候会关闭通道
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, []byte(msgLog.Msg))
		if err != nil {
			log.Println(err)
			return
		}

		c.hub.logger.Log(msgLog)

	}

}

var upgrater = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 这里可以校验客户端  比如是否跨域等
		return true
	},
}

// websocket 连接
func ServeWs(c *gin.Context, hub *Hub) {
	conn, err := upgrater.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// 返回错误消息给客户端
		c.Writer.WriteHeader(http.StatusBadRequest)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("websocket connection error: %s", err)))
		conn.Close()
		return
	}

	// TODO  这里我们实现一个JWT认证
	uid, err := hub.Authenticator.Authenticate(c)
	fmt.Println("uid---------", uid)

	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("authenticate error: %s", err)))
		conn.Close()
		return
	}

	// 告诉客户端连接成功
	conn.WriteMessage(websocket.TextMessage, []byte("connection success"))

	// 创建一个和浏览器连接的客户端
	client := &Client{conn: conn, hub: hub, send: make(chan message.Msg, CLIENT_BUFFER_SIZE), uid: uid}

	// 向hub注册一个客户端
	client.hub.register <- client

	// 启动两个协程 分别进行读和写
	go client.readPump()
	go client.writerPump()
}

/*
 * 业务处理中心往client推送消息
 * client 有个协程writerPump 向客户端发送消息，当收到业务中心的消息，就会往客户端发送
 */
func Send(c *gin.Context, hub *Hub) {
	// 业务中心需要传uid  hub才能知道往哪个client推送消息
	uid := c.PostForm("uid")
	if uid == "" {
		// c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "uid is need",
		})
		return
	}

	// 获取 client
	hub.sm.RLock()
	client, ok := hub.userClients[uid]
	hub.sm.RUnlock()
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "client is not exist,please check uid",
		})
		return
	}

	sendMsg := c.PostForm("msg")

	fmt.Println(sendMsg)
	msgLog := message.Msg{
		Uid: uid,
		Msg: sendMsg,
	}

	client.send <- msgLog
}
