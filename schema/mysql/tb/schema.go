package tb

type Schema struct {
}

// Engine 引擎
func (s *Schema) Engine(engine string) *Schema {
	return s
}

// Charset 编码
func (s *Schema) Charset(charset string) *Schema {
	return s
}

// Collate 排序规则
func (s *Schema) Collate(collate string) *Schema {
	return s
}

// Comment 表注释
func (s *Schema) Comment(comment string) *Schema {
	return s
}

// AutoIncrement 自动递增
func (s *Schema) AutoIncrement(offset int) *Schema {
	return s
}
