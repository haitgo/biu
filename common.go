package biu

import (
	"fmt"
)

//数据结构
type H map[string]interface{}

func debugPrint(data ...interface{}) {

	bata := make([]interface{}, 0)
	bata = append(bata, "[biu]")
	bata = append(bata, data...)
	fmt.Println(bata...)
}
