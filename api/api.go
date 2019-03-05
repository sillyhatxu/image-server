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
		stockRouterGroup.POST("/upload", upload)
	}
	_ = router.Run(config.Conf.Http.Listen)
}

func upload(context *gin.Context) {
	context.JSON(http.StatusOK, response.Success(nil))
}
