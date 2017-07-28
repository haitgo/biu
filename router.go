//路由处理
package biu

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

//请求方式集合
const (
	M_GET    = "GET"
	M_POST   = "POST"
	M_PUT    = "PUT"
	M_DELETE = "DELETE"
	M_ANY    = "ANY"
)

//路由匹配-----------------------------------------------------------------------
type pathMatch struct {
	needMatch bool              //是否需要匹配
	matchPath string            //路径
	pattens   map[string]string //正则式
	match     *regexp.Regexp    //路由参数match对象
}

//添加匹配路径
func (this *pathMatch) addMatchPath(path string) {
	this.matchPath = strings.Replace(path, ".", "\\.", -1)
	this.pattens = make(map[string]string)
	if strings.Index(path, "{") > 0 && strings.Index(path, "}") > 0 {
		this.needMatch = true
	}
}

//添加匹配规则
func (this *pathMatch) addMatchPatten(name, patten string) {
	this.pattens[name] = patten
}

//正则路由匹配，返回查询参数和新路径
func (this *pathMatch) matching(path string) (params map[string]string, newPath string) {
	if !this.needMatch {
		return
	}
	var err error
	if this.match == nil {
		pathPatten := "^" + this.matchPath
		matchName, _ := regexp.Compile(`\{([\w]+)\}`)
		names := matchName.FindAllStringSubmatch(pathPatten, -1)
		if names == nil {
			return
		}
		for _, eq := range names {
			name := eq[1]
			pat := fmt.Sprintf(`(?P<%s>[\d\w_\.]+)`, name)
			if patten, ok := this.pattens[name]; ok {
				pat = fmt.Sprintf(`(?P<%s>%s)`, name, patten)
			}
			pathPatten = strings.Replace(pathPatten, eq[0], pat, -1)
		}
		this.match, err = regexp.Compile(pathPatten)
		if err != nil {
			log.Fatal("[biu]", err)
		}
		//log.Println("正则", path, "\t", pathPatten)
	}
	params = make(map[string]string)
	match := this.match.FindStringSubmatch(path)
	if match == nil {
		return
	}
	newPath = this.matchPath
	for i, name := range this.match.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		params[name] = match[i]
		newPath = strings.Replace(newPath, "{"+name+"}", match[i], -1)
	}
	return
}

//路由节点------------------------------------------------------------------------
type node struct {
	method    string     //请求方法
	path      string     //路径（截取后的段路径)
	handle    HandleFunc //控制器处理方法
	pathMatch            //继承路径匹配对象
}

func (this *node) Match(name, patten string) *node {
	this.addMatchPatten(name, patten)
	return this
}

//回调-----------------------------------------------------------------------
type HandleFunc func(*Context)

//路由-----------------------------------------------------------------------
type Route struct {
	path        string       //路径（截取后的段路径)
	domain      string       //域名设置
	middleware  []HandleFunc //中间件
	nodes       []*node      //终节点
	childRoute  []*Route     //子路由
	domainRoute []*Route     //域名路由
	pathMatch                //继承路径匹配对象
}

//创建路由
func newRoute() *Route {
	r := new(Route)
	r.nodes = make([]*node, 0)
	r.childRoute = make([]*Route, 0)
	r.domainRoute = make([]*Route, 0)
	r.middleware = make([]HandleFunc, 0)
	return r
}

func (this *Route) Match(name, patten string) *Route {
	this.addMatchPatten(name, patten)
	return this
}

//get请求
func (this *Route) Get(path string, call HandleFunc) *node {
	return this.addNode(M_GET, path, call)
}

//post请求
func (this *Route) Post(path string, call HandleFunc) *node {
	return this.addNode(M_POST, path, call)
}

//put请求
func (this *Route) Put(path string, call HandleFunc) *node {
	return this.addNode(M_PUT, path, call)
}

//del请求
func (this *Route) Delete(path string, call HandleFunc) *node {
	return this.addNode(M_DELETE, path, call)
}

//任意请求
func (this *Route) Any(path string, call HandleFunc) *node {
	return this.addNode(M_ANY, path, call)
}

//添加路由节点
func (this *Route) addNode(method, path string, call HandleFunc) *node {
	n := &node{method: method, handle: call}
	n.addMatchPath(path)
	n.path = path
	this.nodes = append(this.nodes, n)
	return n
}

//域名限制,可输入多个域名
func (this *Route) Domain(domain string) *Route {
	rt := newRoute()
	rt.domain = domain
	this.domainRoute = append(this.domainRoute, rt)
	return rt
}

//路由分组
func (this *Route) Group(path string) *Route {
	rt := newRoute()
	rt.addMatchPath(path)
	rt.path = path
	this.childRoute = append(this.childRoute, rt)
	return rt
}

//中间件
func (this *Route) Middleware(call ...HandleFunc) {
	this.middleware = append(this.middleware, call...)
}
