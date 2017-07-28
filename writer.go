//writer装饰过滤器，并对输出内容进行缓存，模拟php的ob缓存功能
package biu

import (
	"net/http"
)

type Writer struct {
	statusCode int                 //状态码
	writer     http.ResponseWriter //
	obStart    bool                //开启ob缓存
	obCache    []byte              //ob缓存
}

//
func newWriter(w http.ResponseWriter) *Writer {
	return &Writer{writer: w}
}

//Writer接口实现，输出请求头
func (this *Writer) Header() http.Header {
	return this.writer.Header()
}

//Writer接口实现，输出内容
func (this *Writer) Write(bt []byte) (int, error) {
	if this.obCache == nil {
		this.obCache = make([]byte, 0)
	}
	if this.obStart {
		this.obCache = append(this.obCache, bt...)
		return len(bt), nil
	}
	return this.writer.Write(bt)
}

//Writer接口实现，输出状态码
func (this *Writer) WriteHeader(code int) {
	this.statusCode = code
	this.writer.WriteHeader(code)
}

//当前返回状态码
func (this *Writer) Status() int {
	return this.statusCode
}

//类似php的ob_start
func (this *Writer) OBstart() {
	this.obStart = true
}

//获取ob缓存数据
func (this *Writer) OBget() []byte {
	return this.obCache
}

//输出ob缓存内容
func (this *Writer) OBflush() {
	this.writer.Write(this.obCache)
}
