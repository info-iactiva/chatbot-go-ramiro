package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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


// Get pid of the current process
func GetPID() string {
	return fmt.Sprintf("%d", os.Getpid())
}