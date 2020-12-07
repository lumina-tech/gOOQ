package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"github.com/google/uuid"
)

type UUIDArray struct {
	UUIDArray []uuid.UUID
	Valid       bool // Valid is true if UUIDArray is not NULL
}

// UUIDArrayFrom creates a new UUIDArray that will never be blank.
func UUIDArrayFrom(value []uuid.UUID) UUIDArray{
	uuidArray := UUIDArray{
		Valid: value != nil,
	}
	if value != nil {
		uuidArray.UUIDArray = make([]uuid.UUID, len(value))
		copy(uuidArray.UUIDArray, value)
	}
	return uuidArray
}

func (self UUIDArray) MarshalText() ([]byte, error) {
	if !self.Valid {
		return []byte{}, nil
	}
	bytes, err := json.Marshal(self.UUIDArray)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (self *UUIDArray) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		self.UUIDArray = nil
		self.Valid = false
		return nil
	}
	arr := make([]uuid.UUID, 0)
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	self.Valid = true
	self.UUIDArray = arr
	return nil
}

// Scan implements the Scanner interface.
func (self *UUIDArray) Scan(v interface{}) error {
	if v == nil {
		self.Valid = false
		return nil
	}
	//switch x := v.(type) {
	switch _ := v.(type) {
	case []byte:
		/*value, err := uuid.Parse(string(x))
		if err != nil {
			return err
		}
		self.Valid = value != uuid.Nil*/
		self.UUIDArray = nil//value
	}
	return nil
}

// Value implements the driver Valuer interface.
func (self UUIDArray) Value() (driver.Value, error) {
	if !self.Valid {
		return nil, nil
	}
	a := self.UUIDArray

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+3*n)
		b[0] = '{'

		b = appendArrayQuotedBytes(b, []byte(a[0].String()))
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendArrayQuotedBytes(b, []byte(a[i].String()))
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

/*func (a *StringArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "StringArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(StringArray, len(elems))
		for i, v := range elems {
			if b[i] = string(v); v == nil {
				return fmt.Errorf("pq: parsing array element index %d: cannot convert nil to string", i)
			}
		}
		*a = b
	}
	return nil
}*/


func appendArrayQuotedBytes(b, v []byte) []byte {
	b = append(b, '"')
	for {
		i := bytes.IndexAny(v, `"\`)
		if i < 0 {
			b = append(b, v...)
			break
		}
		if i > 0 {
			b = append(b, v[:i]...)
		}
		b = append(b, '\\', v[i])
		v = v[i+1:]
	}
	return append(b, '"')
}
