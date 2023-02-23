package main

import (
	"douyin/cache"
	"douyin/config"
	"douyin/dao"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {

	//pprof监听
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	var err error
	//1.初始化配置文件
	err = config.ConfInit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//2.初始化redis数据库
	err = cache.RedisPoolInit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//3.初始化mysql数据库
	err = dao.DbInit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//4.初始化http框架
	r := gin.Default()
	//5.初始化吧路由
	RouterInit(r)
	//6.启动程序
	err = r.Run(":8000")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
