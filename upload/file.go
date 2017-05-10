package upload

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type fileData struct {
	data []byte
}

//保存文件，入股目录不存在则自动创建
func (this *fileData) WriteFile(filePath string) error {
	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, 744); err != nil {
		return err
	}
	if len(this.data) == 0 {
		return errors.New("上传资源无法保存")
	}
	return ioutil.WriteFile(filePath, this.data, 755)
}
