package bind

import (
	"testing"
)

//创建一个结构体
type Test struct {
	Age  int
	Name string
}

//实现Valider接口
func (this *Test) BindMatch() (m []*M) {
	m = append(m, &M{"Age", "Int", "年龄输入错误"})
	m = append(m, &M{"Name", this.checkName, "姓名输入错误"})
	return m
}
func (this *Test) checkName() bool {
	return this.Name == "wang"
}

//测试
func TestBind(t *testing.T) {
	data := map[string]string{"Age": "32", "Name": "王"}
	test := new(Test)
	v := New(data) //.Omit("Age")
	ok, err := v.Give(test)
	t.Log(ok, err)
	t.Log("对象打印", test)
}
