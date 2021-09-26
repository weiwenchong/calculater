package main

import (
	"github.com/gin-gonic/gin"
	"github.com/weiwenchong/calculator/logic"
	"log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	logic.RegisterHttp(router)

	err := router.Run("127.0.0.1:7777")
	if err != nil {
		log.Panicf("router.Run err:%v", err)
	}
}