package message

import (
	"chsir-zy/msg-push-center/impl/util"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileMsgLog struct{}

var _ MsgLogger = &FileMsgLog{}

type fileLog struct {
	Datetime string `json:"datetime"`
	Msg
}

var filelog fileLog

func (f *FileMsgLog) Log(msg Msg) error {
	floder := time.Now().Format("2006-01")
	isPath := util.PathExists("./log/" + floder)
	if !isPath {
		return errors.New(floder + " is not exist")
	}
	day := time.Now().Day()
	dayStr := strconv.Itoa(day)
	fmt.Println("day:", day)
	filename := "./log/" + floder + "/" + dayStr + ".log"

	fmt.Println("filename:", filename)

	filelog.Datetime = time.Now().Format(time.RFC3339)
	filelog.Uid = msg.Uid
	filelog.Message = msg.Message

	jsonMsg, err := json.Marshal(filelog)
	if err != nil {
		return err
	}

	os.WriteFile(filename, jsonMsg, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
