package service

import (
	"douyin/dao"
	"errors"
	"time"
)

// SendMessage 存储消息记录
func SendMessage(userid, toUserId int64, action, content string) error {
	var err error
	switch action {
	case "1":
		err = dao.GetMessageInstance().AddMessage(&dao.Message{
			UserId:   userid,
			ToUserId: toUserId,
			Content:  content,
			//指定日期格式
			CreateTime: time.Now().Unix(),
		})
		if err != nil {
			return err
		}
	default:
		err = errors.New("action is not valid")
		return err
	}
	return nil
}

// GetMessageLists 获得消息记录
func GetMessageLists(userid, preMsgTime, toUserId int64) ([]dao.Message, error) {
	var err error
	var MessageLists []dao.Message
	MessageLists, err = dao.GetMessageInstance().QueryMessageLists(userid, toUserId, preMsgTime)
	if err != nil {
		return nil, err
	}
	if len(MessageLists) > 10 {
		//设置最多只能返回十条消息
		return MessageLists[0:10], nil
	}
	return MessageLists, nil
}
