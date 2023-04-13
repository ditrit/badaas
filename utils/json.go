package utils

import "strings"

// build a json from a string slices
// turn []string {`"a",1`, `"b",1`} to `{"a":1,"b":2}`
func BuildJSONFromStrings(pairs []string) string {
	var b strings.Builder
	// starting the json string
	b.WriteString("{")
	b.WriteString(strings.Join(pairs, ","))
	//ending the json string
	b.WriteString("}")
	return b.String()
}

// build a json list from a string slices
// turn []string {`{"a",1}`, `{"b",1}`} to `["{a":1},{"b":2}]`
func BuildJSONListFromStrings(pairs []string) string {
	var b strings.Builder
	// starting the json list
	b.WriteString("[")
	b.WriteString(strings.Join(pairs, ","))
	//ending the json list
	b.WriteString("]")
	return b.String()
}
