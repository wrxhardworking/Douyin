package handle

import (
	"douyin/common"
	"douyin/dao"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type CommentResponse struct {
	common.Response
	Comment dao.Comment
}
type CommentListsResponse struct {
	common.Response
	CommentLists []dao.Comment `json:"comment_list"`
}

func CommentAction(c *gin.Context) {

	var (
		Comment   *dao.Comment
		videoId   int
		commentId int
		err       error
	)
	//得到userid
	userid := c.MustGet("userid").(int64)
	//得到评论动作
	action := c.Query("action_type")
	//得到评论内容
	commentText := c.Query("comment_text")
	//得到video_id
	videoIdStr := c.Query("video_id")
	if videoIdStr != "" {
		videoId, err = strconv.Atoi(c.Query("video_id"))
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, CommentResponse{
				Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
	}
	//得到评论id
	commentIdStr := c.Query("comment_id")
	if commentIdStr != "" {
		commentId, err = strconv.Atoi(commentIdStr)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, CommentResponse{
				Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
	}
	//进行方法判断
	Comment, err = service.CommentOrDelete(action, userid, int64(videoId), int64(commentId), commentText)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, CommentResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	if Comment == nil {
		c.JSON(http.StatusOK, CommentResponse{
			Response: common.Response{StatusCode: 0, StatusMsg: "successful"},
		})
		return
	}

	c.JSON(http.StatusOK, CommentResponse{
		Response: common.Response{StatusCode: 0, StatusMsg: "successful"},
		Comment:  *Comment,
	})
	return
}

func CommentList(c *gin.Context) {
	var commentLists []dao.Comment

	videoId, err := strconv.Atoi(c.Query("video_id"))

	if err != nil {
		c.JSON(http.StatusOK, CommentListsResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	commentLists, err = service.GetCommentLists(int64(videoId))
	if err != nil {
		c.JSON(http.StatusOK, CommentListsResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, CommentListsResponse{
		Response:     common.Response{StatusCode: 0, StatusMsg: "successful"},
		CommentLists: commentLists,
	})
	return
}
