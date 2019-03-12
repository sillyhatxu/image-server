package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sillyhatxu/sillyhat-cloud-web/jwt"
	log "github.com/sirupsen/logrus"
	"image-server/api/dto"
	"image-server/config"
	"image-server/response"
	"image-server/service"
	"image-server/token"
	"net/http"
	"sillyhat-cloud-utils/cache"
	"time"
)

func InitialAPI() {
	log.Info("---------- initial api start ----------")
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Use(HandlerInterceptorAdapter())
	router.Use(gin.LoggerWithFormatter(Logger))
	router.Use(gin.Recovery())
	loginRouterGroup := router.Group("/server-image/login")
	{
		loginRouterGroup.POST("/", login)
	}
	//stockRouterGroup := router.Group("/server-image/upload").Use(AuthRequired())
	stockRouterGroup := router.Group("/server-image/upload")
	{
		stockRouterGroup.POST("/file", uploadFile)
		stockRouterGroup.POST("/url", uploadURL)
	}
	userRouterGroup := router.Group("/server-image/users").Use(AuthRequired())
	{
		userRouterGroup.GET("/", getUserList)
		userRouterGroup.GET("/{id}", getUserById)
	}

	_ = router.Run(config.Conf.Http.Listen)
}

func Logger(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

func AuthRequired() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenHeader := context.GetHeader("Authorization") //Grab the token from the header
		if tokenHeader == "" {                            //Token is missing, returns with error code 403 Unauthorized
			context.Header("Content-Type", "application/json")
			context.JSON(http.StatusForbidden, response.Error(nil, "Missing auth token"))
			context.Abort() //return
			return
		}
		userToken := token.UserToken{}
		_, err := jwt.ParseTokenString(tokenHeader, &userToken)
		if err != nil {
			context.Header("Content-Type", "application/json")
			context.JSON(http.StatusForbidden, response.Error(nil, "Missing auth token"))
			context.Abort() //return
			return
		}
		_, found := cache.Get(userToken.UserId)
		if !found {
			context.Header("Content-Type", "application/json")
			context.JSON(http.StatusGatewayTimeout, response.Error(nil, "Authorized time out"))
			context.Abort() //return
			return
		}
	}
}

func HandlerInterceptorAdapter() gin.HandlerFunc {
	return func(context *gin.Context) {
		t := time.Now()
		// Set example variable
		context.Set("example", "12345")
		// before request
		context.Next()
		// after request
		latency := time.Since(t)
		log.Print(latency)
		// access the status we are sending
		status := context.Writer.Status()
		log.Println(status)
	}
}

func login(context *gin.Context) {
	var requestBody dto.LoginDTO
	if err := context.ShouldBindJSON(&requestBody); err == nil {
		tokenSrc, err := service.Login(requestBody)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		context.Header("Content-Type", "application/json")
		context.Header("Authorization", tokenSrc)
		context.JSON(http.StatusOK, response.Success(requestBody))
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func getUserList(context *gin.Context) {
	context.JSON(http.StatusOK, response.Success(nil))
}

func getUserById(context *gin.Context) {
	context.JSON(http.StatusOK, response.Success(nil))
}

func uploadFile(context *gin.Context) {
	// Multipart form
	form, err := context.MultipartForm()
	if err != nil {
		context.JSON(http.StatusOK, response.Success(err.Error()))
		return
	}
	files := form.File["file[]"]
	for _, file := range files {
		log.Println(file.Filename)
	}
	context.JSON(http.StatusOK, response.Success(""))
}

func uploadURL(context *gin.Context) {

	context.JSON(http.StatusOK, response.Success(nil))
}
