// Package common 返回时引用的对象
package common

import (
	"douyin/cache"
	"douyin/dao"
)

type Response struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// UserCountSearchStrategy 做了缓存后计数的查找策略
func UserCountSearchStrategy(user *dao.User, userid int64) error {
	var err error
	user.TotalFavorite, err = cache.GetUserTotalFavoriteCount(userid)
	if err != nil {
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
	return nil
}
