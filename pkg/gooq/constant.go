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

type LockingType int

const (
	LockingTypeNone LockingType = iota
	LockingTypeUpdate
	LockingTypeNoKeyUpdate
	LockingTypeShare
	LockingTypeKeyShare
)

func (t LockingType) String() string {
	switch t {
	case LockingTypeUpdate:
		return "FOR UPDATE"
	case LockingTypeNoKeyUpdate:
		return "FOR NO KEY UPDATE"
	case LockingTypeShare:
		return "FOR SHARE"
	case LockingTypeKeyShare:
		return "FOR KEY SHARE"
	default:
		return ""
	}
}

type LockingOption int

const (
	LockingOptionNone LockingOption = iota
	LockingOptionNoWait
	LockingOptionSkipLocked
)

func (t LockingOption) String() string {
	switch t {
	case LockingOptionNoWait:
		return "NOWAIT"
	case LockingOptionSkipLocked:
		return "SKIP LOCKED"
	default:
		return ""
	}
}
