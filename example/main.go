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
func domainA(c *biu.Context) {
	c.Writer.Write([]byte("domain a."))
}
func domainB(c *biu.Context) {
	c.Writer.Write([]byte("domain b."))
}
func call(c *biu.Context) {
	c.Writer.Write([]byte("第一个测试"))
	fmt.Println("第一个测试")
}
func call2(c *biu.Context) {
	c.Writer.Write([]byte(c.Request.Host))
	c.Writer.Write([]byte("哈哈哈"))
	fmt.Println("参数是", c.Params)
	fmt.Println(c.Query("Name"))
}
func middle1(c *biu.Context) {
	fmt.Println("中间件开始1")

	c.Next()
	fmt.Println("中间件结束1")
}
func middle2(c *biu.Context) {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println("错啦", er)
		}
	}()
	c.Writer.OBstart()
	c.Writer.Write([]byte("中间件开始2"))
	// c.Upload("aa").GetFile()
	//c.Abort()
	c.Next()
	c.Writer.Write([]byte("中间件结束2"))
	c.Writer.OBflush()
}
func main() {
	b := biu.New()
	b.Server.WriteTimeout = 1
	b.StaticPath("tmp")
	b.Route(func(rt *biu.Route) {
		at := rt.Domain("127.0.0.1")
		{
			at.Get("/test", domainA)
		}
		bt := rt.Domain("localhost")
		{
			bt.Get("/test", domainB)
		}
		rt.Get("/aaa", call)
		rt.Get("/bbb.html", call)
		g := rt.Group("/bbb")
		{
			g.Get(`/ccc`, call2)
			g.Middleware(middle2)
			g.Get(`/{age}.html`, call2).Match("age", `[\d]{5}`)
		}
	})
	b.Run(":992")
}
