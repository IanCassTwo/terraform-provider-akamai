package tools

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func ConvertToString(data interface{}) (res string) {
	switch v := data.(type) {
	case float32, float64:
		res = fmt.Sprintf("%g", v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		res = fmt.Sprintf("%d", v)
	case json.Number:
		res = v.String()
	case string:
		res = v
	case []byte:
		res = string(v)
	case bool:
		res = strconv.FormatBool(v)
	default:
		res = fmt.Sprintf("%v", data)
	}
	return
}
