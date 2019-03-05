package download

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddSKUs(t *testing.T) {
	var dJPG Download = Image{URL: "http://pic.pc6.com/up/2016-12/20161220195164315314.jpg", Outpath: "/Users/cookie/go/gopath/src/image-server/download/01/10", Outfile: "test.jpg"}
	err := dJPG.DownloadImage()
	assert.Nil(t, err)
	var noSuffix Download = Image{URL: "http://pin.aliyun.com/get_img?sessionid=d88ca3c24732ce2ee07e74ecc7e32b3b&identity=sm-tmallsearch&type=default", Outpath: "/Users/cookie/go/gopath/src/image-server/download/01/05", Outfile: "2151324.jpg"}
	err1 := noSuffix.DownloadImage()
	assert.Nil(t, err1)
}
