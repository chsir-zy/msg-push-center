class SocketConnect {
    constructor(url) {
      this.url = url //WebSocket连接地址
      this.ws = null //WebSocket连接对象
      this.heartInterval = 10000 //心跳间隔时间
      this.heartTimeout = 5000 //心跳超时时间
      this.lockReconnect = false //是否禁止重连
      this.heartTimer = null //心跳定时器
      this.serverTimer = null //服务器超时定时器
      this.reconnectCount = 0 //重连次数
      this.maxReconnectCount = 5 //最大重连次数
  
      this.connect()
    }
  
    //WebSocket连接
    connect() {
      this.ws = new WebSocket(this.url)
  
      this.ws.onopen = () => {
        this.reconnectCount = 0 // 重置重连次数
        this.start() // 开启心跳
      }
  
      this.ws.onclose = (event) => {
        console.log('WebSocket closed:', event)
        this.reconnect()
      }
  
      this.ws.onerror = (error) => {
        console.log('WebSocket error:', error)
        this.reconnect()
      }
  
      this.ws.onmessage = (event) => {
        //收到心跳信息则重置心跳，收到其他信息则触发回调
        if (event.data == 'pong') {
          this.start()
        } else {
          var thisDialog = layer.open({
              type: 1
              ,title: '新消息提醒' 
              ,offset: 'rb' 
              // ,id: 'layerDemorb' //防止重复弹出
              ,skin: 'layui-layer-demo'
              ,content: '<div style="padding: 20px 20px;font-size:14px;">'+ event.data +'</div>'
              ,btn: '关闭'
              ,btnAlign: 'rb' //按钮居中
              ,shade: 0 //不显示遮罩
              // ,zIndex: layer.zIndex
              ,success: function(layero, index){
                  layer.setTop(layero);
              }
              ,yes: function(){
                // layer.closeAll();
                layer.close(thisDialog);
              }
          });
        }
      }
    }
  
    //发送信息
    send(message) {
      this.ws.send(message)
      return this
    }
  
    //开启心跳
    start() {
      this.reset()
  
      this.heartTimer = setTimeout(() => {
        this.send("ping")
  
        //5秒钟还没有返回心跳信息，则认为连接断开，关闭WebSocket并重连
        this.serverTimer = setTimeout(() => {
          this.ws.close()
        }, this.heartTimeout)
      }, this.heartInterval)
    }
  
    //重置心跳定时器/服务超时定时器
    reset() {
      this.heartTimer && clearTimeout(this.heartTimer)
  
      this.serverTimer && clearTimeout(this.serverTimer)
    }
  
    //重连
    reconnect() {
      // 设置lockReconnect变量避免重复连接
      if (this.lockReconnect || this.reconnectCount >= this.maxReconnectCount) return
      this.lockReconnect = true
  
      this.reconnectCount++ //重连次数+1
  
      setTimeout(() => {
        this.connect()
        this.lockReconnect = false
      }, 1000 * this.reconnectCount) //重连次数越多，延时越久
    }
  }
  
var url = "ws://localhost:8080";
const ws = new SocketConnect(url)
  