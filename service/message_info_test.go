package service

import (
	"douyin/dao"
	"fmt"
	"testing"
)

func TestSendMessage(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		t.Error(err.Error())
		return
	}
	type args struct {
		userid   int64
		toUserId int64
		content  string
		action   string
	}
	tests := []struct {
		name string
		args
	}{
		{
			name: "测试1",
			args: args{
				userid:   10,
				toUserId: 11,
				content:  "你好啊",
				action:   "1",
			},
		},
		{
			name: "测试2",
			args: args{
				userid:   11,
				toUserId: 10,
				content:  "你也好啊",
				action:   "1",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err = SendMessage(test.userid, test.toUserId, test.action, test.content)
			if err != nil {
				t.Errorf("ERROR is %v", err)
				return
			}
		})
	}
}

func TestGetMessageLists(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		t.Error(err.Error())
		return
	}
	type args struct {
		userid   int64
		toUserId int64
		preTime  int64
	}
	tests := []struct {
		name string
		args
	}{
		{
			name: "测试1",
			args: args{
				userid:   10,
				toUserId: 11,
				preTime:  213123213,
			},
		},
		{
			name: "测试2",
			args: args{
				userid:   10,
				toUserId: -11,
				preTime:  00213123213,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			MessageLists, err := GetMessageLists(test.userid, test.preTime, test.toUserId)
			fmt.Printf("%#v", MessageLists)
			if err != nil {
				t.Errorf("ERROR is %v", err)
				return
			}
		})
	}
}
