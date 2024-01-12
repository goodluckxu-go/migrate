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

// DropColumn 删除字段
func (t *UpdateTable) DropColumn(column string) {
}

// AddIndexes 添加索引
func (t *UpdateTable) AddIndexes(columns ...string) *Indexes {
	return new(Indexes)
}

// DropIndexes 删除索引
func (t *UpdateTable) DropIndexes(name string) {
}

// AddPrimaryKey 添加主键
func (t *UpdateTable) AddPrimaryKey(columns ...string) {
}

// DropPrimaryKey 删除主键
func (t *UpdateTable) DropPrimaryKey() {
}

// AddForeignKey 添加外键
func (t *UpdateTable) AddForeignKey(columns ...string) *Constraint {
	return new(Constraint)
}

// DropForeignKey 删除外键
func (t *UpdateTable) DropForeignKey(name string) {
}

type CreateTable struct {
	Schema *Schema
}

// Column 创建字段
func (t *CreateTable) Column(column string) *Types {
	return new(Types)
}

// Indexes 创建字段
func (t *CreateTable) Indexes(column string, columns ...string) *Indexes {
	return new(Indexes)
}

// PrimaryKey 主键
func (t *CreateTable) PrimaryKey(column string, columns ...string) {
}

// ForeignKey 外键
func (t *CreateTable) ForeignKey(column string, columns ...string) *Constraint {
	return new(Constraint)
}
