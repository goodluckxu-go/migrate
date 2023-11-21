package tb

type Table struct {
}

// AddColumn 添加字段
func (t *Table) AddColumn(column string) *Types {
	return new(Types)
}

// ModifyColumn 修改字段
func (t *Table) ModifyColumn(column string) *Types {
	return new(Types)
}

// ChangeColumn 重命名字段
func (t *Table) ChangeColumn(oldColumn, newColumn string) *Types {
	return new(Types)
}

// DropColumn 删除字段
func (t *Table) DropColumn(column string) {
}

// AddIndexes 添加索引
func (t *Table) AddIndexes(columns ...string) *Indexes {
	return new(Indexes)
}

// DropIndexes 删除索引
func (t *Table) DropIndexes(name string) {
}

// AddPrimaryKey 添加主键
func (t *Table) AddPrimaryKey(columns ...string) {
}

// DropPrimaryKey 删除主键
func (t *Table) DropPrimaryKey() {
}

// AddForeignKey 添加外键
func (t *Table) AddForeignKey(columns ...string) *Constraint {
	return new(Constraint)
}

// DropForeignKey 删除外键
func (t *Table) DropForeignKey(name string) {
}

type CreateTable struct {
}

// Column 创建字段
func (t *CreateTable) Column(column string) *Types {
	return new(Types)
}

// Indexes 创建字段
func (t *CreateTable) Indexes(columns ...string) *Indexes {
	return new(Indexes)
}

// PrimaryKey 主键
func (t *CreateTable) PrimaryKey(columns ...string) {
}

// ForeignKey 外键
func (t *CreateTable) ForeignKey(columns ...string) *Constraint {
	return new(Constraint)
}
