package migrate

import (
	_ "github.com/goodluckxu-go/migrate/schema"
	"go/token"
	_ "unsafe"
)

const (
	modName = "github.com/goodluckxu-go/migrate" // mod名称
)

// 表ast解析
type tableAst struct {
	Type       string      // 表类型
	Name       string      // 表名称
	Active     string      // 表动作
	Func       string      // 调用方法
	ColumnList []columnAst // 内部字段
}

// 字段ast解析
type columnAst struct {
	LianFuncSort []string            // 链式方法调用顺序
	InternetFunc map[string][]argAst // 内部链式方法
}

// 链表方法参数ast解析
type argAst struct {
	Val  interface{} // 字段值
	Type string      // 字段类型
	Pos  token.Pos   // 位置
}

//go:linkname schemaFuncValid github.com/goodluckxu-go/migrate/schema.schemaFuncValid
var schemaFuncValid map[string]int

//go:linkname funcArgValid github.com/goodluckxu-go/migrate/schema.funcArgValid
var funcArgValid map[string]map[string]string
