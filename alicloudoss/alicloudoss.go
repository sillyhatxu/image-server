package alicloudoss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sillyhatxu/sillyhat-cloud-utils/encryption/hash"
	"github.com/sillyhatxu/sillyhat-cloud-utils/uuid"
	log "github.com/sirupsen/logrus"
	"image-server/constants"
	"image-server/db"
	"image-server/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const default_image_folder = "images"
const default_image_error_folder = "images/563185489"

type AliCloud struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
}

func writeLog(uploadFile, outputFile string) {
	logjson := `{"upload_file": "` + uploadFile + `","output_file": "` + outputFile + `","upload_time": "` + time.Now().Format("2006-01-02 15:04:05") + `"}`
	err := db.Log(constants.LOG_TYP_UPLOAD_IMAGE, logjson)
	if err != nil {
		log.Errorf("insert log error.%v", err)
	}
}

func getClient(Endpoint string, AccessKeyId string, AccessKeySecret string) (*oss.Client, error) {
	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func createFile(filename string) string {
	year, month, day := time.Now().Date()
	suffix := filepath.Ext(filename)
	yearFolder, err := hash.HashValue32(strconv.Itoa(year))
	if err != nil {
		return default_image_error_folder
	}
	monthFolder, err := hash.HashValue32(strconv.Itoa(int(month)))
	if err != nil {
		return default_image_error_folder
	}
	dayFolder, err := hash.HashValue32(strconv.Itoa(day))
	if err != nil {
		return default_image_error_folder
	}
	return default_image_folder + "/" + yearFolder + "/" + monthFolder + "/" + dayFolder + "/" + uuid.UUID() + suffix
}

func (ali AliCloud) ListBuckets() (*oss.ListBucketsResult, error) {
	client, err := getClient(ali.Endpoint, ali.AccessKeyId, ali.AccessKeySecret)
	if err != nil {
		log.Error("Get oss client error.", err)
		return nil, err
	}
	lsRes, err := client.ListBuckets()
	if err != nil {
		log.Error("Get bucket list error.", err)
		return nil, err
	}
	return &lsRes, nil
}

func (ali AliCloud) UploadImageFromPath(bucketName, uploadFile string) (string, error) {
	client, err := getClient(ali.Endpoint, ali.AccessKeyId, ali.AccessKeySecret)
	if err != nil {
		log.Error("Get oss client error.", err)
		return "", err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Errorf("Get bucket [%v] error.%v", bucketName, err)
		return "", err
	}
	outputFile := createFile(uploadFile)
	err = bucket.PutObjectFromFile(outputFile, uploadFile)
	if err != nil {
		log.Errorf("Upload image [%v] to bucket [%v] error.%v", uploadFile, bucketName, err)
		return "", err
	}
	writeLog(uploadFile, outputFile)
	return outputFile, nil
}

func (ali AliCloud) UploadImageFromFile(bucketName, uploadFileName string, uploadFile io.Reader) (string, error) {
	client, err := getClient(ali.Endpoint, ali.AccessKeyId, ali.AccessKeySecret)
	if err != nil {
		log.Error("Get oss client error.", err)
		return "", err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Errorf("Get bucket [%v] error.%v", bucketName, err)
		return "", err
	}
	outputFile := createFile(uploadFileName)
	err = bucket.PutObject(outputFile, uploadFile)
	if err != nil {
		log.Errorf("Upload image [%v] to bucket [%v] error.%v", uploadFile, bucketName, err)
		return "", err
	}
	writeLog(uploadFileName, outputFile)
	return outputFile, nil
}

func getSuffix(url string) string {
	suffix := filepath.Ext(url)
	if strings.ToUpper(suffix) == ".JPG" || strings.ToUpper(suffix) == ".JPEG" || strings.ToUpper(suffix) == ".PNG" || strings.ToUpper(suffix) == ".GIF" {
		return suffix
	}
	return ".jpeg"
}

func tempDownloadImage(url, tempPath string) (string, error) {
	err := utils.CheckFolder(tempPath)
	if err != nil {
		return "", err
	}
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	outputFile := tempPath + "/" + uuid.UUID() + getSuffix(url)
	file, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	return outputFile, nil
}

func (ali AliCloud) UploadImageFromURL(bucketName, url, tempPath string) (string, error) {
	filePath, err := tempDownloadImage(url, tempPath)
	if err != nil {
		log.Errorf("Get bucket [%v] error.%v", bucketName, err)
		return "", err
	}
	defer os.Remove(filePath)
	return ali.UploadImageFromPath(bucketName, filePath)
}

func (ali AliCloud) SetBucketReferer(bucketName string, referers []string) error {
	client, err := getClient(ali.Endpoint, ali.AccessKeyId, ali.AccessKeySecret)
	if err != nil {
		log.Error("Get oss client error.", err)
		return err
	}
	err = client.SetBucketReferer(bucketName, referers, false)
	if err != nil {
		log.Errorf("Set bucket referer [%v] to bucket [%v] error.%v", referers, bucketName, err)
		return err
	}
	return nil
}
