package impl

// Deep clone a map to avoid modifying the original object
func deepClone(obj map[string]interface{}) map[string]interface{} {
	clone := make(map[string]interface{})
	for key, value := range obj {
		clone[key] = deepCloneValue(value)
	}
	return clone
}

func deepCloneValue(value interface{}) interface{} {
	if m, ok := value.(map[string]interface{}); ok {
		return deepClone(m)
	}
	if s, ok := value.([]interface{}); ok {
		clonedSlice := make([]interface{}, len(s))
		for i, v := range s {
			clonedSlice[i] = deepCloneValue(v)
		}
		return clonedSlice
	}
	return value
}
