package biu

import (
	"net/http"
	"path"
	"strings"
	"time"
)

type Biu struct {
	session    Sessioner //session处理接口
	staticPath string    //静态目录
	route      *Route    //路由对象
}

func New() *Biu {
	return new(Biu)
}

//路由
func (this *Biu) Route(call func(*Route)) {
	if this.route == nil {
		this.route = newRoute()
	}
	call(this.route)
}

//[设置]静态目录
func (this *Biu) StaticPath(path string) *Biu {
	this.staticPath = path
	return this
}

//[设置]注册session处理方法（默认内存session)
func (this *Biu) SessionHandle(session Sessioner) {
	this.session = session
}

// ServeHTTP
func (this *Biu) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoratorWriter := newWriter(w)
	c := &Content{
		Writer:  decoratorWriter,
		Request: r,
	}
	if this.session != nil {
		this.session.Server(w, r)
		c.session = this.session
	}
	this.routeHandle(c, strings.ToUpper(r.Method), r.URL.Path)
	println("code is", decoratorWriter.Status())
}

//静态文件处理
func (this *Biu) staticHandle(c *Content) {
	file := path.Base(c.Request.URL.Path)
	if strings.Index(file, ".") <= 1 {
		c.Writer.WriteHeader(404)
		c.Writer.Write([]byte("404 page not found"))
		return
	}
	f := http.FileServer(http.Dir(this.staticPath))
	f.ServeHTTP(c.Writer, c.Request)
}

//路由处理，如果路由无法匹配，则进行静态文件处理
func (this *Biu) routeHandle(c *Content, method, path string) {
	handles, find := this.routeParse(method, path, this.route, c)
	if !find {
		handles = append(handles, this.staticHandle)
	}
	var isAbort = false
	var length = len(handles)
	var called = make(map[int]bool) //已调用集合
	var handleIndex = 0
	var callfunc = func(index int) {
		if index >= length || isAbort {
			return
		}
		handleIndex = index
		call := handles[index]
		called[index] = true
		call(c)
	}
	c.nextHandle = func() {
		callfunc(handleIndex + 1)
	}
	c.abortHandle = func() {
		isAbort = true
	}
	for index, _ := range handles {
		if !called[index] && !isAbort {
			callfunc(index)
		}
	}
	return
}

//路由解析
func (this *Biu) routeParse(method, path string, rt *Route, c *Content) (handles []ControlHandle, find bool) {
	//中间件
	for _, call := range rt.middle {
		handles = append(handles, call)
	}
	//节点
	for pt, r := range rt.nodes {
		succ := false
		if pt == path {
			succ = true
		} else if params, _ := r.matching(path); len(params) > 0 {
			c.paramAppend(params)
			succ = true
		}
		if succ && (method == r.method || method == M_ANY) {
			handles = append(handles, r.handle)
			find = true
			return
		}
	}
	//子节点
	for pt, r := range rt.child {
		succ := false
		if strings.Index(path, pt) == 0 {
			succ = true
		} else if params, npath := r.matching(path); len(params) > 0 {
			c.paramAppend(params)
			pt = npath
			succ = true
		}
		if succ {
			path = path[len(pt):]
			childHandles, childFind := this.routeParse(method, path, r, c)
			handles = append(handles, childHandles...)
			find = childFind
			return
		}
	}
	return
}

//启动http服务
func (this *Biu) Run(addr string) error {
	service := &http.Server{
		Addr:           addr,
		Handler:        this,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return service.ListenAndServe()
}
