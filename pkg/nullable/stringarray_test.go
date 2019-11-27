package nullable_test

import (
	"testing"

	"github.com/lib/pq"

	"github.com/stretchr/testify/require"

	"github.com/lumina-tech/gooq/pkg/nullable"
)

func TestStringArrayScan(t *testing.T) {
	stringArray := nullable.StringArrayFrom(nil)
	require.False(t, stringArray.Valid)
	require.Nil(t, stringArray.StringArray)

	bytes := []byte("{\"hello\",\"world\"}")
	err := stringArray.Scan(bytes)
	require.NoError(t, err)
	require.True(t, stringArray.Valid)
	require.Equal(t, pq.StringArray{"hello", "world"}, stringArray.StringArray)

	value, err := stringArray.Value()
	require.NoError(t, err)
	require.Equal(t, string(bytes), value)
}

func TestStringArrayMarshalText(t *testing.T) {
	stringArray := nullable.StringArrayFrom([]string{"test"})

	err := stringArray.UnmarshalText(nil)
	require.NoError(t, err)
	require.False(t, stringArray.Valid)
	require.Nil(t, stringArray.StringArray)

	text := []byte("[\"hello\",\"world\"]")
	err = stringArray.UnmarshalText(text)
	require.NoError(t, err)
	require.True(t, stringArray.Valid)
	require.Equal(t, pq.StringArray{"hello", "world"}, stringArray.StringArray)

	marshalled, err := stringArray.MarshalText()
	require.NoError(t, err)
	require.Equal(t, text, marshalled)
}
