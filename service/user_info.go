// Package service 不要依赖于上Handle的任何东西
package service

import (
	"douyin/cache"
	"douyin/common"
	"douyin/dao"
	"douyin/utls"
	"errors"
)

// UserRegisInfo 返回给handle层的数据
type UserRegisInfo struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID int64  `json:"user_id"` // 用户id
}

// UserRegister 用户注册
func UserRegister(username string, password string) (*UserRegisInfo, error) {

	var err error
	var info UserRegisInfo
	var token string
	var user *dao.User

	//进行加密
	password = utls.Md5Encryption(password)

	user, err = dao.GetUserInstance().QueryUserByName(username)

	if user.ID != 0 {
		err = errors.New("user existed")
	}
	if err != nil {
		return nil, err
	}

	user = &dao.User{
		Name:            username,
		FollowCount:     0,
		FollowerCount:   0,
		Password:        password,
		Signature:       "随便玩玩罢了",
		Avatar:          "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		BackGroundImage: "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
	}

	//用户不存在则创建用户
	err = dao.GetUserInstance().AddUser(user)
	if err != nil {
		return nil, err
	}
	//添加计数缓存
	err = cache.SetUserCount(user.ID)
	if err != nil {
		return nil, err
	}

	token, err = utls.GenerateToken(username, user.ID)
	//fmt.Printf("token is %s", token)
	if err != nil {
		return nil, err
	}

	info.Token = token
	info.UserID = user.ID
	return &info, nil
}

func UserLogin(username string, password string) (*UserRegisInfo, error) {
	var err error
	var token string
	var user *dao.User
	//进行md5加密
	password = utls.Md5Encryption(password)
	user, err = dao.GetUserInstance().QueryUserByName(username)
	//判断用户是否存在
	if user.ID == 0 {
		err = errors.New("user not exists")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	//验证密码是否正确
	if password != user.Password {
		err = errors.New("password is wrong")
		return nil, err
	}
	//生成token
	token, err = utls.GenerateToken(username, user.ID)
	//成功返回
	return &UserRegisInfo{Token: token, UserID: user.ID}, nil
}

func GetUserInfo(userid int64) (*dao.User, error) {
	//返回给Handle的user
	var err error
	//Dao数据层user
	var user *dao.User
	user, err = dao.GetUserInstance().QueryUserByID(userid)
	if err != nil {
		return nil, err
	}
	//对三个count进行查找赋值
	err = common.UserCountSearchStrategy(user, userid)
	return user, nil
}
