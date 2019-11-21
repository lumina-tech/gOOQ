package nullable

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/lib/pq"
)

type StringArray struct {
	StringArray pq.StringArray
	Valid       bool // Valid is true if StringArray is not NULL
}

// StringArrayFrom creates a new StringArray that will never be blank.
func StringArrayFrom(value []string) StringArray {
	stringArray := StringArray{
		Valid: value != nil,
	}
	if value != nil {
		stringArray.StringArray = make([]string, len(value))
		copy(stringArray.StringArray, value)
	}
	return stringArray
}

func (self StringArray) MarshalText() ([]byte, error) {
	if !self.Valid {
		return []byte{}, nil
	}
	bytes, err := json.Marshal(self.StringArray)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (self *StringArray) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		self.StringArray = nil
		self.Valid = false
		return nil
	}
	arr := make([]string, 0)
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	self.Valid = true
	self.StringArray = arr
	return nil
}

// Scan implements the Scanner interface.
func (self *StringArray) Scan(v interface{}) error {
	if v == nil {
		self.Valid = false
		return nil
	}
	err := self.StringArray.Scan(v)
	if err == nil {
		self.Valid = true
	}
	return err
}

// Value implements the driver Valuer interface.
func (self StringArray) Value() (driver.Value, error) {
	if !self.Valid {
		return nil, nil
	}
	return self.StringArray.Value()
}
