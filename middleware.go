package biu

import (
	"strconv"
	"time"
)

//调试输出
func Debug() HandleFunc {
	return func(c *Context) {
		time1 := time.Now().UnixNano()
		c.Next()
		time2 := time.Now().UnixNano()

		runeMethod := make([]rune, 8)
		copy(runeMethod, []rune(c.Request.Method))

		scode := strconv.Itoa(c.Writer.Status())
		runeCode := make([]rune, 4)
		copy(runeCode, []rune(scode))

		useTime := (time2 - time1) / 1000000
		runeTime := make([]rune, 10)
		useTimeStr := strconv.Itoa(int(useTime)) + "ms"
		copy(runeTime, []rune(useTimeStr))
		echo(string(runeMethod), string(runeCode), string(runeTime), "->", c.Request.URL.Path)
	}
}
