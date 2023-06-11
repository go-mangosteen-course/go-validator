package config

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type MyNullString struct {
	String string
	Valid  bool
}

func (s MyNullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return []byte(`"` + s.String + `"`), nil
	}
	return []byte("null"), nil
}
func (s *MyNullString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		s.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &s.String); err != nil {
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	s.Valid = true
	return nil
}

// 从数据库读值
func (s *MyNullString) Scan(value interface{}) error {
	if value == nil {
		s.Valid = false
		return nil
	}
	s.String, s.Valid = value.(string)
	return nil
}

// 向数据库写值
func (s MyNullString) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}
