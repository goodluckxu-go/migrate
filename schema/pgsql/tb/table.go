package tb

type UpdateTable struct {
}

// AddColumn 添加字段
func (t *UpdateTable) AddColumn(column string) *Types {
	return new(Types)
}

// ModifyColumn 修改字段
func (t *UpdateTable) ModifyColumn(column string) *Types {
	return new(Types)
}

// ChangeColumn 重命名字段
func (t *UpdateTable) ChangeColumn(oldColumn, newColumn string) *Types {
	return new(Types)
}

type CreateTable struct {
	Schema *Schema
}

// Column 创建字段
func (t *CreateTable) Column(column string) *Types {
	return new(Types)
}

// Indexes 创建字段
func (t *CreateTable) Indexes(columns func(column *IndexesColumn)) *Indexes {
	return new(Indexes)
}

// PrimaryKey 主键
func (t *CreateTable) PrimaryKey(column string, columns ...string) {
}

// ForeignKey 外键
func (t *CreateTable) ForeignKey(column string, columns ...string) *Constraint {
	return new(Constraint)
}
