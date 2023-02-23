// Package dao 就是面向底层数据 不管其他的
package dao

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type User struct {
	ID              int64   `gorm:"column:user_id" json:"id"`
	Name            string  `gorm:"column:name"    json:"name"`
	FollowCount     int64   `gorm:"column:follow_count"   json:"follow_count"`
	FollowerCount   int64   `gorm:"column:follower_count" json:"follower_count"`
	Password        string  `gorm:"column:password"       json:"-"`
	IsFollow        bool    `gorm:"-"                     json:"is_follow" `
	Avatar          string  `gorm:"column:avatar"         json:"avatar"`
	BackGroundImage string  `gorm:"column:background_image"        json:"background_image"`
	Signature       string  `gorm:"column:signature"               json:"signature"`
	TotalFavorite   int64   `gorm:"-"               json:"total_favorited"`
	WorkCount       int64   `gorm:"-"               json:"work_count"`
	FavoriteCount   int64   `gorm:"-"               json:"favorite_count"`
	VideoLieLists   []Video `gorm:"many2many:like;" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

type UserDao struct {
}

var userDao *UserDao

// 单例 接口 表示只创建一次对象
var userOnce sync.Once

func GetUserInstance() *UserDao {
	//创建单例 类比为cpp中的局部静态变量
	userOnce.Do(func() {
		userDao = &UserDao{}
	})
	return userDao
}

// QueryUserByID 通过id查找user
func (UserDao) QueryUserByID(userID int64) (*User, error) {
	var user User
	err := db.Where("user_id = ?", userID).Find(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	//逃逸分析
	return &user, nil
}

// QueryUserByName 通过名字找到user
func (UserDao) QueryUserByName(name string) (*User, error) {

	var user User
	err := db.Where("name = ?", name).Find(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	//逃逸分析
	return &user, nil
}

// AddUser 添加user
func (UserDao) AddUser(user *User) error {
	res := db.Create(user)
	err := res.Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateFollowCount 更新关注数量
func (UserDao) UpdateFollowCount(userId, count int64) error {
	err := db.Model(&User{}).Where("user_id = ?", userId).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", count)).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateFollowerCount 更新粉丝数量
func (UserDao) UpdateFollowerCount(userId, count int64) error {
	err := db.Model(&User{}).Where("user_id = ?", userId).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", count)).Error
	if err != nil {
		return err
	}
	return nil
}

// QueryWorkCount 获取作品数量
func (UserDao) QueryWorkCount(userid int64) (int64, error) {
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM video WHERE video.user_id = ?", userid).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// QueryFavoriteCount 获取点赞的总数量
func (UserDao) QueryFavoriteCount(userid int64) (int64, error) {
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM `like` WHERE `like`.user_id = ?", userid).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// QueryTotalFavorite  获取收获的赞总数
func (UserDao) QueryTotalFavorite(userid int64) (int64, error) {
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM `like` WHERE `like`.video_id in (SELECT video_id FROM video WHERE video.user_id  = ? )", userid).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
