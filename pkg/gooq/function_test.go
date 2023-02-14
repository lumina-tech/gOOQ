package gooq

import "testing"

var functionTestCases = []TestCase{
	{
		Constructed:  And(Table1.Column1.Eq(Table2.Column1), Table1.Column2.Eq(Table2.Column2), Table1.Column2.Eq(Table2.Column2)),
		ExpectedStmt: `("table1".column1 = "table2".column1 AND "table1".column2 = "table2".column2 AND "table1".column2 = "table2".column2)`,
	},
	{
		Constructed:  Or(Table1.Column1.Eq(Table2.Column1), Table1.Column2.Eq(Table2.Column2), Table1.Column2.Eq(Table2.Column2)),
		ExpectedStmt: `("table1".column1 = "table2".column1 OR "table1".column2 = "table2".column2 OR "table1".column2 = "table2".column2)`,
	},
	{
		Constructed:  Select(Coalesce(Table1.Column1, Table1.Column2)).From(Table1),
		ExpectedStmt: `SELECT COALESCE("table1".column1, "table1".column2) FROM public.table1`,
	},
	{
		Constructed:  Select(Coalesce(Table1.Column1, Int64(0))).From(Table1),
		ExpectedStmt: `SELECT COALESCE("table1".column1, $1) FROM public.table1`,
	},
	{
		Constructed:  Select(Count()).From(Table1),
		ExpectedStmt: `SELECT COUNT(*) FROM public.table1`,
	},
	{
		Constructed:  Select(Count(Asterisk)).From(Table1),
		ExpectedStmt: `SELECT COUNT(*) FROM public.table1`,
	},
	{
		Constructed:  Select(Distinct(Table1.Column1)).From(Table1),
		ExpectedStmt: `SELECT DISTINCT("table1".column1) FROM public.table1`,
	},
	{
		Constructed:  Greatest(Int64(10), Int64(2), Int64(23)),
		ExpectedStmt: `GREATEST($1, $2, $3)`,
	},
	{
		Constructed:  Least(String("a"), String("b")),
		ExpectedStmt: `LEAST($1, $2)`,
	},
	{
		Constructed:  Ascii(String("abc")),
		ExpectedStmt: `ASCII($1)`,
	},
	{
		Constructed:  Ascii(Table1.Column1),
		ExpectedStmt: `ASCII("table1".column1)`,
	},
	{
		Constructed:  BTrim(String("    abc    ")),
		ExpectedStmt: `BTRIM($1)`,
	},
	{
		Constructed:  LTrim(Table1.Column1, String("xyz")),
		ExpectedStmt: `LTRIM("table1".column1, $1)`,
	},
	{
		Constructed:  RTrim(String("xyzxyzabcxyz"), Table1.Column1),
		ExpectedStmt: `RTRIM($1, "table1".column1)`,
	},
	{
		Constructed:  Chr(Table1.Column3),
		ExpectedStmt: `CHR("table1".column3)`,
	},
	{
		Constructed:  Concat(String("xyzxyzabcxyz"), Table1.Column3, Int64(3)),
		ExpectedStmt: `CONCAT($1, "table1".column3, $2)`,
	},
	{
		Constructed:  ConcatWs(String("x"), Table1.Column3, Int64(3), String("four")),
		ExpectedStmt: `CONCAT_WS($1, "table1".column3, $2, $3)`,
	},
	{
		Constructed:  Format(String("Hello %s, %1$s"), Table1.Column3),
		ExpectedStmt: `FORMAT($1, "table1".column3)`,
	},
	{
		Constructed:  Format(String("no formatting to be done")),
		ExpectedStmt: `FORMAT($1)`,
	},
	{
		Constructed:  InitCap(String("initCap THIS SenTEnce")),
		ExpectedStmt: `INITCAP($1)`,
	},
	{
		Constructed:  Left(Table1.Column1, Int64(3)),
		ExpectedStmt: `LEFT("table1".column1, $1)`,
	},
	{
		Constructed:  Right(String("take my right chars"), Table1.Column3),
		ExpectedStmt: `RIGHT($1, "table1".column3)`,
	},
	{
		Constructed:  Length(Table1.Column1),
		ExpectedStmt: `LENGTH("table1".column1)`,
	},
	{
		Constructed:  Length(String("jose"), String("UTF8")),
		ExpectedStmt: `LENGTH($1, $2)`,
	},
	{
		Constructed:  LPad(String("hi"), Int64(7)),
		ExpectedStmt: `LPAD($1, $2)`,
	},
	{
		Constructed:  RPad(Table1.Column2, Table1.Column3, Table1.Column1),
		ExpectedStmt: `RPAD("table1".column2, "table1".column3, "table1".column1)`,
	},
	{
		Constructed:  Md5(Table1.Column2),
		ExpectedStmt: `MD5("table1".column2)`,
	},
	{
		Constructed:  PgClientEncoding(),
		ExpectedStmt: `PG_CLIENT_ENCODING()`,
	},
	{
		Constructed:  QuoteIdent(String("foo bar")),
		ExpectedStmt: `QUOTE_IDENT($1)`,
	},
	{
		Constructed:  QuoteLiteral(String("foo bar")),
		ExpectedStmt: `QUOTE_LITERAL($1)`,
	},
	{
		Constructed:  QuoteLiteral(Float64(42.5)),
		ExpectedStmt: `QUOTE_LITERAL($1)`,
	},
	{
		Constructed:  QuoteNullable(Table3.ID),
		ExpectedStmt: `QUOTE_NULLABLE("table3".id)`,
	},
	{
		Constructed:  Repeat(String("abc"), Table1.Column3),
		ExpectedStmt: `REPEAT($1, "table1".column3)`,
	},
	{
		Constructed:  Replace(Table1.Column1, String("ab"), String("CD")),
		ExpectedStmt: `REPLACE("table1".column1, $1, $2)`,
	},
	{
		Constructed:  Reverse(String("reversable")),
		ExpectedStmt: `REVERSE($1)`,
	},
	{
		Constructed:  SplitPart(String("abc~@~def~@~ghi"), String("~@~"), Int64(2)),
		ExpectedStmt: `SPLIT_PART($1, $2, $3)`,
	},
	{
		Constructed:  Strpos(Table1.Column1, String("ab")),
		ExpectedStmt: `STRPOS("table1".column1, $1)`,
	},
	{
		Constructed:  Substr(Table1.Column1, Int64(2), Int64(5)),
		ExpectedStmt: `SUBSTR("table1".column1, $1, $2)`,
	},
	{
		Constructed:  StartsWith(String("alphabet"), String("alph")),
		ExpectedStmt: `STARTS_WITH($1, $2)`,
	},
	{
		Constructed:  ToAscii(Table1.Column1),
		ExpectedStmt: `TO_ASCII("table1".column1)`,
	},
	{
		Constructed:  ToAscii(String("Karel"), String("WIN1250")),
		ExpectedStmt: `TO_ASCII($1, $2)`,
	},
	{
		Constructed:  ToHex(Table1.Column3),
		ExpectedStmt: `TO_HEX("table1".column3)`,
	},
	{
		Constructed:  Translate(String("12345"), String("143"), String("ax")),
		ExpectedStmt: `TRANSLATE($1, $2, $3)`,
	},
	{
		Constructed:  TryAdvisoryLock(Int64(43)),
		ExpectedStmt: `pg_try_advisory_lock($1)`,
	},
	{
		Constructed:  ReleaseAdvisoryLock(Int64(52)),
		ExpectedStmt: `pg_advisory_unlock($1)`,
	},
	{
		Constructed:  LessThan(Int64(10), Int64(2)),
		ExpectedStmt: "$1 < $2",
	},
	{
		Constructed:  GreaterThan(Int64(10), Int64(2)),
		ExpectedStmt: "$1 > $2",
	},
	{
		Constructed:  LessOrEqual(Int64(10), Int64(2)),
		ExpectedStmt: "$1 <= $2",
	},
	{
		Constructed:  GreaterOrEqual(Int64(10), Int64(2)),
		ExpectedStmt: "$1 >= $2",
	},
	{
		Constructed:  Equal(Int64(10), Int64(2)),
		ExpectedStmt: "$1 = $2",
	},
	{
		Constructed:  NotEqual(Int64(10), Int64(2)),
		ExpectedStmt: "$1 <> $2",
	},
	{
		Constructed:  Concat(String("xyzxyzabcxyz"), Int64(2)),
		ExpectedStmt: "$1 || $2",
	},
	{
		Constructed:  OctetLength(String("xyzxyzabcxyz")),
		ExpectedStmt: "octet_length($1)",
	},
	{
		Constructed:  Overlay(String("xyzxyzabcxyz"), String("xyz"), Int64(2), Int64(3)),
		ExpectedStmt: "overlay($1 placing $2 from $3 for $4)",
	},
	{
		Constructed:  Overlay(String("xyzxyzabcxyz"), String("xyz"), Int64(2)),
		ExpectedStmt: "overlay($1 placing $2 from $3)",
	},
}

func TestFunctions(t *testing.T) {
	runTestCases(t, functionTestCases)
}
