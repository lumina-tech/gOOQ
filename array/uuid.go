package array

import (
	"fmt"

	"github.com/google/uuid"
)

type UUIDSlice []uuid.UUID

func (s *UUIDSlice) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Scan source was not []bytes")
	}
	result := []uuid.UUID{}
	strarray := parseArray(string(bytes))
	for _, str := range strarray {
		value, err := uuid.Parse(str)
		if err != nil {
			return err
		}
		result = append(result, value)
	}
	(*s) = result
	return nil
}
