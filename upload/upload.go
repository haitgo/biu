//为什么要多个文件同时上传，传一个挺好，可以一个一个传啊
package upload

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type Upload struct {
	r        *http.Request //http请求
	allowExt []string      //允许上传的文件类型
	maxSize  int           //允许上传的文件大小，单位kb
	ext      string        //文件后缀
	formName string        //上传表单名
}

//创建上传
func NewUpload(r *http.Request, formName string) *Upload {
	o := new(Upload)
	o.r = r
	o.formName = formName
	return o
}

//设置允许上传的文件类型后缀名
func (this *Upload) AllowExt(ext ...string) *Upload {
	this.allowExt = ext
	return this
}

//允许上传最大文件大小
func (this *Upload) MaxSize(size int) *Upload {
	this.maxSize = size
	return this
}

//过滤文件后缀名
func (this *Upload) filterExt(ext string) bool {
	if len(this.allowExt) == 0 {
		return true
	}
	for _, e := range this.allowExt {
		e = strings.ToLower(e)
		if e == ext {
			return true //允许上传的类型
		}
	}
	return false
}

//保存文件，入股目录不存在则自动创建
func (this *Upload) WriteFile(filePath string) error {
	bt, err := this.handle(this.formName)
	if err != nil {
		return err
	}
	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, 744); err != nil {
		return err
	}
	if len(bt) == 0 {
		return errors.New("上传资源无法保存")
	}
	return ioutil.WriteFile(filePath, bt, 755)
}

//上传图片
func (this *Upload) GetImage() (*ImageCut, error) {
	bt, err := this.handle(this.formName)
	if err != nil {
		return nil, err
	}
	imgCut := NewImageCutReadInByte(bt, this.ext)
	return imgCut, err
}

//上传处理
func (this *Upload) handle(name string) ([]byte, error) {
	file, head, err := this.r.FormFile(name)
	if err != nil {
		return nil, err //上传错误
	}
	defer file.Close()
	this.ext = strings.ToLower(path.Ext(head.Filename))
	if !this.filterExt(this.ext) {
		return nil, errors.New("不允许上传的文件类型")
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err //无法读取上传文件
	}
	if this.maxSize > 0 && len(bytes) > this.maxSize {
		return nil, errors.New("上传文件过大。")
	}
	return bytes, nil

}
