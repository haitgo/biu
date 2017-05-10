//http请求上下文
package biu

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/haitgo/biu/upload"
	"github.com/haitgo/biu/validate"
	"github.com/haitgo/biu/websocket"
	"github.com/haitgo/tool/any"
)

type Content struct {
	session     Sessioner
	Params      map[string]string
	Request     *http.Request
	Writer      *Writer
	nextHandle  func()
	abortHandle func()
}

//发起请求的ip地址
func (this *Content) ClientIP() string {
	return ""
}

//读取查询参数
func (this *Content) Query(key string) *any.Any {
	req := this.Request
	var dt interface{}
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		dt = values[0]
	}
	return any.New(dt)
}

//绑定查询参数
func (this *Content) BindQuery() *validate.Bind {
	data := map[string]string{}
	for k, v := range this.Request.URL.Query() {
		data[k] = strings.TrimSpace(v[0])
	}
	return validate.New(data)
}

//读取表单参数
func (this *Content) Form(key string) *any.Any {
	req := this.Request
	req.ParseMultipartForm(32 << 20) // 32 MB
	if values := req.PostForm[key]; len(values) > 0 {
		return any.New(values[0])
	}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			return any.New(values[0])
		}
	}
	return any.New(nil)
}

//绑定表单参数
func (this *Content) BindForm() *validate.Bind {
	req := this.Request
	req.ParseMultipartForm(32 << 20) // 32 MB
	data := map[string]string{}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		for k, v := range req.MultipartForm.Value {
			data[k] = strings.TrimSpace(v[0])
		}
	} else {
		for k, v := range req.PostForm {
			data[k] = strings.TrimSpace(v[0])
		}
	}
	return validate.New(data)
}

//请求头读取
func (this *Content) Header(name string) {

}

//输出请求头
func (this *Content) SetHeader(name, value string) {

}

//打开session对象
func (this *Content) Session() Sessioner {
	return this.session
}

//读取cookie
func (this *Content) Cookie(name string) {

}

//设置cookie
func (this *Content) SetCookie(name string, value interface{}) {

}

//url参数读取
func (this *Content) Param(name string) *any.Any {
	return any.New(this.Params[name])
}

//追加参数
func (this *Content) paramAppend(params map[string]string) {
	if this.Params == nil {
		this.Params = make(map[string]string)
	}
	for k, v := range params {
		this.Params[k] = v
	}
}

//上传文件接收（针对单文件）
func (this *Content) Upload(name string) *upload.Upload {
	return upload.NewUpload(this.Request, name)
}

//websocket对象
func (this *Content) WebSocket(opt *websocket.Option) (*websocket.Websocket, error) {
	return websocket.New(this.Writer, this.Request, opt)
}

//返回状态码
func (this *Content) Status(code int) *Content {
	this.Writer.WriteHeader(code)
	return this
}

//返回html页面
func (this *Content) Html(tplFile string, data interface{}) {

}

//返回字符串
func (this *Content) String(str string) error {
	_, e := this.Writer.Write([]byte(str))
	return e
}

//返回json
func (this *Content) Json(data interface{}) error {
	bt, _ := json.Marshal(data)
	_, e := this.Writer.Write(bt)
	return e
}

//返回文件
func (this *Content) File(file string) {

}

//执行下一个中间件
func (this *Content) Next() {
	this.nextHandle()
}

//结束访问
func (this *Content) Abort() {
	this.abortHandle()
}

//跳转
func (this *Content) Redirect(url string) {

}
