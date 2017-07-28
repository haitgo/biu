//http请求上下文
package biu

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/haitgo/biu/bind"
	"github.com/haitgo/biu/upload"
)

type Context struct {
	session     SessionBase
	Params      map[string]string
	Request     *http.Request
	Writer      *Writer
	nextHandle  func()
	abortHandle func()
}

//发起请求的ip地址
func (this *Context) ClientIP() string {
	return ""
}

//读取查询参数
func (this *Context) Query(key string) string {
	req := this.Request
	var dt string
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		dt = values[0]
	}
	return dt
}

//绑定查询参数
func (this *Context) BindQuery() *bind.Bind {
	data := map[string]string{}
	for k, v := range this.Request.URL.Query() {
		data[k] = strings.TrimSpace(v[0])
	}
	return bind.New(data)
}

//读取表单参数
func (this *Context) Form(key string) string {
	req := this.Request
	req.ParseMultipartForm(32 << 20) // 32 MB
	var ret string
	if values := req.PostForm[key]; len(values) > 0 {
		ret = values[0]
	} else if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			ret = values[0]
		}
	}
	return ret
}

//绑定表单参数
func (this *Context) BindForm() *bind.Bind {
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
	return bind.New(data)
}

//请求头读取
func (this *Context) Header(key string) string {
	return this.Request.Header.Get(key)
}

//输出请求头
func (this *Context) SetHeader(key, value string) {
	this.Writer.Header().Set(key, value)
}

//打开session对象
func (this *Context) Session() SessionBase {
	return this.session
}

//读取cookie
func (this *Context) Cookie(name string) string {
	var ret string
	if ck, e := this.Request.Cookie(name); e != nil {
		ret = ck.Value
	}
	return ret
}

//设置cookie
func (this *Context) SetCookie(ck *http.Cookie) {
	this.Request.AddCookie(ck)
}

//url参数读取
func (this *Context) Param(name string) string {
	return this.Params[name]
}

//追加参数
func (this *Context) paramAppend(params map[string]string) {
	if this.Params == nil {
		this.Params = make(map[string]string)
	}
	for k, v := range params {
		this.Params[k] = v
	}
}

//上传文件接收（针对单文件）
func (this *Context) Upload(name string) *upload.Upload {
	return upload.NewUpload(this.Request, name)
}

//返回状态码
func (this *Context) Status(code int) *Context {
	this.Writer.WriteHeader(code)
	return this
}

//返回html页面
func (this *Context) Html(tplFile string, data interface{}) error {
	this.SetHeader("content-type", "text/html")
	tpl, e1 := template.ParseFiles(tplFile)
	if e1 != nil {
		return e1
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	e2 := tpl.Execute(buf, data)
	if e2 != nil {
		return e2
	}
	_, e3 := this.Writer.Write(buf.Bytes())
	return e3
}

//返回字符串
func (this *Context) String(str string) error {
	this.SetHeader("content-type", "text/plain")
	_, e := this.Writer.Write([]byte(str))
	return e
}

//返回json
func (this *Context) Json(data interface{}) error {
	this.SetHeader("content-type", "text/json")
	wt := json.NewEncoder(this.Writer)
	return wt.Encode(data)
}

//常用接口json数据返回
//例如：c.Api(1,biu.H{});
//结构：{code:1,data:"sss"}
func (this *Context) Api(code int, data interface{}) error {
	ret := make(map[string]interface{})
	ret["code"] = code
	ret["data"] = data
	return this.Json(ret)
}

//返回文件
func (this *Context) File(file string) {

}

//执行下一个中间件
func (this *Context) Next() {
	this.nextHandle()
}

//结束访问
func (this *Context) Abort() {
	this.abortHandle()
}

//跳转
func (this *Context) Redirect(url string) {
	this.Writer.Header().Set("location", url)
}
