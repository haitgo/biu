package biu

import (
	"fmt"
)

//数据结构
type H map[string]interface{}

func echo(data ...interface{}) {
	bata := make([]interface{}, 0)
	bata = append(bata, "[biu-debug]")
	bata = append(bata, data...)
	fmt.Println(bata...)
}
