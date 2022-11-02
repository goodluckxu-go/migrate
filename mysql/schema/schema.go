package schema

import (
	"github.com/goodluckxu/migrate/mysql/schema/tb"
)

// Table 修改表
func Table(table string, fn func(table tb.Table)) {
}

// Table 创建表
func Create(table string, fn func(table tb.CreateTable)) *tb.Schema {
	return new(tb.Schema)
}

// CreateIfNotExists 表不存在创建表
func CreateIfNotExists(table string, fn func(table tb.CreateTable)) *tb.Schema {
	return new(tb.Schema)
}

// Drop 删除表
func Drop(table string) {
}

// DropIfExists 表存在则删除
func DropIfExists(table string) {
}
