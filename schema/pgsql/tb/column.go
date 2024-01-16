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

// Dimension 维度，几维数组
func (c *Column) Dimension(length int) *Column {
	return c
}

// Comment 注释
func (c *Column) Comment(column string) *Column {
	return c
}

// CollateMode 排序规则模式
func (c *Column) CollateMode(mode string) *Column {
	return c
}

// Collate 排序规则
func (c *Column) Collate(collate string) *Column {
	return c
}

// GeneratedAlwaysAsIdentity 生成always标志
// 只能系统生成，无法被用户修改覆盖
// maxvalue不填则根据类型生成最大值
func (c *Column) GeneratedAlwaysAsIdentity(increment, start, minvalue int64, maxvalue ...int64) *Column {
	return c
}

// GeneratedByDefaultAsIdentity 生成default标志
// 系统生成，可被用户修改覆盖
// maxvalue不填则根据类型生成最大值
func (c *Column) GeneratedByDefaultAsIdentity(increment, start, minvalue int64, maxvalue ...int64) *Column {
	return c
}
