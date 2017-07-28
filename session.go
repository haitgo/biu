//session接口
package biu

import (
	"net/http"
)

type SessionBase interface {
	Server(http.ResponseWriter, *http.Request)
	//session编号
	SessionId() string
	//设置session参数
	Set(interface{}, interface{})
	//读取session参数
	Get(interface{}) interface{}
}
