package service

import (
	"douyin/cache"
	"douyin/common"
	"douyin/dao"
	"douyin/utls"
	"errors"
	"fmt"
	"time"
)

func VideoStream(token string) ([]dao.Video, int64, error) {
	VideoLists, err := dao.GetVideoInstance().QueryVideo()
	//如果token不为空
	if token != "" {
		_, clime, _ := utls.ParseToken(token)
		var userLists []dao.User
		var likeVideoLists []dao.Video
		userLists, err = dao.GetFollowInstance().QueryFollowLists(clime.UserId)
		if err != nil {
			return VideoLists, 0, err
		}
		likeVideoLists, err = dao.GetLikeInstance().QueryLikeByUserid(clime.UserId)
		if err != nil {
			return VideoLists, 0, err
		}
		//fixme 缓存
		//设置 用户是否关注 视频是否点赞
		for index := range VideoLists {
			if cache.IsUserRelation(clime.UserId, VideoLists[index].Author.ID) {
				VideoLists[index].Author.IsFollow = true
			} else {
				for _, user := range userLists {
					if VideoLists[index].Author.ID == user.ID {
						VideoLists[index].Author.IsFollow = true
					}
				}
			}

			if cache.IsUserVideoRelation(clime.UserId, VideoLists[index].ID) {
				VideoLists[index].IsFavorite = true
			} else {
				for _, likeVideo := range likeVideoLists {
					if likeVideo.ID == VideoLists[index].ID {
						VideoLists[index].IsFavorite = true
					}
				}
			}

			err := common.UserCountSearchStrategy(&VideoLists[index].Author, VideoLists[index].Author.ID)
			if err != nil {
				return nil, 0, err
			}
		}
	}

	if len(VideoLists) == 0 {
		err = errors.New("video lists not exists")
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, 0, err
	}
	//得到最早的时间返回过去
	nextTime := VideoLists[len(VideoLists)-1].TimeStamp
	return VideoLists, nextTime, nil

}

// PublishedVideoLists 已发布视频的列表
func PublishedVideoLists(userid int64) ([]dao.Video, error) {
	var (
		//fixme 不查询数据库
		//LikeVideoLists []dao.Video
		VideoLists []dao.Video
		err        error
	)
	VideoLists, err = dao.GetVideoInstance().QueryVideoByUserId(userid)
	//fixme
	if len(VideoLists) == 0 {
		err = errors.New("video lists not exists")
		return nil, err
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	//fixme
	//LikeVideoLists, err = dao.GetLikeInstance().QueryLikeByUserid(userid)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	//fixme 缓存
	for index1 := range VideoLists {
		if cache.IsUserVideoRelation(userid, VideoLists[index1].ID) {
			VideoLists[index1].IsFavorite = true
		}
		//for index2 := range LikeVideoLists {
		//	if VideoLists[index1].ID == LikeVideoLists[index2].ID {
		//		VideoLists[index1].IsFavorite = true
		//	}
		//}
		err := common.UserCountSearchStrategy(&VideoLists[index1].Author, userid)
		if err != nil {
			return nil, err
		}
	}

	if len(VideoLists) >= 30 {
		//限制只能返回最多三十条视频
		return VideoLists[0:30], nil
	}

	return VideoLists, nil
}

// PublishVideo 投稿视频 并且更新作品数量
func PublishVideo(userid int64, playUrl, coverUrl, title string) error {
	err := dao.GetVideoInstance().AddVideo(&dao.Video{
		UserId:        userid,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
		TimeStamp:     time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	err = cache.IncrByUserWorkCount(userid)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
