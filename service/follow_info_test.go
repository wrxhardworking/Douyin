package service

import (
	"douyin/dao"
	"fmt"
	"testing"
)

func TestFollowOrCancel(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		fmt.Println(err.Error())
	}
	tests := []struct {
		name     string
		userid   int64
		touserid int64
		action   string
	}{
		{
			name:     "test1",
			userid:   10,
			touserid: 11,
			action:   "1",
		},
		{
			name:     "test2",
			userid:   11,
			touserid: 11,
			action:   "1",
		},
	}

	for _, test := range tests {
		t.Run("测试", func(t *testing.T) {
			err := FollowOrCancel(test.userid, test.touserid, test.action)
			if err != nil {
				t.Errorf("UserRegister ERROR is %v", err)
				return
			}
		})

	}
}

func TestFollowLists(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		fmt.Println(err.Error())
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
			userid: 16,
		},
	}

	for _, test := range tests {
		t.Run("测试", func(t *testing.T) {
			userLists, err := GetFollowLists(test.userid)
			if err != nil {
				t.Errorf("UserRegister ERROR is %v", err)
				return
			}
			fmt.Printf("%#v", userLists)
		})

	}
}

func TestFollowerLists(t *testing.T) {
	err := dao.DbInit()
	if err != nil {
		fmt.Println(err.Error())
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
		t.Run("测试", func(t *testing.T) {
			userLists, err := GetFollowerLists(test.userid)
			if err != nil {
				t.Errorf("UserRegister ERROR is %v", err)
				return
			}
			fmt.Printf("%#v", userLists)
		})

	}
}
