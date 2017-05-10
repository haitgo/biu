package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/haitgo/biu"
)

func test() {
	aaa := gin.New()
	aaa.GET("/a", func(c *gin.Context) {
		c.Writer.Status()
	})
}
func call(c *biu.Content) {
	c.Writer.Write([]byte("第一个测试"))
	fmt.Println("第一个测试")
}
func call2(c *biu.Content) {
	c.Writer.Write([]byte("哈哈哈"))
	fmt.Println("参数是", c.Params)
	fmt.Println(c.Query("Name"))
	fmt.Println("表示成功啦")
}
func middle1(c *biu.Content) {
	fmt.Println("中间件开始1")
	c.Next()
	fmt.Println("中间件结束1")
}
func middle2(c *biu.Content) {

	c.Writer.ObStart()
	c.Writer.Write([]byte("中间件开始2"))
	// c.Upload("aa").GetFile()
	//c.Abort()
	c.Next()
	c.Writer.Write([]byte("中间件结束2"))
	c.Writer.ObFlush()
}
func main() {
	b := biu.New()
	b.StaticPath("tmp")
	b.Route(func(rt *biu.Route) {
		rt.Middleware(middle1)
		rt.Get("/aaa", call)
		rt.Get("/bbb", call)
		g := rt.Group("/bbb")
		{
			g.Get(`/ccc`, call2)
			g.Middleware(middle2)
			g.Get(`/{age}.html`, call2).Match("age", `[\d]{5}`)
		}
	})
	b.Run(":992")
}
