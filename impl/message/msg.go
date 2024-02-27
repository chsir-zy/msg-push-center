package message

// 消息推送的内容模板
type Msg struct {
	Uid     string
	Message string
}

// 定义一个消息记录的接口
//
// 记录消息需要实现这个接口
// 可以把日志写入控制台、文件、mysql数据库等
type MsgLogger interface {
	Log(msg Msg) error
}
