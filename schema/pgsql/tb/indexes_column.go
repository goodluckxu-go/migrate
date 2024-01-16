package tb

type IndexesColumn struct {
}

// Field 字段
func (i *IndexesColumn) Field(column string) *IndexesColumn {
	return i
}

// CollateMode 排序规则模式
func (i *IndexesColumn) CollateMode(mode string) *IndexesColumn {
	return i
}

// Collate 排序规则
func (i *IndexesColumn) Collate(collate string) *IndexesColumn {
	return i
}

// OperationSymbolMode 运算符号类别模式
func (i *IndexesColumn) OperationSymbolMode(mode string) *IndexesColumn {
	return i
}

// OperationSymbol 运算符号类别
func (i *IndexesColumn) OperationSymbol(mode string) *IndexesColumn {
	return i
}

// ASC 排序升序
func (i *IndexesColumn) ASC() *IndexesColumn {
	return i
}

// DESC 排序降序
func (i *IndexesColumn) DESC() *IndexesColumn {
	return i
}

// NULLSLAST NULLS排序last
func (i *IndexesColumn) NULLSLAST() *IndexesColumn {
	return i
}

// NULLSFIRST NULLS排序first
func (i *IndexesColumn) NULLSFIRST() *IndexesColumn {
	return i
}
