package gjson

import "encoding/json"

func MustPrettyString(obj any) string {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func MustString(obj any) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(b)
}
