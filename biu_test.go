package biu

import (
	"fmt"
	"testing"
)

func call(c *Content) {
	fmt.Println("第一个测试")
}
func call2(c *Content) {
	fmt.Println("表示成功啦")
}
func middle1(c *Content) {
	fmt.Println("中间件开始1")
	c.Next()
	fmt.Println("中间件结束1")
}
func middle2(c *Content) {
	fmt.Println("中间件开始2")
	c.Next()
	fmt.Println("中间件结束2")
}
func TestBiu(t *testing.T) {
	b := New()
	rt := b.Route()
	//rt.Middleware(middle1)
	rt.Get("/aaa", call)
	rt.Get("/bbb", call)
	g := rt.Group("/bbb").Match("name", "水电费")
	{
		g.Middleware(middle2)
		g.Get("/ddd", call2)
	}
	b.Run(":992")
}
