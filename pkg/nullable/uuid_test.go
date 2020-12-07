package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/lumina-tech/gooq/pkg/nullable"
	"github.com/stretchr/testify/require"
)

var (
	nilJSON       = []byte(`"00000000-0000-0000-0000-000000000000"`)
	uuidJSON1     = []byte(`"4925be64-f5dc-4a49-a682-98134ca5286d"`)
	uuidJSON2     = []byte(`"4925be64-f5dc-4a49-a682-98134ca5286d"`)
	blankUUIDJSON = []byte(`""`)
	nullJSON      = []byte(`null`)
	invalidJSON   = []byte(`:)`)
	boolJSON      = []byte(`true`)
)

type stringInStruct struct {
	Test nullable.UUID `json:"test,omitempty"`
}

func TestUUIDFrom(t *testing.T) {
	uuid1, _ := uuid.Parse("4925be64-f5dc-4a49-a682-98134ca5286d")
	u := nullable.UUIDFrom(uuid1)
	require.True(t, u.Valid)
	require.Equal(t, uuid1, u.UUID)

	// test nil case
	u = nullable.UUIDFrom(uuid.Nil)
	require.False(t, u.Valid)
	require.Equal(t, uuid.Nil, u.UUID)
}

func TestUUIDFromPtr(t *testing.T) {
	uuid1, _ := uuid.Parse("4925be64-f5dc-4a49-a682-98134ca5286d")
	uptr := &uuid1
	u := nullable.UUIDFromPtr(uptr)
	require.True(t, u.Valid)
	require.Equal(t, uuid1, u.UUID)

	u = nullable.UUIDFromPtr(&uuid.Nil)
	require.False(t, u.Valid)

	u = nullable.UUIDFromPtr(nil)
	require.False(t, u.Valid)
}

func TestUnmarshalUUID(t *testing.T) {
	expected, _ := uuid.Parse("4925be64-f5dc-4a49-a682-98134ca5286d")

	var u nullable.UUID
	err := json.Unmarshal(uuidJSON1, &u)
	require.NoError(t, err)
	require.True(t, u.Valid)
	require.Equal(t, expected, u.UUID)

	var blank nullable.UUID
	err = json.Unmarshal(blankUUIDJSON, &blank)
	require.Error(t, err)
	require.False(t, blank.Valid)
	require.Equal(t, uuid.Nil, blank.UUID)

	var nilUUID nullable.UUID
	err = json.Unmarshal(nilJSON, &nilUUID)
	require.NoError(t, err)
	require.False(t, nilUUID.Valid)
	require.Equal(t, uuid.Nil, nilUUID.UUID)

	var null nullable.UUID
	err = json.Unmarshal(nullJSON, &null)
	require.NoError(t, err)
	require.False(t, null.Valid)
	require.Equal(t, uuid.Nil, null.UUID)

	var badType nullable.UUID
	err = json.Unmarshal(boolJSON, &badType)
	require.Error(t, err)
	require.False(t, badType.Valid)
	require.Equal(t, uuid.Nil, badType.UUID)

	var invalid nullable.UUID
	err = invalid.UnmarshalJSON(invalidJSON)
	require.Error(t, err)
	require.False(t, invalid.Valid)
	require.Equal(t, uuid.Nil, invalid.UUID)
}

func TestTextUnmarshalUUID(t *testing.T) {
	uuid1, _ := uuid.Parse("4925be64-f5dc-4a49-a682-98134ca5286d")

	var u nullable.UUID
	err := u.UnmarshalText([]byte(uuid1.String()))
	require.NoError(t, err)
	require.True(t, u.Valid)
	require.Equal(t, uuid1, u.UUID)

	var nilUUID nullable.UUID
	err = nilUUID.UnmarshalText([]byte("00000000-0000-0000-0000-000000000000"))
	require.NoError(t, err)
	require.False(t, nilUUID.Valid)
	require.Equal(t, uuid.Nil, nilUUID.UUID)

	var null nullable.UUID
	err = null.UnmarshalText([]byte(""))
	require.NoError(t, err)
	require.False(t, null.Valid)
	require.Equal(t, uuid.Nil, null.UUID)
}
