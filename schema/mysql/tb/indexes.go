package tb

type Indexes struct {
}

// Name 名称
func (i *Indexes) Name(name string) *Indexes {
	return i
}

// FULLTEXT 文本索引
func (i *Indexes) FULLTEXT() *Indexes {
	return i
}

// UNIQUE 唯一索引
func (i *Indexes) UNIQUE() *Indexes {
	return i
}

// SPATIAL 空间索引
func (i *Indexes) SPATIAL() *Indexes {
	return i
}

// BTREE 存储方法
func (i *Indexes) BTREE() *Indexes {
	return i
}

// HASH 存储方法
func (i *Indexes) HASH() *Indexes {
	return i
}
