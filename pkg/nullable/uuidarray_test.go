package nullable_test

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/lumina-tech/gooq/pkg/nullable"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUUIDArrayScan(t *testing.T) {
	uuidArray := nullable.UUIDArrayFrom(nil)
	require.False(t, uuidArray.Valid)
	require.Nil(t, uuidArray.UUIDArray)

	bytes := []byte(fmt.Sprintf("{\"%s\",\"%s\"}", uuidJSON1, uuidJSON2))
	err := uuidArray.Scan(bytes)
	require.NoError(t, err)
	require.True(t, uuidArray.Valid)
	require.Equal(t, pq.StringArray{"hello", "world"}, uuidArray.UUIDArray)

	value, err := uuidArray.Value()
	require.NoError(t, err)
	require.Equal(t, string(bytes), value)
}

func TestUUIDArrayMarshalText(t *testing.T) {
	uuidArray := nullable.StringArrayFrom([]string{"test"})

	err := uuidArray.UnmarshalText(nil)
	require.NoError(t, err)
	require.False(t, uuidArray.Valid)
	require.Nil(t, uuidArray.StringArray)

	text := []byte("[\"hello\",\"world\"]")
	err = uuidArray.UnmarshalText(text)
	require.NoError(t, err)
	require.True(t, uuidArray.Valid)
	require.Equal(t, pq.StringArray{"hello", "world"}, uuidArray.StringArray)

	marshalled, err := uuidArray.MarshalText()
	require.NoError(t, err)
	require.Equal(t, text, marshalled)
}

