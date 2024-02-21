# msg-push-center
    websocket message push center

# 说明
## 客户端连接
    执行go run main.go 启动服务，监听在8080端口，浏览器 访问/ws/connect 连接websocket服务器。

## 推送消息
    业务系统调用 /send 将消息发送给hub hub推送给浏览器
