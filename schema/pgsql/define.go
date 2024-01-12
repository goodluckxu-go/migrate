package pgsql

type Create func(table int)

type CreateIfNotExists func(table int)

type Drop uint8

type DropIfExists uint8
