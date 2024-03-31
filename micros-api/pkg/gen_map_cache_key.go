package pkg

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
)

// sortAny Helper function to recursively sort maps and arrays of maps
func sortAny(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		return sortMap(v)
	case []interface{}:
		// Iterate over each element in the array
		for i, item := range v {
			v[i] = sortAny(item) // Recursively sort each item
		}
		return v
	default:
		return data // Return the data unchanged if not a map or array
	}
}

// sortMap Recursively sorts maps, including handling arrays of maps
func sortMap(m map[string]interface{}) map[string]interface{} {
	// Convert map keys to a slice and sort them
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create a new map with sorted keys
	sortedMap := make(map[string]interface{})
	for _, k := range keys {
		sortedMap[k] = sortAny(m[k]) // Recursively sort
	}

	return sortedMap
}

// GenMapCacheKey Serialize and optionally hash the map for use as a cache key
func GenMapCacheKey(m map[string]interface{}) (string, error) {
	// Recursively sort the map, including nested maps and arrays of maps
	sortedMap := sortAny(m)
	//fmt.Print(sortedMap)
	// Serialize the sorted map to JSON
	jsonBytes, err := json.Marshal(sortedMap)
	if err != nil {
		return "", err
	}

	// Optional: Hash the JSON string to get a fixed-size cache key
	hash := sha256.Sum256([]byte(jsonBytes))
	cacheKey := fmt.Sprintf("%x", hash)

	return cacheKey, nil
}
