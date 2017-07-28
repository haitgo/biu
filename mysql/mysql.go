package mysql

func Connect(link string) *Session {
	c := new(Conn)
	return c

}

//事务处理回调
func Trans(call func()) {
	Trans(func(c *Conn) error {

	})
}

type Session struct {
	Link     string
	andSql   []string
	orSql    []string
	defSql   []string
	orderSql []string
	giveObj  interface{}
}

//查询字段
func (this *Session) Field() *Session {
	return this
}

//填充数据（insert和update时使用）
func (this *Session) Data() *Session {
	return this
}

//条件
func (this *Session) Where() *Session {
	return this
}

//条件
func (this *Session) And(sql string) *Session {
	this.andSql = append(this.andSql, sql)
	return this
}

//条件
func (this *Session) Or(sql string) *Conn {
	this.orSql = append(this.orSql, sql)
	return this
}

//排序
func (this *Session) Order(order string) *Session {
	this.orderSql = append(this.orderSql, order)
	return this
}

//查询导出,将查询出来的结果赋值给obj对象
func (this *Session) Give(obj ...interface{}) {
	this.giveObj = obj
}

//查询
func (this *Session) Select() {

}

//查询one
func (this *Session) SelectOne() {

}

//查询并分页
func (this *Session) SelectPage() {

}

//插入数据
func (this *Session) Insert() {

}

//更新数据
func (this *Session) Update() {

}

//修改一条数据
func (this *Session) UpdateOne() {

}

//删除数据
func (this *Session) Delete() {

}

//删除一条数据
func (this *Session) DeleteOne() {

}

//sql执行
func (this *Session) Query(sql string) {

}
