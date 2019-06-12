package loader

import (
	"testing"

	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestPostgresLoader(t *testing.T) {
	t.Skip("comment out flaky test")

	db, err := sqlx.Connect("postgres", "user=lumina password=lumens dbname=lumina sslmode=disable")
	require.NoError(t, err)

	loader := NewPostgresLoader()
	tables, err := loader.TableList(db, "public")
	require.NoError(t, err)

	for _, table := range tables {
		_, err := loader.ColumnList(db, "public", table.TableName)
		require.NoError(t, err)
	}
	_, err = loader.EnumList(db, "public")
	require.NoError(t, err)
}
