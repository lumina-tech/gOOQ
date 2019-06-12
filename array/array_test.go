package array_test

import (
	"testing"

	"github.com/lumina-tech/lumina/apps/server/pkg/gooq/array"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestScanUUIDSlice(t *testing.T) {
	uuid1, _ := uuid.Parse("74e71edd-a652-4395-a2bf-0680252e5a3a")
	uuid2, _ := uuid.Parse("e564aafb-504b-465d-950d-8f77598e281a")
	uuids := []uuid.UUID{uuid1, uuid2}
	testScanUUIDSlice(t, []uuid.UUID{}, "{}")
	testScanUUIDSlice(t, uuids,
		"{74e71edd-a652-4395-a2bf-0680252e5a3a,e564aafb-504b-465d-950d-8f77598e281a}")
	testScanUUIDSlice(t, uuids,
		"{ 74e71edd-a652-4395-a2bf-0680252e5a3a , e564aafb-504b-465d-950d-8f77598e281a }")
	testScanUUIDSlice(t, uuids,
		" { 74e71edd-a652-4395-a2bf-0680252e5a3a , e564aafb-504b-465d-950d-8f77598e281a } ")
}

func testScanUUIDSlice(
	t *testing.T, uuids []uuid.UUID, str string,
) {
	uuidSlice := array.UUIDSlice{}
	err := uuidSlice.Scan([]byte(str))
	require.NoError(t, err)
	require.Len(t, uuidSlice, len(uuids))
	for index := range uuids {
		require.Equal(t, uuids[index], uuidSlice[index])
	}
}
