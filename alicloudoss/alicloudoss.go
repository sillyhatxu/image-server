package alicloudoss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sillyhatxu/sillyhat-cloud-utils/encryption/hash"
	"github.com/sillyhatxu/sillyhat-cloud-utils/uuid"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strconv"
	"time"
)

const default_image_folder = "images"
const default_image_error_folder = "images/563185489"

type AliCloud struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
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

func (ali AliCloud) UploadImage(bucketName, uploadFile string) error {
	client, err := getClient(ali.Endpoint, ali.AccessKeyId, ali.AccessKeySecret)
	if err != nil {
		log.Error("Get oss client error.", err)
		return err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Errorf("Get bucket [%v] error.%v", bucketName, err)
		return err
	}
	err = bucket.PutObjectFromFile(createFile(uploadFile), uploadFile)
	if err != nil {
		log.Errorf("Upload image [%v] to bucket [%v] error.%v", uploadFile, bucketName, err)
		return err
	}
	return nil
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
