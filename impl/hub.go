package impl

import (
	"chsir-zy/msg-push-center/impl/message"
	"chsir-zy/msg-push-center/impl/util"
	"sync"
)

/*
 *	从业务处理系统处接受消息，推送给client
 *	处理所有客户端的连接
 */

type Hub struct {
	// 注册客户端
	register chan *Client

	// 注销客户端
	unregister chan *Client

	// uid为map的key，记录每个uid对应的Client
	userClients map[string]*Client

	// 控制每个连接只有一个读写
	sm *sync.RWMutex

	// 消息日志记录器
	logger message.MsgLogger

	// 验证连接过来的用户
	Authenticator util.Authenticator
}

func NewHub() *Hub {
	return &Hub{
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		userClients:   make(map[string]*Client),
		sm:            &sync.RWMutex{},
		logger:        &message.FileMsgLog{},
		Authenticator: &util.JWTAuthenticator{},
	}
}

// 单独协程 专门处理客户端的注册和注销
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// 注册client
			h.sm.RLock()
			h.userClients[client.uid] = client
			h.sm.RUnlock()
		case client := <-h.unregister:
			// 注销client
			h.sm.Lock()
			close(client.send)
			delete(h.userClients, client.uid)
			h.sm.Unlock()
		}
	}
}
