package handle

import (
	"douyin/common"
	"douyin/config"
	"douyin/dao"
	"douyin/service"
	"douyin/utls"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type PublishedResponse struct {
	common.Response
	VideoLists []dao.Video `json:"video_list,omitempty"`
}

func PublishList(c *gin.Context) {
	//返回所有视频信息
	userid, _ := strconv.Atoi(c.Query("user_id"))

	videoLists, err := service.PublishedVideoLists(int64(userid))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	} else {
		c.JSON(http.StatusOK, FeedResponse{
			Response:   common.Response{StatusCode: 0, StatusMsg: "successful"},
			VideoLists: videoLists,
		})
		return
	}
}

func Publish(c *gin.Context) {
	//返回所有视频信息
	title := c.PostForm("title")
	//得到token 获取userid
	userid := c.MustGet("userid").(int64)
	//获取文件
	file, err := c.FormFile("data")

	if err != nil {
		//得到的文件错误
		log.Println(err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	//保证了文件名的随机性 得到文件
	filename := filepath.Base(file.Filename)
	//得到文件的后缀名
	finalName := fmt.Sprintf("%d_%s", userid, filename)
	//选择路径
	//映射成了对外的static文件
	saveFilePath := filepath.Join("./public/", finalName)
	//保存文件到对应的路径
	err = c.SaveUploadedFile(file, saveFilePath)
	if err != nil {
		//得到的文件错误
		log.Println(err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	//生成对应的快照
	savePagePath := "./public/" + finalName
	err = utls.GenerateSnapshot(saveFilePath, savePagePath, 1)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	//192.168.0.103
	//本机ip地址
	//静态资源的地址
	playUrl := "http://" + config.C.Resouece.Ipaddress + ":" + config.C.Resouece.Port + "/" + "static/" + finalName
	coverUrl := "http://" + config.C.Resouece.Ipaddress + ":" + config.C.Resouece.Port + "/" + "static/" + finalName + ".png"

	//发布视频
	err = service.PublishVideo(userid, playUrl, coverUrl, title)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response: common.Response{StatusCode: 0, StatusMsg: "upload successful"},
	})
	return
}
