package service

import (
	"douyin/cache"
	"douyin/config"
	"douyin/dao"
	"fmt"
	"log"
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

func BenchmarkVideoStream(b *testing.B) {
	var err error
	//1.初始化配置文件

	err = config.ConfInit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = cache.RedisPoolInit()
	err = dao.DbInit()
	if err != nil {
		return
	}
	//重新设置时间
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//这是没有登陆情况下的测试
		//_, _, err := VideoStream("")
		//这是在登陆状态下
		_, _, err := VideoStream("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozMSwiYWRtaW4iOiIxMjM0NTYiLCJleHAiOjE2Nzc2NTAyNDd9.Tjq_DMF5lBJZTmKcWQeL5xpSBiH3J_BwK8fqZfvxSQA")
		if err != nil {
			log.Println(err.Error())
		}
	}
}
