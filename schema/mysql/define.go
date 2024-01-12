package mysql

import (
	"github.com/goodluckxu-go/migrate/schema/mysql/tb"
)

var (
	funcArgValid = map[string]string{
		"Create":            "tb.CreateTable",
		"CreateIfNotExists": "tb.CreateTable",
		"Update":            "tb.UpdateTable",
	}
)

type Create func(table tb.CreateTable)

type CreateIfNotExists func(table tb.CreateTable)

type Update func(table tb.UpdateTable)

type Drop uint8

type DropIfExists uint8
