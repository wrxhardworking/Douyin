package service

import (
	"douyin/common"
	"douyin/dao"
	"errors"
	"fmt"
	"time"
)

// CommentOrDelete 1.进行评论 2.删除评论
func CommentOrDelete(CommentAction string, userid int64, videoId int64, commentId int64, commentText string) (*dao.Comment, error) {
	var err error

	switch CommentAction {
	case "1":
		comment := dao.Comment{
			UserId:      userid,
			VideoId:     videoId,
			CommentText: commentText,
			CreateTime:  time.Now().Format("01-02"),
			TimeStamp:   time.Now().Unix(),
		}
		err = dao.GetCommentInstance().AddComment(&comment)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		err = dao.GetVideoInstance().UpdateCommentCount(videoId, 1)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		err := common.UserCountSearchStrategy(&comment.User, userid)
		if err != nil {
			return nil, err
		}
		return &comment, nil
	case "2":
		err = dao.GetCommentInstance().DeleteCommentById(commentId)
		if err != nil {
			fmt.Println(err.Error())
		}
		err = dao.GetVideoInstance().UpdateCommentCount(videoId, -1)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return nil, nil
	}

	return nil, nil
}

// GetCommentLists 按照评论发布时间的倒叙返回评论列表
func GetCommentLists(videoId int64) ([]dao.Comment, error) {
	var err error
	var CommentLists []dao.Comment

	CommentLists, err = dao.GetCommentInstance().QueryCommentByVideoId(videoId)
	for index := range CommentLists {
		err := common.UserCountSearchStrategy(&CommentLists[index].User, CommentLists[index].User.ID)
		if err != nil {
			return nil, err
		}
	}
	
	if len(CommentLists) == 0 {
		err = errors.New("comment lists not exists")
		return nil, err
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return CommentLists, nil
}
