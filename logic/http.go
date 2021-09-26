package logic

import "github.com/gin-gonic/gin"

func RegisterHttp(router *gin.Engine) {
	router.POST("/calculate/", Calculate)
}