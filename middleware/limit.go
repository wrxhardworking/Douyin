package middleware

import (
	"douyin/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

func Check() gin.HandlerFunc {

	return func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("username")
		//传给handle层
		//5～16字节，允许字母、数字、下划线，以字母开头
		matchString := "^[a-zA-Z][a-zA-Z0-9_]{4,15}$"
		usernameMatch, _ := regexp.MatchString(matchString, username)
		passwordMatch, _ := regexp.MatchString(matchString, password)

		if usernameMatch && passwordMatch != true {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1, StatusMsg: "Account or password is illegal",
			})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Set("password", password)
		//挂起
		c.Next()
	}

}
