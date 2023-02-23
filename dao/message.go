package dao

import "sync"

type Message struct {
	ID         int64  `gorm:"column:msg_id"       json:"id,omitempty"`
	UserId     int64  `gorm:"column:from_user_id" json:"from_user_id,omitempty"`
	ToUserId   int64  `gorm:"column:to_user_id"   json:"to_user_id,omitempty"`
	Content    string `gorm:"column:content"      json:"content,omitempty"`
	CreateTime int64  `gorm:"column:create_time"  json:"create_time,omitempty"`
}

func (Message) TableName() string {
	return "message"
}

type MessageDao struct {
}

var messageDao *MessageDao
var messageOnce sync.Once

func GetMessageInstance() *MessageDao {
	messageOnce.Do(func() {
		messageDao = &MessageDao{}
	})
	return messageDao
}

// AddMessage 添加消息记录
func (MessageDao) AddMessage(message *Message) error {
	err := db.Create(message).Error
	if err != nil {
		return err
	}
	return nil
}

// QueryMessageLists 消息记录
func (MessageDao) QueryMessageLists(userid, ToUserId int64, preMsgTime int64) ([]Message, error) {
	var MessageLists []Message
	MessageLists = make([]Message, 0, 10)
	//fixme
	err := db.Raw("SELECT * FROM `message` WHERE message.create_time>? and ((to_user_id = ? and message.from_user_id  = ?) or (to_user_id = ? and message.from_user_id = ?)) ORDER BY  create_time", preMsgTime, userid, ToUserId, ToUserId, userid).Scan(&MessageLists).Error
	if err != nil {
		return nil, err
	}
	return MessageLists, nil
}
