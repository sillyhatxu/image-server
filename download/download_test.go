package download

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDownloadImage(t *testing.T) {
	var dJPG Download = Image{URL: "http://pic.pc6.com/up/2016-12/20161220195164315314.jpg", Outpath: "/Users/shikuanxu/go/src/image-server/download", Outfile: "test.jpg"}
	err := dJPG.DownloadImage()
	assert.Nil(t, err)

}

func TestDownloadImageNoSuffix(t *testing.T) {
	var noSuffix Download = Image{URL: "http://pin.aliyun.com/get_img?sessionid=d88ca3c24732ce2ee07e74ecc7e32b3b&identity=sm-tmallsearch&type=default", Outpath: "/Users/shikuanxu/go/src/image-server/download", Outfile: "2151324.jpg"}
	err := noSuffix.DownloadImage()
	assert.Nil(t, err)
}
