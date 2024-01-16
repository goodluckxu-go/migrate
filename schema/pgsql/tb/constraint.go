package tb

type Constraint struct {
}

// Name 外键名称
func (c *Constraint) Name(name string) *Constraint {
	return c
}

// QuoteTable 引用表
func (c *Constraint) QuoteTable(table string) *Constraint {
	return c
}

// QuoteColumn 引用字段
func (c *Constraint) QuoteColumn(column string, columns ...string) *Constraint {
	return c
}

// DeleteCascade 删除同步删除
func (c *Constraint) DeleteCascade() *Constraint {
	return c
}

// UpdateCascade 更新同步更新
func (c *Constraint) UpdateCascade() *Constraint {
	return c
}

// Comment 注释
func (c *Constraint) Comment(column string) *Constraint {
	return c
}
