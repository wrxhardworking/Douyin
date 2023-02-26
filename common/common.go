// Package common 返回时引用的对象
package common

import (
	"douyin/cache"
	"douyin/dao"
	"log"
)

type Response struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// UserCountSearchStrategy todo 进行缓存
// UserCountSearchStrategy 做了缓存后计数的查找策略
func UserCountSearchStrategy(user *dao.User, userid int64) error {
	var err error
	user.TotalFavorite, err = cache.GetUserTotalFavoriteCount(userid)
	//是否在redis中的标志
	var isExists = true
	if err != nil {
		if err.Error() == "is not exists" {
			//说明redis中没找到
			isExists = false
		}
		user.TotalFavorite, err = dao.GetUserInstance().QueryTotalFavorite(userid)
		if err != nil {
			return err
		}
	}
	user.WorkCount, err = cache.GetUserWorkCount(userid)
	if err != nil {
		user.WorkCount, err = dao.GetUserInstance().QueryWorkCount(userid)
		if err != nil {
			return err
		}
	}
	user.FavoriteCount, err = cache.GetUserFavoriteCount(userid)
	if err != nil {
		user.FavoriteCount, err = dao.GetUserInstance().QueryFavoriteCount(userid)
		if err != nil {
			return err
		}
	}

	if !isExists {
		log.Println("caching")
		//说明不缓存中不存在 重新进行设置
		err := cache.SetUserCount(userid, map[string]int64{
			"followCount":   0,
			"followerCount": 0,
			"workCount":     user.WorkCount,
			"favoriteCount": user.FavoriteCount,
			"totalFavorite": user.TotalFavorite})
		if err != nil {
			log.Println("redis 服务器出问题")
		}
	}
	return nil
}
