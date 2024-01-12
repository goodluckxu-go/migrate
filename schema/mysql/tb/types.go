package tb

type Types struct {
}

// Tinyint 类型
func (t *Types) Tinyint(length ...int) *Column {
	return new(Column)
}

// Smallint 类型
func (t *Types) Smallint(length ...int) *Column {
	return new(Column)
}

// Int 类型
func (t *Types) Int(length ...int) *Column {
	return new(Column)
}

// Bigint 类型
func (t *Types) Bigint(length ...int) *Column {
	return new(Column)
}

// Decimal 类型
func (t *Types) Decimal(length, point int) *Column {
	return new(Column)
}

// Float 类型
func (t *Types) Float(length, point int) *Column {
	return new(Column)
}

// Char 类型
func (t *Types) Char(length int) *Column {
	return new(Column)
}

// Varchar 类型
func (t *Types) Varchar(length int) *Column {
	return new(Column)
}

// Tinytext 类型
func (t *Types) Tinytext() *Column {
	return new(Column)
}

// Mediumtext 类型
func (t *Types) Mediumtext() *Column {
	return new(Column)
}

// Text 类型
func (t *Types) Text() *Column {
	return new(Column)
}

// Longtext 类型
func (t *Types) Longtext() *Column {
	return new(Column)
}

// Date 类型
func (t *Types) Date() *Column {
	return new(Column)
}

// Time 类型
func (t *Types) Time(length ...int) *Column {
	return new(Column)
}

// Year 类型
func (t *Types) Year() *Column {
	return new(Column)
}

// Datetime 类型
func (t *Types) Datetime(length ...int) *Column {
	return new(Column)
}

// Timestamp 类型
func (t *Types) Timestamp(length ...int) *Column {
	return new(Column)
}

// Enum 类型
func (t *Types) Enum(args ...string) *Column {
	return new(Column)
}

// Json 类型
func (t *Types) Json() *Column {
	return new(Column)
}
