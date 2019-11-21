package nullable

import (
	"database/sql/driver"
)

type Jsonb struct {
	Jsonb []byte
	Valid bool // Valid is true if Jsonb is not NULL
}

// JsonbFrom creates a new Jsonb that will never be blank.
func JsonbFrom(value []byte) Jsonb {
	jsonb := Jsonb{
		Valid: value != nil,
	}
	if value != nil {
		jsonb.Jsonb = make([]byte, len(value))
		copy(jsonb.Jsonb, value)
	}
	return jsonb
}

func (self Jsonb) MarshalText() ([]byte, error) {
	if !self.Valid {
		return []byte{}, nil
	}
	return []byte(self.Jsonb), nil
}

func (self *Jsonb) UnmarshalText(data []byte) error {
	str := string(data)
	if str == "" {
		self.Jsonb = nil
		self.Valid = false
		return nil
	}
	self.Valid = data != nil
	self.Jsonb = data
	return nil
}

// Scan implements the Scanner interface.
func (self *Jsonb) Scan(v interface{}) error {
	if v == nil {
		self.Valid = false
		return nil
	}
	switch x := v.(type) {
	case []byte:
		// must copy bytes to Jsonb.
		// cannot just assign jsonb to v because v might be reused
		self.Valid = true
		self.Jsonb = make([]byte, len(x))
		copy(self.Jsonb, x)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (self Jsonb) Value() (driver.Value, error) {
	if !self.Valid {
		return nil, nil
	}
	return self.Jsonb, nil
}
