package godooj

import (
	"encoding/json"
	"time"
)

// OFloat is Odoo Float
type OFloat float64

// UnmarshalJSON parses float value from JSON
func (f *OFloat) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	switch v := i.(type) {
	case float64:
		*f = OFloat(v)
	case bool:
		*f = 0.0
	}
	return nil
}

// OInt is Odoo Int
type OInt int

// UnmarshalJSON parses Int value from JSON
func (i *OInt) UnmarshalJSON(b []byte) error {
	var ie interface{}
	if err := json.Unmarshal(b, &ie); err != nil {
		return err
	}
	switch v := ie.(type) {
	case int:
		*i = OInt(v)
	case bool:
		*i = 0
	}
	return nil
}

// OMany2one is a struct that hold Many2one value from Odoo
type OMany2one struct {
	ID   int
	Name string
}

// UnmarshalJSON parses Many2one struct from JSON
func (m *OMany2one) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	switch i.(type) {
	case []interface{}:
		slice := i.([]interface{})
		*m = OMany2one{
			ID:   slice[0].(int),
			Name: slice[1].(string),
		}
	case bool:
		*m = OMany2one{}
	}
	return nil
}

// ODateTime is Odoo datetime
type ODateTime time.Time

// UnmarshalJSON parses ODateTime from JSON
func (t *ODateTime) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	switch v := i.(type) {
	case string:
		ptime, err := time.Parse(time.RFC3339, v)
		if err == nil {
			*t = ODateTime(ptime)
		} else {
			*t = ODateTime(time.Time{})
		}
	case bool:
		*t = ODateTime(time.Time{})
	}
	return nil
}

// OString os Odoo String
type OString string

// UnmarshalJSON parses OString from JSON
func (s *OString) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	switch v := i.(type) {
	case string:
		*s = OString(v)
	case bool:
		*s = ""
	}
	return nil
}
