package service

import (
	"douyin/cache"
	"douyin/dao"
	"fmt"
	"testing"
)

func TestThumbUpOrCancel(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = cache.RedisPoolInit()
	if err != nil {
		return
	}
	tests := []struct {
		name    string
		userid  int64
		videoId int64
		action  string
	}{
		{
			name:    "test1",
			userid:  10,
			videoId: 1,
			action:  "1",
		},
		{

			name:    "test2",
			userid:  11,
			videoId: 2,
			action:  "1",
		},
		{
			name:    "test3",
			userid:  10,
			videoId: 2,
			action:  "1",
		},
	}

	for _, test := range tests {
		t.Run("测试", func(t *testing.T) {
			err := ThumbUpOrCancel(test.action, test.userid, test.videoId)
			if err != nil {
				t.Errorf("UserRegister ERROR is %v", err)
				return
			}
		})
	}
}

func TestGetLikeLists(t *testing.T) {
	err := cache.RedisPoolInit()
	if err != nil {
		return
	}
	err = dao.DbInit()
	if err != nil {
		t.Error(err.Error())
		return
	}

	tests := []struct {
		name   string
		userid int64
	}{
		{
			name:   "test1",
			userid: 10,
		},
		{
			name:   "test2",
			userid: 11,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			likeLists, err := GetLikeLists(test.userid)
			if err != nil {
				t.Errorf("UserRegister ERROR is %v", err)
				return
			}
			fmt.Printf("%#v", likeLists)
		})
	}
}
