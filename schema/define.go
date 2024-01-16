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
		"pgsql": funcArgValidPgSQL,
	}
)

//go:linkname funcArgValidMySQL github.com/goodluckxu-go/migrate/schema/mysql.funcArgValid
var funcArgValidMySQL map[string]string

//go:linkname funcArgValidPgSQL github.com/goodluckxu-go/migrate/schema/pgsql.funcArgValid
var funcArgValidPgSQL map[string]string

type EditTable interface {
	mysql.Create | mysql.CreateIfNotExists | mysql.Update |
		pgsql.Create | pgsql.CreateIfNotExists | pgsql.Update
}

type DropTable interface {
	mysql.Drop | mysql.DropIfExists |
		pgsql.Drop | pgsql.DropIfExists
}
