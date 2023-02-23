package service

import (
	"douyin/cache"
	"douyin/dao"
	"fmt"
	"testing"
)

func TestVideoStream(t *testing.T) {
	cache.RedisPoolInit()
	err := dao.DbInit()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Run("test1", func(t *testing.T) {
		token := "dasdadsadjnaidnd"
		videoInfo, _, err := VideoStream(token)
		if err != nil {
			t.Errorf("%v\n", err)
			return
		}
		fmt.Println(videoInfo)
	})
}

func TestPublishedVideoLists(t *testing.T) {
	cache.RedisPoolInit()
	err := dao.DbInit()
	if err != nil {
		t.Error(err.Error())
		return
	}
	type args struct {
		userid int64
	}
	tests := []struct {
		name string
		args
	}{
		{
			"测试1",
			args{
				1,
			},
		},
		{
			"测试2",
			args{
				100000,
			},
		},
		{
			"测试3",
			args{
				-10,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			videoInfo, err := PublishedVideoLists(test.userid)
			if err != nil {
				t.Errorf("%v\n", err)
				return
			}
			fmt.Printf("%#v", videoInfo)
		})
	}
}
