package app

import (
	"github.com/gin-gonic/gin"
	"log"
)

var router *gin.Engine

func init() {
	router = gin.Default()
}

func StartApp() {
	mapUrls()

	log.Fatal(router.Run(""))
}
