package upload

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"strings"
)

//图片裁剪加水印
type ImageCut struct {
	rawImage    image.Image //原始img资源
	handleImage image.Image //正在处理的img资源
	err         error       //排除当前错误
}

//通过文件创建图片裁剪对象
func NewImageCutReadInFile(filePath string) (o *ImageCut) {
	o = new(ImageCut)
	o.rawImage, _ = o.readingImageInFile(filePath)
	return
}

//使用byte实例化imgCut对象
func NewImageCutReadInByte(bt []byte, ext string) (o *ImageCut) {
	o = new(ImageCut)
	rd := bytes.NewReader(bt)
	o.rawImage, _ = o.readingImage(ext, rd)
	return
}

//压缩至指定大小,width和height其中1个可以为0，另一个不能为0
//如果为0则是自适应,都为0则是原始大小
func (this *ImageCut) Resize(x, y int) *ImageCut {
	if this.err != nil {
		return this
	}
	p := this.rawImage.Bounds() //原图大小
	var bl float64
	bl = float64(p.Max.X) / float64(p.Max.Y) //缩放比例
	if x == 0 && y > 0 {                     //如果宽为0，则以高等比缩放
		x = int(float64(y) * bl)
	} else if y == 0 && x > 0 { //如果高为0，则以宽等比缩放
		y = int(float64(x) / bl)
	} else {
		x = p.Max.X
		y = p.Max.Y
	}
	this.handleImage = this.resizeHandle(this.rawImage, this.rawImage.Bounds(), x, y)
	return this
}

//添加水印，为了美观，只支持图片水印
//file 水印图片路径，
//position 水印坐标 1 左上角，2右上角，3左下角，4右下角
func (this *ImageCut) WaterMark(filePath string, position int) *ImageCut {
	if this.err != nil {
		return this
	}
	img, err := this.readingImageInFile(filePath)
	if err != nil {
		this.err = err //水印图片图片解码错误
		return this
	}
	if this.handleImage == nil {
		this.handleImage = this.rawImage
	}
	var rect image.Rectangle
	aBounds := img.Bounds()
	bBounds := this.handleImage.Bounds()
	ax, ay := aBounds.Max.X, aBounds.Max.Y
	bx, by := bBounds.Max.X, bBounds.Max.Y
	switch position {
	case 1: //上
		rect = image.Rect(0, 0, ax, ay)
	case 2: //右
		rect = image.Rect(bx-ax, 0, bx, ay)
	case 3: //下
		rect = image.Rect(bx-ax, by-ay, bx, by)
	case 4: //左
		rect = image.Rect(0, by-ay, ax, by)
	}
	m := image.NewNRGBA(bBounds)
	draw.Draw(m, bBounds, this.handleImage, image.ZP, draw.Src)
	draw.Draw(m, rect, img, image.Pt(0, 0), draw.Src)
	this.handleImage = m
	return this
}

//将图片写入到文件
func (this *ImageCut) WriteFile(filePath string) error {
	if this.err != nil {
		return this.err
	}
	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, 744); err != nil {
		return err
	}
	f, err := os.Create(filePath)
	defer func() {
		f.Close()
		this.handleImage = nil
	}()
	if err != nil {
		return err
	}
	ext := strings.ToLower(path.Ext(filePath)) //获取文件后缀
	if ext == ".jpg" {
		err = jpeg.Encode(f, this.handleImage, &jpeg.Options{90})
	} else if ext == ".png" {
		err = png.Encode(f, this.handleImage)
	} else {
		return errors.New("图片生成错误")
	}
	return err
}

//从文件读取图片资源
func (this *ImageCut) readingImageInFile(filePath string) (img image.Image, err error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, errors.New("图片文件读取失败")
	}
	return this.readingImage(path.Ext(filePath), f)
}

//读取图片资源
func (this *ImageCut) readingImage(ext string, r io.Reader) (img image.Image, err error) {
	ext = strings.ToLower(ext) //获取文件后缀
	if ext == ".jpg" {
		img, err = jpeg.Decode(r)
	} else if ext == ".png" {
		img, err = png.Decode(r)
	}
	this.err = err
	return
}

//缩放处理
func (this *ImageCut) resizeHandle(m image.Image, r image.Rectangle, w, h int) image.Image {
	if w < 0 || h < 0 {
		return nil
	}
	if w == 0 || h == 0 || r.Dx() <= 0 || r.Dy() <= 0 {
		return image.NewRGBA64(image.Rect(0, 0, w, h))
	}
	curw, curh := r.Dx(), r.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Get a source pixel.
			subx := x * curw / w
			suby := y * curh / h
			r32, g32, b32, a32 := m.At(subx, suby).RGBA()
			r := uint8(r32 >> 8)
			g := uint8(g32 >> 8)
			b := uint8(b32 >> 8)
			a := uint8(a32 >> 8)
			img.SetRGBA(x, y, color.RGBA{r, g, b, a})
		}
	}
	return img
}
