package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sillyhatxu/sillyhat-cloud-utils/cache"
	"github.com/sillyhatxu/sillyhat-cloud-web/jwt"
	log "github.com/sirupsen/logrus"
	"image-server/alicloudoss"
	"image-server/api/dto"
	"image-server/config"
	"image-server/response"
	"image-server/service"
	"image-server/token"
	"math/rand"
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
curl -X POST http://localhost:8080/image-server/upload-file \
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
	images := []string{"images/3308014985/25/72a2f9ee-87f5-45c1-bda8-0e83a3b16c94.png",
		"images/3308014985/15/ce2b3298-be2e-47db-92b6-b4c962437d42.png",
		"images/3308014985/50/b3f565a6-f1c6-4b1b-b65d-e557c1304e02.png",
		"images/3308014985/41/eb70ed48-02a9-4945-8fcc-3c22fd3afc6e.png",
		"images/3308014985/26/e61888ed-2736-4c7a-beb2-cdf495d01aac.png",
		"images/3308014985/14/fe578ae1-1d81-45f9-8193-9217ab8942a4.png",
		"images/3308014985/39/c05967e6-933b-4402-bad4-886250ac8d8b.png",
		"images/3308014985/12/980ff08a-d64a-457f-b815-a72d8ab985c1.png",
		"images/3308014985/41/43098daf-15e2-479c-8e13-10ac020b1204.png",
		"images/3308014985/18/e19fd329-5b0e-4f15-9c79-6bbd22b1c55a.png",
		"images/3308014985/24/55089a57-760f-4ac0-90f8-5da7fa559671.png",
		"images/3308014985/28/87d06c25-b230-4765-87a3-f1299626032b.png",
		"images/3308014985/15/ce7b521b-8b92-4270-bb87-affcca2a76e0.png",
		"images/3308014985/19/ee5350ea-635f-4ef6-b5ac-3a391686ccd2.png",
		"images/3308014985/67/abab8ec7-57bc-4106-b4fb-5c50423dfbad.png",
		"images/3308014985/14/a8626ecc-2b84-4bb3-b286-44a30ad1c318.png",
		"images/3308014985/36/cdf85ba4-9fb7-4297-a517-2af7f83713a9.png",
		"images/3308014985/42/bcae0608-a719-4277-bf83-69a84068e77d.png",
		"images/3308014985/86/fa7528aa-4a18-4e4e-b8df-44533ea713a3.png",
		"images/3308014985/30/b771de15-81b4-4f6a-a9dd-0cb4422cc5fe.png",
		"images/3308014985/40/7a74d0ec-32ab-4c92-bbc8-36a4b3901b71.png",
		"images/3308014985/42/c9749676-a328-4d12-9e61-3805722749c6.png",
		"images/3308014985/23/8d0c68ce-2ab5-4a33-9b2f-e24e00d84615.png",
		"images/3308014985/32/ca117bb1-3f24-48bf-bd9f-70cd3681102f.png",
		"images/3308014985/33/361e0d34-6aea-4268-a6c4-4f501bc605b5.png",
		"images/3308014985/34/2836ba84-95a2-42e4-9855-74fca68754ad.png",
		"images/3308014985/25/2a71935f-8296-4569-8d15-f0b6f57350ac.png",
		"images/3308014985/42/7164d562-3d54-4f36-b353-ceda29fdbe81.png",
		"images/3308014985/10/a1f0d5e3-8452-4a84-8429-c89a7c8993d7.png",
		"images/3308014985/14/d0a7f7e0-d44a-49ae-b5b6-34b9e1e1b7f5.png",
		"images/3308014985/39/7ddea680-ca39-456a-bacb-030229b950f3.png",
		"images/3308014985/33/5ed8c140-1642-4c4d-9cfe-7c4491125202.png",
		"images/3308014985/27/101cb5ad-c73f-413f-b98b-b01a6ce9ea02.png",
		"images/3308014985/33/a2fe050d-1166-4e7b-bcbb-38e96f3ca424.png",
		"images/3308014985/26/22a65dbc-e704-45ea-ae83-378d20a85d45.png",
		"images/3308014985/60/162e8be4-8d55-4a0c-a919-38117368df44.png",
		"images/3308014985/27/72bdcecd-3746-462c-adbd-39b106c00abb.png",
		"images/3308014985/33/e0a10c5b-7a59-4e40-967d-6ed74be2e55d.png",
		"images/3308014985/11/2a727081-401c-4e28-85d4-8ff6f6258e5f.png",
		"images/3308014985/16/6eb73cd8-f347-4aca-98fe-08a84ad3ad88.png",
		"images/3308014985/41/006b488c-59d5-499a-9c47-4c855988b032.png",
		"images/3308014985/42/54cb75b9-f894-4352-aeb4-989bde6073e3.png",
		"images/3308014985/36/67df622e-ca26-422b-970d-091c9fc8ceb4.png",
		"images/3308014985/35/409848e2-4ab2-4694-8ef5-9f58fdcdccf2.png",
		"images/3308014985/24/d440bed3-7276-48c7-8e00-2a8b3ce3a15a.png",
		"images/3308014985/42/985a12cf-49d1-46cb-acff-ea9ae6b39f0d.png",
		"images/3308014985/25/a8f8579f-3a69-426a-a9ac-10b4c1e302d3.png",
		"images/3308014985/41/e0a56cca-ba63-4e41-9d4f-584b2c527d50.png",
		"images/3308014985/11/63f773ba-42f5-4278-9cef-88a31a00086e.png",
		"images/3308014985/39/64165da6-725e-4743-a042-e1deefe582f0.png",
		"images/3308014985/14/91d2ec8a-9939-44e8-888e-c300e72b7795.png",
		"images/3308014985/18/61bad9ee-9a7c-42b5-ae9b-70f8b2e36780.png",
		"images/3308014985/32/0709a1ea-3d46-4739-bdd1-92d7dcf4e2ee.png",
		"images/3308014985/14/8949d812-a479-4321-bb2c-17d04ca4bbf5.png",
		"images/3308014985/26/03c4455d-131d-4a59-af7d-3b6120666b85.png",
		"images/3308014985/28/3402b2cb-a4ec-4cf1-9a0b-f6bbd1477cd3.png",
		"images/3308014985/37/c05720bd-2f8d-409a-8cc6-1477a2af2d54.png",
		"images/3308014985/18/4e109846-560e-4716-9a42-aa72c03f6d6d.png",
		"images/3308014985/16/32dc7f26-1169-4eba-9a94-b7e201f959d7.png",
		"images/3308014985/32/c3d74d6e-4d80-4fc9-9db4-bc50d0179d28.png",
		"images/3308014985/11/cd0cdc80-e523-48e1-a8e8-84f9e1e18e2b.png",
		"images/3308014985/16/6a5d7e86-c0a0-4ce7-b41b-ebb40068cac9.png",
		"images/3308014985/25/48981564-0c05-416f-ae66-29092eaa38c6.png",
		"images/3308014985/41/38dd02bf-16b1-4c42-87ab-48f446842dba.png",
		"images/3308014985/40/eddd816d-f757-4601-87f6-3089149a1b1b.png",
		"images/3308014985/28/b18a5b8f-a567-420a-99f0-0f5000b4773e.png",
		"images/3308014985/55/6d96e797-282c-47e2-9ea2-de0689f27462.png",
		"images/3308014985/57/c953555b-9e4d-4af8-b97c-b34dc93bcfab.png",
		"images/3308014985/22/15583c2f-200a-45b2-a31a-e3ccc3012397.png",
		"images/3308014985/23/230add8a-4cc0-4a71-83c3-b02c68bd428c.png",
		"images/3308014985/17/b382c3f7-c646-4d94-ba93-0e39a14cb333.png",
		"images/3308014985/14/15742b41-fe9d-4210-b0e2-066ba4d04727.png",
		"images/3308014985/13/73092245-2f6c-4f1e-95f7-ef823e10a646.png",
		"images/3308014985/22/d84f106b-8cc4-4c47-9a4e-be874f6b7175.png",
		"images/3308014985/17/3b5b0725-643f-414f-a608-dd06c94b445b.png",
		"images/3308014985/21/c6b6d13c-8722-40b3-98b0-707b78065fa9.png",
		"images/3308014985/13/6580033f-8f5e-493e-b24e-24bc3254bfaa.png",
		"images/3308014985/22/c27be0c3-14fd-4657-8227-30e8993b5fb7.png",
		"images/3308014985/15/c79554fe-af74-4caa-aad9-a128d8f7f579.png",
		"images/3308014985/21/c5c12b6b-e2be-4df1-87ed-78198e641df7.png",
		"images/3308014985/38/0de61a72-5dbe-45a0-b28d-e60a1959138f.png",
		"images/3308014985/11/ad13713f-417b-48c1-abce-86a4ca2ead67.png",
		"images/3308014985/58/0b1cc17d-b8dc-461e-9473-85f0d0a56966.png",
		"images/3308014985/30/096d6dbc-0a50-4099-bafb-73a188a7540c.png"}
	context.JSON(http.StatusOK, response.Success(images[rand.Intn(len(images))]))
}
