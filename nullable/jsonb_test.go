package nullable_test

import (
	"testing"

	"github.com/lumina-tech/gooq/nullable"
	"github.com/stretchr/testify/require"
)

func TestJsonbScan(t *testing.T) {
	bytes := []byte("hello world")
	jsonb := nullable.JsonbFrom(nil)
	require.False(t, jsonb.Valid)
	require.Nil(t, jsonb.Jsonb)

	err := jsonb.Scan(bytes)
	require.NoError(t, err)
	require.True(t, jsonb.Valid)
	require.Equal(t, bytes, jsonb.Jsonb)

	bytes[0] = byte('H')
	require.NotEqual(t, bytes, jsonb.Jsonb)
}
