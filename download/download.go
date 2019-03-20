package download

import (
	"github.com/sillyhatxu/sillyhat-cloud-utils/uuid"
	"image-server/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Download interface {
	DownloadImage() error
}

type Image struct {
	URL     string
	Outpath string
	Outfile string
}

func getFileName() string {
	filename := uuid.UUID()
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	return name
}

func (i Image) DownloadImage() error {
	err := utils.CheckFolder(i.Outpath)
	if err != nil {
		return err
	}

	response, err := http.Get(i.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	file, err := os.Create(i.Outpath + "/" + i.Outfile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
