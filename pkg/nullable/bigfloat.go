package nullable

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math/big"
)

type BigFloat struct {
	BigFloat big.Float
	Valid    bool //Valid is true if BigFloat is not NULL
}

func BigFloatFrom(value big.Float) BigFloat {
	bigFloat := BigFloat{
		Valid: true,
	}
	bigFloat.BigFloat = value

	return bigFloat
}

// Scan implements the Scanner interface.
func (b *BigFloat) Scan(v interface{}) error {
	if v == nil {
		b.Valid = false
		return nil
	}
	var i sql.NullString
	if err := i.Scan(v); err != nil {
		return err
	}
	if _, ok := b.BigFloat.SetString(i.String); ok {
		return nil
	}
	return fmt.Errorf("Could not scan type %T into BigFloat", v)
}

// Value implements the driver Valuer interface.
func (b BigFloat) Value() (driver.Value, error) {
	return b.Value()
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank BigFloat when this BigFloat is null.
func (b BigFloat) MarshalText() ([]byte, error) {
	return b.BigFloat.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null BigFloat if the input is a null BigFloat.
func (b *BigFloat) UnmarshalText(text []byte) error {
	return b.BigFloat.UnmarshalText(text)
}
