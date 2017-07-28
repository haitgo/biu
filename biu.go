package biu

import (
	"net/http"
	"path"
	"strings"
	"time"
)

type Biu struct {
	session    SessionBase  //session处理接口
	staticPath string       //静态文件目录
	route      *Route       //路由对象
	Server     *http.Server //服务器配置
}

func New() *Biu {
	biu := new(Biu)
	biu.Server = &http.Server{
		Handler:        biu,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return biu
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
func (this *Biu) SessionHandle(handler SessionBase) {
	this.session = handler
}

// ServeHTTP
func (this *Biu) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoratorWriter := newWriter(w)
	c := &Context{
		Writer:  decoratorWriter,
		Request: r,
	}
	if this.session != nil {
		this.session.Server(w, r)
		c.session = this.session
	}
	handles, find := this.routeParse(strings.ToUpper(r.Method), r.URL.Path, this.route, c)
	//如果路由未匹配成功，则试图使用静态文件
	if !find {
		handles = append(handles, this.staticHandle)
	}
	var isAbort = false
	var handlesLength = len(handles)
	var called = make(map[int]bool)  //已调用集合
	var nextFuncIndex = 0            //回调索引编号
	var nextFunc = func(index int) { //下一步回调
		if index >= handlesLength || isAbort {
			return
		}
		nextFuncIndex = index
		call := handles[index]
		called[index] = true
		call(c)
	}
	c.nextHandle = func() {
		nextFunc(nextFuncIndex + 1)
	}
	c.abortHandle = func() {
		isAbort = true
	}
	for index, _ := range handles {
		if !called[index] && !isAbort {
			nextFunc(index)
		}
	}
	return
}

//静态文件处理
func (this *Biu) staticHandle(c *Context) {
	file := path.Base(c.Request.URL.Path)
	if strings.Index(file, ".") <= 1 {
		c.Writer.WriteHeader(404)
		c.Writer.Write([]byte("404 page not found"))
		return
	}
	f := http.FileServer(http.Dir(this.staticPath))
	f.ServeHTTP(c.Writer, c.Request)
}

//路由解析
func (this *Biu) routeParse(method, path string, rt *Route, c *Context) (handles []HandleFunc, find bool) {
	//中间件
	for _, call := range rt.middleware {
		handles = append(handles, call)
	}
	//节点
	var succ bool
	for _, n := range rt.nodes {
		debugPrint("节点", n, "路径", path)
		succ = false //是否匹配成功
		debugPrint("节点查询结果", n.path, path)
		if n.path == path {
			succ = true
		} else if params, _ := n.matching(path); len(params) > 0 {
			c.paramAppend(params)
			succ = true
		}
		if succ && (method == n.method || method == M_ANY) {
			return append(handles, n.handle), true
		}
	}
	//子路由
	for _, r := range rt.childRoute {
		debugPrint("子路由", r)
		succ = false //是否匹配成功
		if strings.Index(path, r.path) == 0 {
			succ = true
		} else if params, npath := r.matching(path); len(params) > 0 {
			c.paramAppend(params)
			r.path = npath
			succ = true
		}
		if succ {
			path = path[len(r.path):]
			subHandles, subFind := this.routeParse(method, path, r, c)
			return append(handles, subHandles...), subFind
		}
	}
	//域名路由
	for _, r := range rt.domainRoute {
		debugPrint("域名路由", r)
		if strings.Index(c.Request.Host, r.domain) >= 0 {
			subHandles, subFind := this.routeParse(method, path, r, c)
			return append(handles, subHandles...), subFind
		}
	}
	return
}

//启动http服务
func (this *Biu) Run(addr string) error {
	this.Server.Addr = addr
	return this.Server.ListenAndServe()
}

//启动https服务
func (this *Biu) RunTLS(addr, certFile, keyFile string) error {
	this.Server.Addr = addr
	return this.Server.ListenAndServeTLS(certFile, keyFile)
}
