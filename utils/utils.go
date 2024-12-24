package utils

import (
	"encoding/json"
	"log"
)

// PrettyPrint returns a pretty-printed JSON string of the given data
func PrettyPrint(data interface{}) string {
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Printf("Failed to generate pretty JSON: %v", err)
		return ""
	}
	return string(prettyJSON)

}
