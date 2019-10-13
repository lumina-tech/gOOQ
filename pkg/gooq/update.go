package gooq

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UpdateSetStep interface {
	UpdateFromStep
	Set(f Field, v interface{}) UpdateSetStep
}

type UpdateFromStep interface {
	UpdateWhereStep
	From(Selectable) UpdateOnStep
}

type UpdateWhereStep interface {
	UpdateOnConflictStep
	Where(conditions ...Expression) UpdateOnConflictStep
}

type UpdateOnStep interface {
	UpdateOnConflictStep
	On(...Expression) UpdateOnConflictStep
}

type UpdateOnConflictStep interface {
	UpdateReturningStep
	OnConflictDoNothing() UpdateReturningStep
	OnConflictDoUpdate() UpdateReturningStep
}

type UpdateReturningStep interface {
	UpdateFinalStep
	Returning(...Expression) UpdateResultStep
}

type UpdateResultStep interface {
	Fetchable
}

type UpdateFinalStep interface {
	Executable
}

///////////////////////////////////////////////////////////////////////////////
// Implementation
///////////////////////////////////////////////////////////////////////////////

type setPredicate struct {
	field Field
	value interface{}
}

type update struct {
	table          Table
	setPredicates  []setPredicate // set predicates
	conditions     []Expression   // where conditions
	fromSelection  Selectable     // selection for from clause
	joinPredicate  []Expression   // join predicates for from clause
	conflictAction ConflictAction
	returning      []Expression
}

func Update(t Table) UpdateSetStep {
	return &update{table: t}
}

func (u *update) Set(field Field, value interface{}) UpdateSetStep {
	u.setPredicates = append(u.setPredicates, setPredicate{field, value})
	return u
}

func (u *update) From(s Selectable) UpdateOnStep {
	u.fromSelection = s
	return u
}

func (u *update) On(c ...Expression) UpdateOnConflictStep {
	u.joinPredicate = c
	return u
}

func (u *update) Where(c ...Expression) UpdateOnConflictStep {
	u.conditions = c
	return u
}

func (u *update) OnConflictDoNothing() UpdateReturningStep {
	u.conflictAction = ConflictActionDoNothing
	return u
}

func (u *update) OnConflictDoUpdate() UpdateReturningStep {
	u.conflictAction = ConflictActionDoUpdate
	panic("not implemented")
	return u
}

func (u *update) Returning(f ...Expression) UpdateResultStep {
	u.returning = f
	return u
}

///////////////////////////////////////////////////////////////////////////////
// Executable
///////////////////////////////////////////////////////////////////////////////

func (u *update) Exec(d Dialect, db DBInterface) (sql.Result, error) {
	//return exec(dl, d, db)
	return nil, nil
}

///////////////////////////////////////////////////////////////////////////////
// Fetchable
///////////////////////////////////////////////////////////////////////////////

func (u *update) Fetch(dl Dialect, db DBInterface) (*sqlx.Rows, error) {
	builder := u.Build(dl)
	return db.Queryx(builder.String(), builder.arguments...)
}

func (u *update) FetchRow(dl Dialect, db DBInterface) *sqlx.Row {
	builder := u.Build(dl)
	return db.QueryRowx(builder.String(), builder.arguments...)
}

///////////////////////////////////////////////////////////////////////////////
// Renderable
///////////////////////////////////////////////////////////////////////////////

func (u *update) Build(d Dialect) *Builder {
	builder := Builder{}
	u.Render(&builder)
	return &builder
}

func (u *update) Render(
	builder *Builder,
) {
	// UPDATE table_name SET
	builder.Printf("UPDATE %s", u.table.GetQualifiedName())

	if len(u.setPredicates) > 0 {
		// render SET clause
		builder.Print(" SET ")
		builder.RenderSetPredicates(u.setPredicates)
		if len(u.conditions) > 0 {
			// render WHERE clause
			builder.Print(" WHERE ")
			builder.RenderConditions(u.conditions)
		}
	}

	if u.fromSelection != nil {
		// render WHERE clause
		builder.Print(" FROM ")
		switch sub := u.fromSelection.(type) {
		case Table:
			builder.Print(sub.GetQualifiedName())
		case *selection:
			builder.Print("(")
			sub.Render(builder)
			builder.Print(")")
		}
		builder.Print(" WHERE ")
		builder.RenderConditions(u.joinPredicate)
	}

	// render on conflict
	if u.conflictAction != ConflictActionNil {
		builder.Printf(" ON CONFLICT %s", string(u.conflictAction))
	}

	// render returning
	if u.returning != nil {
		builder.Print(" RETURNING ")
		builder.RenderProjections(u.returning)
	}
}
