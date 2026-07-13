package transformer

func StripNullValues(input map[string]any) map[string]any {
	result := make(map[string]any)

	for key, value := range input {
		switch v := value.(type) {
		case map[string]any:
			// Recursively strip null values from nested maps
			nestedResult := StripNullValues(v)
			if len(nestedResult) > 0 {
				result[key] = nestedResult
			}
		case []any:
			// Recursively strip null values from slices
			var filteredSlice []any
			for _, item := range v {
				if itemMap, ok := item.(map[string]any); ok {
					// If the item is a map, strip null values from it
					strippedItem := StripNullValues(itemMap)
					if len(strippedItem) > 0 {
						filteredSlice = append(filteredSlice, strippedItem)
					}
				} else if item != nil {
					// If the item is not a map and not nil, keep it
					filteredSlice = append(filteredSlice, item)
				}
			}
			if len(filteredSlice) > 0 {
				result[key] = filteredSlice
			}
		default:
			// Keep non-nil values
			if v != nil {
				result[key] = v
			}
		}
	}
	return result
}
