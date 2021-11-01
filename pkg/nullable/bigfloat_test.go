package nullable_test

import (
	"github.com/lumina-tech/gooq/pkg/nullable"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestBigFloatScan(t *testing.T) {
	b := new(big.Float)
	b.SetInt64(0)
	bigFloat := nullable.BigFloatFrom(*b)
	require.True(t, bigFloat.Valid)
}
