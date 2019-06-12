package gooq

type CommitAction string

const (
	CommitActionDeleteRows   CommitAction = "DELETE ROWS"
	CommitActionDrop         CommitAction = "DROP"
	CommitActionPreserveRows CommitAction = "PRESERVE ROWS"
)
