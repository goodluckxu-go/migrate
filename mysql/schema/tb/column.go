package tb

type Column struct {
}

// Nullable 允许为空
func (c *Column) Nullable() *Column {
	return c
}

// Default 默认值
func (c *Column) Default(value interface{}) *Column {
	return c
}

// Unsigned 无符号
func (c *Column) Unsigned() *Column {
	return c
}

// Comment 注释
func (c *Column) Comment(column string) *Column {
	return c
}

// Charset 编码
func (c *Column) Charset(charset string) *Column {
	return c
}

// Collate 排序规则
func (c *Column) Collate(collate string) *Column {
	return c
}

// AutoIncrement 自动递增
func (c *Column) AutoIncrement() *Column {
	return c
}

// Zerofill 填充零
func (c *Column) Zerofill() *Column {
	return c
}

// After 在某个字段之后
func (c *Column) After(column string) *Column {
	return c
}

// First 显示第一个字段
func (c *Column) First() *Column {
	return c
}
