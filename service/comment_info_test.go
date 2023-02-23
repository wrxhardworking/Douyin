package service

import (
	"douyin/cache"
	"douyin/dao"
	"fmt"
	"testing"
)

func TestCommentOrDelete(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		fmt.Println(err.Error())
	}
	cache.RedisPoolInit()
	tests := []struct {
		name        string
		videoId     int64
		userId      int64
		action      string
		commentText string
		commentId   int64
	}{
		{
			name:        "test1",
			videoId:     1,
			userId:      10,
			action:      "1",
			commentText: "test_comment1",
			commentId:   13,
		},
		{
			name:        "test2",
			videoId:     2,
			userId:      10,
			action:      "1",
			commentText: "test_comment2",
			commentId:   11,
		},
		{
			name:        "test3",
			videoId:     3,
			userId:      12,
			action:      "1",
			commentText: "test_comment3",
			commentId:   12,
		},
	}

	for _, test := range tests {
		t.Run("测试", func(t *testing.T) {
			comment, err := CommentOrDelete(test.action, test.userId, test.videoId, test.commentId, test.commentText)
			if err != nil {
				t.Errorf("UserRegister ERROR is %v", err)
				return
			}
			fmt.Printf("%#v", comment)
		})
	}
}

func TestGetCommentLists(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		fmt.Println(err.Error())
	}
	cache.RedisPoolInit()
	tests := []struct {
		name    string
		videoId int64
	}{
		{
			name:    "test1",
			videoId: 1,
		},
		{
			name:    "test2",
			videoId: 2,
		},
		{
			name:    "test3",
			videoId: 3,
		},
	}
	for _, test := range tests {
		t.Run("测试", func(t *testing.T) {
			commentLists, err := GetCommentLists(test.videoId)
			if err != nil {
				t.Errorf("UserRegister ERROR is %v", err)
				return
			}
			fmt.Printf("%#v", commentLists)
		})

	}
}

func TestGetFriendLists(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		fmt.Println(err.Error())
	}
	cache.RedisPoolInit()
	tests := []struct {
		name   string
		userId int64
	}{
		{
			name:   "test1",
			userId: 10,
		},
		//{
		//	name:    "test2",
		//	videoId: 2,
		//},
		//{
		//	name:    "test3",
		//	videoId: 3,
		//},
	}
	for _, test := range tests {
		t.Run("测试", func(t *testing.T) {
			commentLists, err := GetFriendLists(test.userId)
			if err != nil {
				t.Errorf("%#v", err)
				return
			}
			fmt.Printf("%#v", commentLists)
		})

	}
}
