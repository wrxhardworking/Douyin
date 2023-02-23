package dao

import (
	"sync"
)

type Like struct {
	UserId  int64 `gorm:"column:user_id"`
	VideoId int64 `gorm:"column:video_id"`
}

func (Like) TableName() string {
	return "like"
}

type LikeDao struct {
}

var likeDao *LikeDao

var likeOnce sync.Once

func GetLikeInstance() *LikeDao {
	likeOnce.Do(func() {
		likeDao = &LikeDao{}
	})
	return likeDao
}

// AddLike 添加映射
func (LikeDao) AddLike(like *Like) error {
	err := db.Create(like).Error
	if err != nil {
		return err
	}
	//要对对应的视频点赞数量进行更新
	//db.Update("")
	return nil
}

// DeleteLike 删除映射
func (LikeDao) DeleteLike(like *Like) error {
	//会删除所有符合条件得到对象 但是无所谓 都行
	err := db.Where("user_id = ? and video_id = ?", like.UserId, like.VideoId).Delete(like).Error
	if err != nil {
		return err
	}
	//要对视频点赞数进行更新
	return nil
}

// QueryLikeByUserid DeleteLike 查找映射 并且返回lists
func (LikeDao) QueryLikeByUserid(userid int64) ([]Video, error) {
	user := &User{}
	err := db.Preload("VideoLieLists").Where("user_id = ?", userid).Preload("VideoLieLists.Author").Find(user).Error
	if err != nil {
		return nil, err
	}
	return user.VideoLieLists, nil
}
