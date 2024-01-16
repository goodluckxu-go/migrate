package tb

type Schema struct {
}

// OwnerTo 所有者
func (s *Schema) OwnerTo(owner string) *Schema {
	return s
}

// ClusterOn 集群
func (s *Schema) ClusterOn(cluster string) *Schema {
	return s
}

// Comment 表注释
func (s *Schema) Comment(comment string) *Schema {
	return s
}

// Inherits 继承
func (s *Schema) Inherits(table string) *Schema {
	return s
}

// Fillfactor 填充因子
func (s *Schema) Fillfactor(percentage int) *Schema {
	return s
}
