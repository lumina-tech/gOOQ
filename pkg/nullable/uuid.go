package nullable

import (
	"database/sql/driver"
	"encoding/json"
	fmt "fmt"
	"reflect"

	"github.com/google/uuid"
)

// UUID is a nullable UUID. It supports SQL and JSON serialization.
// It will marshal to null if null. Blank UUID input will be considered null.
type UUID struct {
	uuid.UUID
	Valid bool // Valid is true if String is not NULL
}

// newUUID creates a new UUID
func newUUID(value uuid.UUID) UUID {
	return UUID{
		UUID:  value,
		Valid: value != uuid.Nil,
	}
}

// UUIDFrom creates a new UUID that will never be blank.
func UUIDFrom(v uuid.UUID) UUID {
	return newUUID(v)
}

// UUIDFromPtr creates a new UUID that be null if s is nil.
func UUIDFromPtr(v *uuid.UUID) UUID {
	if v == nil {
		return newUUID(uuid.Nil)
	}
	return newUUID(*v)
}

// Scan implements the Scanner interface.
func (u *UUID) Scan(v interface{}) error {
	if v == nil {
		u.Valid = false
		return nil
	}
	switch x := v.(type) {
	case []byte:
		value, err := uuid.Parse(string(x))
		if err != nil {
			return err
		}
		u.Valid = value != uuid.Nil
		u.UUID = value
	}
	return nil
}

// Value implements the driver Valuer interface.
func (u UUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.UUID.String(), nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this UUID is null.
func (u UUID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(u.UUID)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports UUID and null input. Blank UUID input does not produce a null UUID.
// It also supports unmarshalling a sql.UUID.
func (u *UUID) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		value, err := uuid.Parse(x)
		if err != nil {
			return err
		}
		u.Valid = value != uuid.Nil
		u.UUID = value
		return nil
	case nil:
		u.Valid = false
		u.UUID = uuid.Nil
		return nil
	}
	return fmt.Errorf("json: cannot unmarshal %v into Go value of type null.UUID", reflect.TypeOf(v).Name())
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank UUID when this UUID is null.
func (u UUID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(u.UUID.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null UUID if the input is a blank UUID.
func (u *UUID) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" {
		u.UUID = uuid.Nil
		u.Valid = false
		return nil
	}
	value, err := uuid.Parse(string(text))
	if err != nil {
		return err
	}
	u.Valid = value != uuid.Nil
	u.UUID = value
	return nil
}

// SetValue changes this UUID's value and also sets it to be non-null.
func (u *UUID) SetValue(value uuid.UUID) {
	u.UUID = value
	u.Valid = value != uuid.Nil
}

// Ptr returns a pointer to this UUID's value, or a nil pointer if this UUID is null.
func (u UUID) Ptr() *uuid.UUID {
	if !u.Valid {
		return nil
	}
	return &u.UUID
}
