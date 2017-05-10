package upload

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func Test_HTTP(t *testing.T) {
	g := gin.Default()
	g.GET("/", func(c *gin.Context) {
		html := `<html>
		<body>
			<form method="POST" enctype="multipart/form-data"/>
				<input type="file" name="file" />
				<input type="text" name="test" />
				<input type="submit" value="提交" />
			</form>
		</body>
		</html>`
		c.Header("content-type", "text/html;charset=utf-8")
		c.Header("Cache-Control", "max-age=10")
		c.String(200, html)
	})
	g.POST("/", func(c *gin.Context) {
		cut, err := NewUpload(c.Request).AllowExt(".jpg", ".png").Image("file") //上传的表单为file的图片
		if err != nil {
			fmt.Println(err) //上传失败
		} else {
			tm := strconv.Itoa(int(time.Now().Unix()))
			err := cut.Resize(0, 0).WriteFile("abc/" + tm + ".jpg")
			//err := cut.WriteFile("abc/a.jpg")
			//c.String(200, err.Error())
			fmt.Println(err)
		}
	})
	g.Run(":992")
}
