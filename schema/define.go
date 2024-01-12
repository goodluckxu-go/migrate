package schema

import (
	"github.com/goodluckxu-go/migrate/schema/mysql"
	"github.com/goodluckxu-go/migrate/schema/pgsql"
	_ "unsafe"
)

var (
	schemaFuncValid = map[string]int{
		"Edit": 2,
		"Drop": 1,
	}
	funcArgValid = map[string]map[string]string{
		"mysql": funcArgValidMySQL,
	}
)

//go:linkname funcArgValidMySQL github.com/goodluckxu-go/migrate/schema/mysql.funcArgValid
var funcArgValidMySQL map[string]string

type EditTable interface {
	mysql.Create | mysql.CreateIfNotExists | mysql.Update |
		pgsql.Create | pgsql.CreateIfNotExists
}

type DropTable interface {
	mysql.Drop | mysql.DropIfExists |
		pgsql.Drop | pgsql.DropIfExists
}
