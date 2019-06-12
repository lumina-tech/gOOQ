package gooq

// Base on https://www.jooq.org/javadoc/latest/

type ConflictAction string
type Dialect int
type JoinType int

const (
	Sqlite Dialect = iota
	MySQL
	Postgres
)

const (
	ConflictActionNil                      = ConflictAction("")
	ConflictActionDoNothing ConflictAction = "DO NOTHING"
	ConflictActionDoUpdate  ConflictAction = "DO UPDATE"
)

const (
	Join JoinType = iota
	LeftOuterJoin
	NotJoined
)
