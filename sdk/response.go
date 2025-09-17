package sdk

import "encoding/json"

func ParseResponse(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
