package tb

type Indexes struct {
}

// Name 名称
func (i *Indexes) Name(name string) *Indexes {
	return i
}

// UNIQUE 唯一索引
func (i *Indexes) UNIQUE() *Indexes {
	return i
}

// BTREE 索引方法
func (i *Indexes) BTREE() *Indexes {
	return i
}

// HASH 索引方法
func (i *Indexes) HASH() *Indexes {
	return i
}

// GiST 索引方法
func (i *Indexes) GiST() *Indexes {
	return i
}

// GIN 索引方法
func (i *Indexes) GIN() *Indexes {
	return i
}

// SPGiST 索引方法
func (i *Indexes) SPGiST() *Indexes {
	return i
}

// BRIN 索引方法
func (i *Indexes) BRIN() *Indexes {
	return i
}

// Comment 注释
func (i *Indexes) Comment(column string) *Indexes {
	return i
}
