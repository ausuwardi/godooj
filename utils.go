package godooj

type Many2oneFK struct {
	Id   int
	Name string
}

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

func Val2String(f interface{}) (string, bool) {
	switch v := f.(type) {
	case string:
		return v, true
	default:
		return "", false
	}
}

func Val2Many2one(f interface{}) *Many2oneFK {
	switch v := f.(type) {
	case []interface{}:
		if len(v) == 2 {
			id, idOk := Val2Int(v[0])
			name, nameOk := Val2String(v[1])
			if idOk && nameOk {
				return &Many2oneFK{
					Id:   id,
					Name: name,
				}
			}
		}
	}
	return nil
}

func Val2One2many(f interface{}) []int {
	switch v := f.(type) {
	case []float64:
		res := []int{}
		for _, ff := range v {
			res = append(res, int(ff))
		}
		return res
	}
	return nil
}
