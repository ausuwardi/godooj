package godooj

import (
	"fmt"
)

// Many2oneFK struct holds data for a Many2one field returned from Odoo
type Many2oneFK struct {
	ID   int
	Name string
}

// Val2Float converts a field value to float64
func Val2Float(f interface{}) (float64, bool) {
	switch val := f.(type) {
	case int:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

// FloatField gets float64 value from a record
func FloatField(rec interface{}, name string) (float64, error) {
	m, ok := rec.(map[string]interface{})
	if !ok {
		return -1, fmt.Errorf("Not a record")
	}

	iv, ok := m[name]
	if !ok {
		return -1, fmt.Errorf("Field not found: %s", name)
	}

	val, ok := Val2Float(iv)
	if !ok {
		return -1, fmt.Errorf("Error reading float from %s", name)
	}

	return val, nil

}

// Val2Int converts a field value to int
func Val2Int(f interface{}) (int, bool) {
	switch val := f.(type) {
	case float64:
		return int(val), true
	case int:
		return val, true
	default:
		return 0, false
	}
}

// IntField gets int value from a record
func IntField(rec interface{}, name string) (int, error) {
	m, ok := rec.(map[string]interface{})
	if !ok {
		return -1, fmt.Errorf("Not a record")
	}

	iv, ok := m[name]
	if !ok {
		return -1, fmt.Errorf("Field not found: %s", name)
	}

	val, ok := Val2Int(iv)
	if !ok {
		return -1, fmt.Errorf("Error reading int from %s", name)
	}

	return val, nil
}

// Val2String converts a field value to string
func Val2String(f interface{}) (string, bool) {
	switch v := f.(type) {
	case string:
		return v, true
	default:
		return "", false
	}
}

// StringField gets int value from a record
func StringField(rec interface{}, name string) (string, error) {
	m, ok := rec.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Not a record")
	}

	iv, ok := m[name]
	if !ok {
		return "", fmt.Errorf("Field not found: %s", name)
	}

	val, ok := Val2String(iv)
	if !ok {
		return "", fmt.Errorf("Error reading string from %s", name)
	}

	return val, nil
}

// Val2Many2one converts field value to Many2oneFK
func Val2Many2one(f interface{}) *Many2oneFK {
	switch v := f.(type) {
	case []interface{}:
		if len(v) == 2 {
			id, idOk := Val2Int(v[0])
			name, nameOk := Val2String(v[1])
			if idOk && nameOk {
				return &Many2oneFK{
					ID:   id,
					Name: name,
				}
			}
		}
	}
	return nil
}

// Many2oneField gets Many2one struct from a record
func Many2oneField(rec interface{}, name string) (*Many2oneFK, error) {
	m, ok := rec.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Not a record")
	}

	iv, ok := m[name]
	if !ok {
		return nil, fmt.Errorf("Field not found: %s", name)
	}

	val := Val2Many2one(iv)
	if val == nil {
		return nil, fmt.Errorf("Error reading Many2oneFK from %s", name)
	}

	return val, nil
}

// Val2One2many converts field value to []int
func Val2One2many(f interface{}) []int {
	switch v := f.(type) {
	case []float64:
		res := []int{}
		for _, ff := range v {
			res = append(res, int(ff))
		}
		return res
	case []interface{}:
		res := []int{}
		for _, ff := range v {
			res = append(res, int(ff.(float64)))
		}
		return res
	}
	return nil
}

// One2manyField gets one2many []int from a record
func One2manyField(rec interface{}, name string) ([]int, error) {
	m, ok := rec.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Not a record")
	}

	iv, ok := m[name]
	if !ok {
		return nil, fmt.Errorf("Field not found: %s", name)
	}

	val := Val2One2many(iv)
	if val == nil {
		return nil, fmt.Errorf("Error reading One2many from %s", name)
	}

	return val, nil
}
