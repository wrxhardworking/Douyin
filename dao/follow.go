package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type Follow struct {
	//关注者
	FollowId int64 `gorm:"column:follow_id"`
	//
	FollowedId int64 `gorm:"column:followed_id"`
}

func (Follow) TableName() string {
	return "follow"
}

type FollowDao struct {
}

var followDao *FollowDao

var FollowOnce sync.Once

func GetFollowInstance() *FollowDao {
	FollowOnce.Do(func() {
		followDao = &FollowDao{}
	})
	return followDao
}

func (FollowDao) QueryAllFollow() ([]Follow, error) {

	var FollowLists []Follow
	err := db.Find(&FollowLists).Error
	if err != nil {
		return nil, err
	}
	return FollowLists, err
}

// AddFollow 添加关注映射
func (FollowDao) AddFollow(follow *Follow) error {
	err := db.Create(follow).Error
	if err != nil {
		return err
	}
	//进行关注数量的更新
	return nil
}

// DeleteFollow 删除关注映射
func (FollowDao) DeleteFollow(follow *Follow) error {
	err := db.Where("follow_id = ? and followed_id = ?", follow.FollowId, follow.FollowedId).Delete(follow).Error
	if err != nil {
		return err
	}
	//进行关注数量的更新
	return nil
}

// QueryFollowLists 获得关注列表
func (FollowDao) QueryFollowLists(userid int64) ([]User, error) {
	var userLists []User
	err := db.Raw("SELECT * FROM  `user` WHERE user.user_id IN ( SELECT follow.followed_id FROM follow WHERE follow.follow_id = ? )", userid).Scan(&userLists).Error
	if err != nil {
		return nil, err
	}
	return userLists, nil
}

// QueryFollowerLists 获得粉丝列表
func (FollowDao) QueryFollowerLists(userid int64) ([]User, error) {
	var userLists []User
	err := db.Raw("SELECT * FROM  `user` WHERE user.user_id IN ( SELECT follow.follow_id FROM follow WHERE follow.followed_id = ? )", userid).Scan(&userLists).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", userLists)
	return userLists, nil
}

func (FollowDao) QueryEachFollow(userid int64) ([]User, error) {
	var userLists []User
	err := db.Raw("SELECT * FROM  `user` WHERE user.user_id != ?	 and user.user_id in \n(SELECT DISTINCT follow.followed_id FROM  follow \njoin\n(SELECT follow.follow_id FROM follow WHERE follow.followed_id = ?) a \non\na.follow_id = follow.followed_id) ", userid, userid).Scan(&userLists).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", userLists)
	return userLists, nil
}

// Exists 判断某个关注关系是否存在
func (FollowDao) Exists(followId, followedId int64) bool {
	var follow Follow
	if err := db.Where("follow_id = ?", followId).Where("followed_id = ?", followedId).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false // 记录不存在
		}
		return true // 查询过程中发生了其他错误，记录存在
	}
	return true // 记录存在
}
