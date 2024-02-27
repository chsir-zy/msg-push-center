package message

import (
	"chsir-zy/msg-push-center/config"
	"strconv"

	"gorm.io/gorm"
)

/*
 *  发送消息记录到mysql
 */
// var db *gorm.DB

// // 初始化
// func init() {
// 	var err error
// 	// dsn := "root:123456@tcp(127.0.0.1:3306)/msg_push_center?charset=utf8mb4&parseTime=True&loc=Local"
// 	var configMysql = config.CONFIG.Mysql
// 	dsn := configMysql.Username + ":" + configMysql.Password + "@tcp(" + configMysql.Host + ":" + configMysql.Port + ")/msg_push_center?charset=utf8mb4&parseTime=True&loc=Local"
// 	fmt.Println(dsn)
// 	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic(fmt.Sprintf("mysql connect error: %s", err))
// 	}
// }

type messageLog struct {
	gorm.Model
	Uid     uint32 `gorm:"uid"`
	Message string `gorm:"message"`
}

// 保住 MysqlMsgLog 实现了MsgLogger接口
var _ MsgLogger = &MysqlMsgLog{}

type MysqlMsgLog struct{}

func (mysql *MysqlMsgLog) Log(msg Msg) error {
	config.GORM_DB.AutoMigrate(&messageLog{})

	uid, _ := strconv.Atoi(msg.Uid)
	err := config.GORM_DB.Create(&messageLog{
		Uid:     uint32(uid),
		Message: msg.Message,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
