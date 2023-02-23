package dao

import (
	"sync"
)

type Comment struct {
	ID          int64  `gorm:"column:comment_id" json:"id,omitempty"`
	User        User   `gorm:"foreignKey:UserId"     json:"user"`
	UserId      int64  `gorm:"column:user_id"    json:"-"`
	VideoId     int64  `gorm:"column:video_id"   json:"-"`
	CommentText string `gorm:"column:comment_text"   json:"content,omitempty"`
	CreateTime  string `gorm:"column:create_time"    json:"create_date,omitempty"`
	TimeStamp   int64  `gorm:"column:timestamp"      json:"-"`
}

func (Comment) TableName() string {
	return "comment"
}

type CommentDao struct {
}

var commentDao *CommentDao

var commentOnce sync.Once

func GetCommentInstance() *CommentDao {
	commentOnce.Do(func() {
		commentDao = &CommentDao{}
	})
	return commentDao
}

func (CommentDao) AddComment(comment *Comment) error {
	err := db.Create(comment).Error
	if err != nil {
		return err
	}
	return nil
}
func (CommentDao) DeleteCommentById(commentId int64) error {
	//根据主键删除评论
	err := db.Delete(&Comment{}, commentId).Error
	if err != nil {
		return err
	}
	return nil
}

func (CommentDao) QueryCommentByVideoId(videoId int64) ([]Comment, error) {
	var commentLists []Comment
	//按时间的倒叙排序
	err := db.Preload("User").Where("video_id =?", videoId).Order("timestamp desc").Find(&commentLists).Error
	if err != nil {
		return nil, err
	}
	return commentLists, nil
}
