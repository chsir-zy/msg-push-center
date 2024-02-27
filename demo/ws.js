var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxIn0.3shGQJ4Um1To5hFQ7vQTVLpRKWGxZWDKLX29j53rH1A"
// const socket = new WebSocket("ws://localhost:8080/ws/connect?token=" + token);
// socket.addEventListener("message", function (event) {
//     if (event.data != 'connection success') {
//         layer.open({
//             type: 1
//             ,title: '新消息提醒' 
//             ,offset: 'rb' 
//             ,id: 'layerDemorb' //防止重复弹出
//             ,content: '<div style="padding: 20px 20px;font-size:14px;">'+ event.data +'</div>'
//             ,btn: '关闭'
//             ,btnAlign: 'rb' //按钮居中
//             ,shade: 0 //不显示遮罩
//             ,zIndex: layer.zIndex
//             ,success: function(layero, index){
//                 layer.setTop(layero);
//             }
//             ,yes: function(){
//               layer.closeAll();
//             }
//         });
//     }
    
// });

// WebSocket
var ws = null;
// WebSocket连接地址
var ws_url = "ws://localhost:8080/ws/connect?token=" + token

// 创建WebSocket
ws_create(ws_url);

// 创建WebSocket
function ws_create(url) {
     try{
        // 判断是否支持 WebSocket
        if('WebSocket' in window){
            // 连接WebSocket
            ws = new WebSocket(url);
            // 初始化WebSocket事件(WebSocket对象, WebSocket连接地址)
            ws_event(ws, url);
        }
    }catch(e){
        // 重新连接WebSocket
        ws_recontent(url);
    }
}

// WebSocket 事件创建
function ws_event(ws, url) {
    ws.onopen = function (event) {
        // 心跳检测重置
        ws_heartCheck_func();
        console.log("WebSocket已连接");
    };

    ws.onclose = function (event) {
        // 重新连接WebSocket
        // ws_recontent(url);
        console.log("WebSocket连接已关闭");
    };

    ws.onerror = function (event) {
        console.log("WebSocket错误：", event);
    };

    ws.onmessage = function (event) {
        console.log(event.data)
        console.log(layer)
        // 只要有数据，那就说明连接正常
        // 处理数据，只处理非心跳检测的数据
        if (event.data != 'check') {
          // 处理数据
            if (event.data != 'connection success') {
                layer.open({
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
                      layer.closeAll();
                    }
                });

                layer.alert(event.data);

                layer.msg(event.data)
            }
        }
    };
}

// 重新连接websocker(WebSocket连接地址)
function ws_recontent(url) {
    // 延迟避免请求过多
    // setTimeout(function () {
    //     ws_create(url);
    // }, 2000);

    ws_create(url);
}

// 监听窗口关闭事件，当窗口关闭时，主动去关闭websocket连接，防止连接还没断开就关闭窗口，这样服务端会抛异常。
window.onbeforeunload = function() {
    console.log("window.onbeforeunload")
    ws.close();
} 

// WebSocket心跳检测
var ws_heartCheck = {
    timeout: 5000,          // 5秒一次心跳
    timeoutObj: null,       // 执行心跳的定时器
    serverTimeoutObj: null, // 服务器超时定时器
    reset: function(){      // 重置方法
        clearTimeout(this.timeoutObj);
        clearTimeout(this.serverTimeoutObj);
        return this;
    },
    start: function(){      // 启动方法
        var self = this;
        this.timeoutObj = setTimeout(function(){
            // 这里发送一个心跳信息，后端收到后，返回一个消息，在onmessage拿到返回的心跳（信息）就说明连接正常
            ws.send("check");
            // 如果超过一定时间还没重置，说明后端主动断开了
            // self.serverTimeoutObj = setTimeout(function(){
            //     // 如果onclose会执行reconnect，我们执行ws.close()就行了.如果直接执行reconnect 会触发onclose导致重连两次
            //     ws.close();
            // }, self.timeout);
        }, this.timeout);

        
    }
}

var heartBeat = {
    type: "ping",
    timestamp: new Date().getTime()
}
function ws_heartCheck_func(){
    setInterval(function(){
        ws.send(JSON.stringify(heartBeat));
    }, 5000)
}