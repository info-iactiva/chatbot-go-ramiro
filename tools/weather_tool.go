package tools

import (
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type WeatherArgs struct {
	Location string `json:"location"`
	Unit     string `json:"unit"`
}

func GetWeatherTool() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "getCurrentWeather",
			Description: "Fetches current weather for a given location",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "City and state, e.g., Boston, MA",
					},
					"unit": map[string]interface{}{
						"type": "string",
						"enum": []string{"fahrenheit", "celsius"},
					},
				},
				"required": []string{"location"},
			},
		},
	}
}

func ExecuteWeatherTool(argsJson string) (string, error) {
	var args WeatherArgs
	if err := json.Unmarshal([]byte(argsJson), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}

	// Simula una respuesta.
	return fmt.Sprintf("Weather in %s is 72°F and sunny.", args.Location), nil
}
