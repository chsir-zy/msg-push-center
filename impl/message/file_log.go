package message

import (
	"bufio"
	"chsir-zy/msg-push-center/impl/util"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileMsgLog1 struct{}

var _ MsgLogger = &FileMsgLog1{}

type fileLog struct {
	Datetime string `json:"datetime"`
	Msg
}

var filelog fileLog

func (f *FileMsgLog1) Log(msg Msg) error {
	floder := time.Now().Format("2006-01")
	isPath := util.PathExists("./log/" + floder)
	if !isPath {
		return errors.New(floder + " is not exist")
	}
	day := time.Now().Day()
	dayStr := strconv.Itoa(day)
	filename := "./log/" + floder + "/" + dayStr + ".log"

	filelog.Datetime = time.Now().Format(time.RFC3339)
	filelog.Uid = msg.Uid
	filelog.Message = msg.Message

	jsonMsg, err := json.Marshal(filelog)
	if err != nil {
		return err
	}

	// os.WriteFile(filename, jsonMsg, os.ModePerm)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString(string(jsonMsg) + "\n")
	writer.Flush()

	return nil
}

type FileMsgLog struct{}

var _ MsgLogger = &FileMsgLog{}

func (f *FileMsgLog) Log(msg Msg) error {
	now := time.Now()
	day := now.Day()
	floder := now.Format("2006-01")

	logPath := filepath.Join("log", floder)
	if !util.PathExists(logPath) {
		return fmt.Errorf("log path : %s, is not exist", logPath)
	}

	fileName := filepath.Join(logPath, strconv.Itoa(day)+".log")
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	filelog.Datetime = time.Now().Format(time.RFC3339)
	filelog.Uid = msg.Uid
	filelog.Message = msg.Message

	jsonMsg, err := json.Marshal(filelog)
	if err != nil {
		return err
	}

	logger := log.Default()
	logger.SetOutput(file)
	logger.SetPrefix("[send message log]")
	logger.Print(string(jsonMsg))

	return nil
}
