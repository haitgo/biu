package bind

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

//不赋值
const (
	NOT_SET  = "NOT_SET"
	NOT_NULL = "不能为空"
)

/*
说明：验证器是自动将map[string]interface{}的值赋对象，并对每个值的格式进行验证匹配
以保证数据输入的准确性。
使用方法：
	//创建一个结构体
	type Test struct{
		Age string
	}
	//实现Valider接口
	func (this *Test)BindMatch()(m []*validate.M){
		m=append(m,bind.M{Attr:"Age", Check:this.checkName,Note:"输入错误"})
	}
	func (this *Test)checkName(name interface{})bool{
		return false
	}
	data:=map[string]interface{}{"Age":"28"}
	t:=new(Test)
	v:=validate.New(data)
	ok,err:=v.Give(t)
	fmt.Println(ok,err)
*/

//验证器接口
type Valider interface {
	BindMatch() []*M
}

//回调方法
type CallFunc func(string) bool
type BindCall func() bool

//正则式注册结构
type M struct {
	Attr  string      //字段名称
	Check interface{} //数据验证 string|BindCall
	Note  string      //错误说明
}

//数据绑定
type Bind struct {
	Datas  map[string]string
	matchs map[string][]*M
	must   map[string]bool
	omit   map[string]bool
}

//
func New(datas map[string]string) *Bind {
	o := new(Bind)
	o.Datas = datas
	o.matchs = make(map[string][]*M)
	o.must = make(map[string]bool)
	o.omit = make(map[string]bool)
	return o
}

//增加必须验证的属性
func (this *Bind) Must(attr ...string) *Bind {
	for _, a := range attr {
		this.must[a] = true
	}
	return this
}

//增加忽略验证的属性
func (this *Bind) Omit(attr ...string) *Bind {
	for _, a := range attr {
		this.omit[a] = true
	}
	return this
}

//给o赋值并验证数据
func (this *Bind) Give(obj interface{}) (err error, attr string) {
	//验证集合
	if v, ok := obj.(Valider); ok {
		for _, m := range v.BindMatch() {
			if this.matchs[m.Attr] == nil {
				this.matchs[m.Attr] = make([]*M, 0)
			}
			this.matchs[m.Attr] = append(this.matchs[m.Attr], m)
		}
	}
	ref := reflect.ValueOf(obj).Elem()
	typ := ref.Type()
	return this.reflectSet(ref, typ)
}

//反射赋值
func (this *Bind) reflectSet(ref reflect.Value, typ reflect.Type) (err error, attr string) {
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Kind() == reflect.Struct { //如果为结构体，则循环从结构体里面去绑定参数
			return this.reflectSet(ref.Field(i), typ.Field(i).Type)
		}
		field := typ.Field(i)       //
		attr = field.Name           //结构体成员属性名称
		if !ref.Field(i).CanSet() { //如果该字段不允许赋值
			return errors.New("没有权限"), attr
		}
		//类型断言转换
		value, err := this.typeConversion(this.Datas[attr], field.Type.Name())
		if err != nil { //如果断言转换类型出现错误
			return err, attr
		}
		ok, note := this.match(attr, this.Datas[attr])
		if ok {
			ref.Field(i).Set(value) //给结构体赋值
		} else {
			switch note {
			case NOT_SET:
				continue
			case NOT_NULL:
				return errors.New(note), NOT_NULL
			default:
				return errors.New(note), attr
			}
		}
	}
	return nil, ""
}

//类型断言
func (this *Bind) typeConversion(value string, ntype string) (reflect.Value, error) {
	intvalue := value
	if intvalue == "" {
		intvalue = "0"
	}
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.ParseInt(intvalue, 10, 64)
		return reflect.ValueOf(int(i)), err
	} else if ntype == "uint" {
		i, err := strconv.ParseInt(intvalue, 10, 64)
		return reflect.ValueOf(uint(i)), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(intvalue, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int16" {
		i, err := strconv.ParseInt(intvalue, 10, 64)
		return reflect.ValueOf(int16(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(intvalue, 10, 64)
		return reflect.ValueOf(int32(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(intvalue, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(intvalue, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(intvalue, 64)
		return reflect.ValueOf(i), err
	}
	//else if .......增加其他一些类型的转换
	return reflect.ValueOf(value), errors.New("未知的类型" + ntype)
}

//单项验证，如果没有需要该字段的验证，则忽略过去不映射到对应的结构体属性上去
func (this *Bind) match(attr string, v string) (ok bool, err string) {
	matchs := this.matchs[attr]
	if matchs == nil {
		return false, NOT_SET
	}
	if this.must[attr] && v == "" {
		return false, NOT_NULL
	}
	if this.omit[attr] || v == "" {
		return false, NOT_SET
	}
	check := true
	for _, m := range matchs {
		switch m.Check.(type) {
		case string:
			check = Match(v, m.Check.(string))
		case BindCall:
			call := m.Check.(BindCall)
			check = call()
		}
		if !check {
			return false, m.Note
		}
	}
	return check, ""
}
