package meta

type TypeInfo struct {
	Prefix  string
	Literal string
}

var Types = []TypeInfo{
	TypeInfo{Prefix: "Bool", Literal: "bool"},
	TypeInfo{Prefix: "Float32", Literal: "float32"},
	TypeInfo{Prefix: "Float64", Literal: "float64"},
	TypeInfo{Prefix: "Int", Literal: "int"},
	TypeInfo{Prefix: "Int64", Literal: "int64"},
	TypeInfo{Prefix: "Jsonb", Literal: "[]byte"},
	TypeInfo{Prefix: "String", Literal: "string"},
	TypeInfo{Prefix: "StringArray", Literal: "pq.StringArray"},
	TypeInfo{Prefix: "Time", Literal: "time.Time"},
	TypeInfo{Prefix: "UUID", Literal: "uuid.UUID"},

	TypeInfo{Prefix: "NullBool", Literal: "null.Bool"},
	TypeInfo{Prefix: "NullFloat32", Literal: "null.Float"},
	TypeInfo{Prefix: "NullFloat64", Literal: "null.Float"},
	TypeInfo{Prefix: "NullInt", Literal: "null.Int"},
	TypeInfo{Prefix: "NullInt64", Literal: "null.Int"},
	TypeInfo{Prefix: "NullJsonb", Literal: "nullable.Jsonb"},
	TypeInfo{Prefix: "NullString", Literal: "null.String"},
	TypeInfo{Prefix: "NullTime", Literal: "null.Time"},
	TypeInfo{Prefix: "NullUUID", Literal: "nullable.UUID"},
}

type FunctionInfo struct {
	Name string
	Expr string
}

var Funcs = []FunctionInfo{
	FunctionInfo{Name: "ArrayAgg", Expr: "ARRAY_AGG(%s)"},
	FunctionInfo{Name: "Asc", Expr: "%s ASC"},
	FunctionInfo{Name: "Avg", Expr: "AVG(%s)"},
	FunctionInfo{Name: "Cast", Expr: "CAST(%s AS %s)"},
	FunctionInfo{Name: "Ceil", Expr: "CEIL(%s)"},
	FunctionInfo{Name: "Coalesce", Expr: "COALESCE(%s, %v)"},
	FunctionInfo{Name: "Count", Expr: "COUNT(%s)"},
	FunctionInfo{Name: "Date", Expr: "DATE(%s)"},
	FunctionInfo{Name: "DateTrunc", Expr: "DATE_TRUNC('%s', %v)"},
	FunctionInfo{Name: "Desc", Expr: "%s DESC"},
	FunctionInfo{Name: "Distinct", Expr: "DISTINCT(%s)"},
	FunctionInfo{Name: "Div", Expr: "%s / %v"},
	FunctionInfo{Name: "Mult", Expr: "%s * %v"},
	FunctionInfo{Name: "Hex", Expr: "HEX(%s)"},
	FunctionInfo{Name: "Lower", Expr: "LOWER(%s)"},
	FunctionInfo{Name: "Max", Expr: "MAX(%s)"},
	FunctionInfo{Name: "Md5", Expr: "MD5(%s)"},
	FunctionInfo{Name: "Min", Expr: "MIN(%s)"},
	FunctionInfo{Name: "First", Expr: "first(%s, %s)"},
	FunctionInfo{Name: "Last", Expr: "last(%s, %s)"},
	FunctionInfo{Name: "Substr2", Expr: "SUBSTR(%s, %v)"},
	FunctionInfo{Name: "Substr3", Expr: "SUBSTR(%s, %v, %v)"},
	FunctionInfo{Name: "Sum", Expr: "SUM(%s)"},
	FunctionInfo{Name: "Upper", Expr: "UPPER(%s)"},
}
