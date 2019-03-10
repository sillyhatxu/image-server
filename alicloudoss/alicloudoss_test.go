package alicloudoss

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

const (
	Endpoint        = "xxx.com"
	AccessKeyId     = "xxx"
	AccessKeySecret = "xxx"
)

func TestListBuckets(t *testing.T) {
	alicloud := &AliCloud{Endpoint: Endpoint, AccessKeyId: AccessKeyId, AccessKeySecret: AccessKeySecret}
	buckets, err := alicloud.ListBuckets()
	assert.Nil(t, err)
	for _, bucket := range buckets.Buckets {
		fmt.Println("Buckets:", bucket.Name)
	}
}

func TestUploadImage(t *testing.T) {
	//uploadFiles := []string{"/Users/shikuanxu/Downloads/images/flutter01/bash-profile.png",""}
	alicloud := &AliCloud{Endpoint: Endpoint, AccessKeyId: AccessKeyId, AccessKeySecret: AccessKeySecret}
	err := alicloud.UploadImage("sillyhat-blog", "/Users/shikuanxu/Downloads/images/flutter01/bash-profile.png")
	assert.Nil(t, err)
}

func TestUploadImageFromFolder(t *testing.T) {
	folder := "/Users/shikuanxu/Downloads/images/flutter01/"
	bucketName := "sillyhat-blog"
	alicloud := &AliCloud{Endpoint: Endpoint, AccessKeyId: AccessKeyId, AccessKeySecret: AccessKeySecret}
	files, err := ioutil.ReadDir(folder)
	assert.Nil(t, err)
	for _, f := range files {
		if f.Name() == ".DS_Store" {
			continue
		}
		err := alicloud.UploadImage(bucketName, folder+f.Name())
		assert.Nil(t, err)
	}
}

func TestCreateFile(t *testing.T) {
	filename := "/Users/shikuanxu/Downloads/images/flutter01/bash-profile.png"
	fmt.Println(filepath.Ext(filename))
	fmt.Println(strings.TrimRight(filename, filepath.Ext(filename)))
	fmt.Println(strings.TrimSuffix(filename, path.Ext(filename)))
	assert.EqualValues(t, createFile(filename), "/Users/shikuanxu/Downloads/images/flutter01/bash-profile")
}

func TestSetBucketReferer(t *testing.T) {
	referers := []string{"http://test.com", "http://*.test.com", "*.console.aliyun.com"}
	alicloud := &AliCloud{Endpoint: Endpoint, AccessKeyId: AccessKeyId, AccessKeySecret: AccessKeySecret}
	err := alicloud.SetBucketReferer("sillyhat-blog", referers)
	assert.Nil(t, err)
}
