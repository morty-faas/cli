package debug

import "encoding/json"

// JSON is a simple helper function that return a pretty-print JSON string
func JSON(v any) string {
	d, _ := json.MarshalIndent(v, "", "\t")
	return string(d)
}
