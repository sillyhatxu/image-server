package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sillyhatxu/microlog"
	"image-server/config"
	"image-server/response"
	"net/http"
)

func InitialAPI() {
	log.Info("---------- initial api start ----------")
	router := gin.Default()
	stockRouterGroup := router.Group("/server-image")
	{
		stockRouterGroup.POST("/upload/file", uploadFile)
		stockRouterGroup.POST("/upload/url", uploadURL)
	}
	_ = router.Run(config.Conf.Http.Listen)
}

func uploadFile(context *gin.Context) {

	message := c.PostForm("message")
	nick := c.DefaultPostForm("nick", "anonymous")

	c.JSON(200, gin.H{
		"status":  "posted",
		"message": message,
		"nick":    nick,
	})

	context.JSON(http.StatusOK, response.Success(nil))
}

func uploadURL(context *gin.Context) {

	message := c.PostForm("message")
	nick := c.DefaultPostForm("nick", "anonymous")

	c.JSON(200, gin.H{
		"status":  "posted",
		"message": message,
		"nick":    nick,
	})

	context.JSON(http.StatusOK, response.Success(nil))
}
