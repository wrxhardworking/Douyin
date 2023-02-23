package service

import (
	"douyin/cache"
	"douyin/common"
	"douyin/dao"
	"errors"
	"fmt"
	"log"
)

func FollowOrCancel(userid int64, toUserid int64, action string) error {
	var err error
	switch action {
	//关注
	case "1":
		err = dao.GetFollowInstance().AddFollow(&dao.Follow{
			FollowId:   userid,
			FollowedId: toUserid,
		})
		if err != nil {
			return err
		}

		err = dao.GetUserInstance().UpdateFollowCount(userid, 1)
		if err != nil {
			return err
		}

		err = dao.GetUserInstance().UpdateFollowerCount(toUserid, 1)
		if err != nil {
			return err
		}
		err = cache.SetUserRelation(userid, toUserid)
		if err != nil {
			//fixme 这里不能直接返回错误只是缓存而已
			log.Println(err.Error())
		}
	//取消关注
	case "2":
		err = dao.GetFollowInstance().DeleteFollow(&dao.Follow{
			FollowId:   userid,
			FollowedId: toUserid,
		})
		if err != nil {
			return err
		}
		err = dao.GetUserInstance().UpdateFollowCount(userid, -1)
		if err != nil {
			return err
		}
		err = dao.GetUserInstance().UpdateFollowerCount(toUserid, -1)
		if err != nil {
			return err
		}
		err = cache.DeleteUserRelation(userid, toUserid)
		if err != nil {
			//fixme 这里不能直接返回错误只是缓存而已
			log.Println(err.Error())
		}
	}
	return nil
}

// GetFollowLists   获得关注者的列表
func GetFollowLists(userid int64) ([]dao.User, error) {
	var err error
	var UserLists []dao.User
	UserLists, err = dao.GetFollowInstance().QueryFollowLists(userid)
	for index := range UserLists {
		err := common.UserCountSearchStrategy(&UserLists[index], UserLists[index].ID)
		if err != nil {
			return nil, err
		}
		UserLists[index].IsFollow = true
	}
	if len(UserLists) == 0 {
		err = errors.New("users not exists")
		return nil, err
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return UserLists, nil
}

func GetFollowerLists(userid int64) ([]dao.User, error) {
	var err error
	var UserLists []dao.User
	var FollowUserLists []dao.User
	UserLists, err = dao.GetFollowInstance().QueryFollowerLists(userid)
	FollowUserLists, err = dao.GetFollowInstance().QueryFollowLists(userid)
	//fixme redis缓存
	for index1 := range UserLists {
		if cache.IsUserRelation(userid, UserLists[index1].ID) {
			UserLists[index1].IsFollow = true
		} else {
			for index2 := range FollowUserLists {
				if UserLists[index1].ID == FollowUserLists[index2].ID {
					UserLists[index1].IsFollow = true
				}
			}
		}
		err := common.UserCountSearchStrategy(&UserLists[index1], UserLists[index1].ID)
		if err != nil {
			return nil, err
		}
	}
	if len(UserLists) == 0 {
		err = errors.New("users not exists")
		return nil, err
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return UserLists, nil
}

func GetFriendLists(userid int64) ([]dao.User, error) {
	var err error
	var UserLists []dao.User
	UserLists, err = dao.GetFollowInstance().QueryEachFollow(userid)
	//将状态设置为已经关注
	for index := range UserLists {
		UserLists[index].IsFollow = true
		err := common.UserCountSearchStrategy(&UserLists[index], UserLists[index].ID)
		if err != nil {
			return nil, err
		}
	}
	if len(UserLists) == 0 {
		err = errors.New("users not exists")
		return nil, err
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return UserLists, nil
}
