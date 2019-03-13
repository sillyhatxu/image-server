package alicloudoss

import (
	"fmt"
	"github.com/sillyhatxu/mysql-client"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

const (
	DB              = "xxxx:xxxx@tcp(127.0.0.1:3306)/database"
	Endpoint        = "xxx.com"
	AccessKeyId     = "xxx"
	AccessKeySecret = "xxx"
	BucketName      = "sillyhat-blog"
)

func TestListBuckets(t *testing.T) {
	alicloud := &AliCloud{Endpoint: Endpoint, AccessKeyId: AccessKeyId, AccessKeySecret: AccessKeySecret}
	buckets, err := alicloud.ListBuckets()
	assert.Nil(t, err)
	for _, bucket := range buckets.Buckets {
		fmt.Println("Buckets:", bucket.Name)
	}
}

func TestWriteLog(t *testing.T) {
	dbclient.InitialDBClient(DB, 4, 4)
	writeLog("/Users/shikuanxu/go/src/image-server/test.png", "images/3308014985/906799682/518729469/25feabc0-40f2-407a-a1a2-d771eecf0f54.png")
}

func TestUploadImageFromURL(t *testing.T) {
	dbclient.InitialDBClient(DB, 4, 4)
	url := "http://pic.pc6.com/up/2016-12/20161220195164315314.jpg"
	tempPath := "/Users/shikuanxu/Downloads/images/temp"
	alicloud := &AliCloud{Endpoint: Endpoint, AccessKeyId: AccessKeyId, AccessKeySecret: AccessKeySecret}
	outputFile, err := alicloud.UploadImageFromURL(BucketName, url, tempPath)
	assert.Nil(t, err)
	fmt.Println(outputFile)
}

func TestUploadImageFromURLNoSuffix(t *testing.T) {
	dbclient.InitialDBClient(DB, 4, 4)
	tempPath := "/Users/shikuanxu/Downloads/images/temp"
	urls := []string{
		"https://img.fonwall.ru/o/62/yaponka-makiyaj-kimono-veer.jpg?route=mid&amp;h=750",
		"https://lh3.googleusercontent.com/ctm2WTdH0Ft5k8vNsKhOEA3dYmwmDLqG5oBftsQyXohuTmbzN0txKP_OqrKxh4PG1PvGM3jtIWfHDT9zJqEBctB4x6MwACH3pnTAZd4SOl8YakWA-5Vdf3ae2Jc3ERBcQ9h3gHGwwql5ZnfqVufeP8qwAnP1QbuPqnTe5NVdijD3UO6DkvMtPmjIJU9oExfpOKnctoWezcqjnIrouGHbhmgAk33p5Hhz2DV5oxaoylKp1QRLrIzaE20eons-k4tK9M7YKZEg1VTvv9ZMzZ_842G0W4XZVYnP9dc635KptlihxTg_rMmW05cSYDW0Gw8jjoPSfBR6J7idHkmUM7Ec-w7rXpp9OvDoiSzYXVXPqQtCzTxIg5Ig9rjb-K3vrJ5gKbXxWklbOsu_E7TinXK_2VucIDOkm8hg8r7E71USg9Yo29OOV0sXn_rDy8WBoEd5RTNrz3Sbgfk7hl9dHjS5irtZxp2NtPqRw6XlbcCRWLF73SSLAXcI7lWQF0UyoCEceNVygn0A2vkNVOKbXxD5auQ9apUdw_-f9i_mrkeW-pSo2NP3QUk0EQe6Vw6gZy1DsmUfxEJoiB_Eo_4bhdTOY1AT7f9seZs5kB5T6XZBw-HLXpGiQowgrzJIKCPXwCyNuO-w4G-0SWqVneEb8TlgodRq4Ung0ek=w1436-h897-no",
	}
	for _, url := range urls {
		alicloud := &AliCloud{Endpoint: Endpoint, AccessKeyId: AccessKeyId, AccessKeySecret: AccessKeySecret}
		outputFile, err := alicloud.UploadImageFromURL(BucketName, url, tempPath)
		assert.Nil(t, err)
		fmt.Println(outputFile)
	}
}

func TestUploadImage(t *testing.T) {
	//uploadFiles := []string{"/Users/shikuanxu/Downloads/images/flutter01/bash-profile.png",""}
	alicloud := &AliCloud{Endpoint: Endpoint, AccessKeyId: AccessKeyId, AccessKeySecret: AccessKeySecret}
	outputFile, err := alicloud.UploadImageFromPath(BucketName, "/Users/shikuanxu/Downloads/images/flutter01/bash-profile.png")
	assert.Nil(t, err)
	fmt.Println(outputFile)
}

func TestUploadImageFromFolder(t *testing.T) {
	folder := "/Users/shikuanxu/Downloads/images/flutter01/"
	alicloud := &AliCloud{Endpoint: Endpoint, AccessKeyId: AccessKeyId, AccessKeySecret: AccessKeySecret}
	files, err := ioutil.ReadDir(folder)
	assert.Nil(t, err)
	for _, f := range files {
		if f.Name() == ".DS_Store" {
			continue
		}
		outputFile, err := alicloud.UploadImageFromPath(BucketName, folder+f.Name())
		assert.Nil(t, err)
		fmt.Println(outputFile)
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
	err := alicloud.SetBucketReferer(BucketName, referers)
	assert.Nil(t, err)
}
