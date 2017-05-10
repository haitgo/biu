package validate

import (
	"testing"
)

//创建一个结构体
type Test struct {
	Age  int
	Name string
}

//实现Valider接口
func (this *Test) Match() (m []*M) {
	m = append(m, &M{Attr: "Age", Patten: "Int", Note: "年龄输入错误"})
	m = append(m, &M{Attr: "Name", Call: this.checkName, Note: "姓名输入错误"})
	return m
}
func (this *Test) checkName() bool {
	return this.Name == "wang"
}

//测试
func TestValid(t *testing.T) {
	data := map[string]string{"Age": "28", "Name": "王"}
	test := new(Test)
	v := New(data) //.Omit("Name")
	ok, err := v.Give(test)
	t.Log(ok, err)
	t.Log("对象打印", test)
}
