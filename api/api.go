package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sillyhatxu/sillyhat-cloud-utils/cache"
	"github.com/sillyhatxu/sillyhat-cloud-utils/uuid"
	"github.com/sillyhatxu/sillyhat-cloud-web/jwt"
	log "github.com/sirupsen/logrus"
	"image-server/alicloudoss"
	"image-server/api/dto"
	"image-server/config"
	"image-server/response"
	"image-server/service"
	"image-server/token"
	"net/http"
	"time"
)

func InitialAPI() {
	log.Info("---------- initial api start ----------")
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Use(HandlerInterceptorAdapter())
	router.Use(gin.LoggerWithFormatter(Logger))
	router.Use(gin.Recovery())
	loginRouterGroup := router.Group("/image-server/login")
	{
		loginRouterGroup.POST("/", login)
	}
	//stockRouterGroup := router.Group("/server-image/upload").Use(AuthRequired())
	stockRouterGroup := router.Group("/image-server")
	{
		stockRouterGroup.POST("/upload-file", uploadFile)
		stockRouterGroup.POST("/upload-multi-file", uploadFileMultipart)
		stockRouterGroup.POST("/upload-url", uploadURL)
		stockRouterGroup.POST("/upload-post", uploadImageByUrl)
	}
	userRouterGroup := router.Group("/image-server/users").Use(AuthRequired())
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

/*
curl -X POST http://localhost:8080/server-image/upload-file \
  -F "file=@/Users/shikuanxu/Downloads/images/course-outline.png" \
  -H "Content-Type: multipart/form-data"
*/
func uploadFile(context *gin.Context) {
	fileHeader, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusOK, response.Success(err.Error()))
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		context.JSON(http.StatusOK, response.Success(err.Error()))
		return
	}
	alicloud := &alicloudoss.AliCloud{Endpoint: config.Conf.AliCloud.Endpoint, AccessKeyId: config.Conf.AliCloud.AccessKeyId, AccessKeySecret: config.Conf.AliCloud.AccessKeySecret}
	uploadPath, err := alicloud.UploadImageFromFile(config.Conf.AliCloud.ImageBlogBucketName, fileHeader.Filename, file)
	// Upload the file to specific dst.
	// c.SaveUploadedFile(file, dst)
	context.JSON(http.StatusOK, response.Success(uploadPath))
}

/*
curl -X POST http://localhost:8080/server-image/upload-multi-file \
  -F "file=@/Users/shikuanxu/Downloads/images/course-outline.png" \
  -F "file=@/Users/shikuanxu/Downloads/images/create-scrapy-project.png" \
  -F "file=@/Users/shikuanxu/Downloads/images/widget-tree-for-this-ui.png" \
  -H "Content-Type: multipart/form-data"
*/
func uploadFileMultipart(context *gin.Context) {
	form, err := context.MultipartForm()
	if err != nil {
		context.JSON(http.StatusOK, response.Success(err.Error()))
		return
	}
	fileHeaders := form.File["file"]
	result := make([]dto.MultiFileDTO, len(fileHeaders))
	for index, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			context.JSON(http.StatusOK, response.Success(err.Error()))
			return
		}
		alicloud := &alicloudoss.AliCloud{Endpoint: config.Conf.AliCloud.Endpoint, AccessKeyId: config.Conf.AliCloud.AccessKeyId, AccessKeySecret: config.Conf.AliCloud.AccessKeySecret}
		outputFile, err := alicloud.UploadImageFromFile(config.Conf.AliCloud.ImageBlogBucketName, fileHeader.Filename, file)
		result[index] = *&dto.MultiFileDTO{UploadFile: fileHeader.Filename, OutputFile: outputFile}
	}
	context.JSON(http.StatusOK, response.Success(result))
}

func uploadURL(context *gin.Context) {
	var requestBody dto.UploadURLDTO
	err := context.ShouldBindJSON(&requestBody)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	outputFile, err := service.UploadFileByURL(requestBody.URL)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, response.Success(outputFile))
}

func uploadImageByUrl(context *gin.Context) {
	var requestBody dto.UploadURLDTO
	err := context.ShouldBindJSON(&requestBody)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if requestBody.URL == "www.error.com" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "testerror"})
		return
	}
	context.JSON(http.StatusOK, response.Success(uuid.UUID()))
}
