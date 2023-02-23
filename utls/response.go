package utls

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response[T any] struct {
	StatusCode int    `json:"status_code,omitempty"`
	StatusMsg  string `json:"status_msg,omitempty"`
	Message    T
}

type ResponseHelper[T any] struct {
	c           *gin.Context
	GoodMessage string
	BadMessage  string
	Message     T
}

func (helper ResponseHelper[T]) BadResponse() {
	helper.c.JSON(http.StatusOK, Response[T]{
		StatusCode: 1,
		StatusMsg:  helper.BadMessage,
		Message:    helper.Message,
	})
}

func (helper ResponseHelper[T]) GoodResponse() {
	helper.c.JSON(http.StatusOK, Response[T]{
		StatusCode: 0,
		StatusMsg:  helper.GoodMessage,
		Message:    helper.Message,
	})
}
