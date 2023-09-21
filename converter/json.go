package converter

import (
	"encoding/json"
)

// ConvertJSONToInterface convertion data from json to interface
func ConvertJSONToInterface(data []byte, converted interface{}) (interface{}, error) {
	err := json.Unmarshal(data, &converted)
	return converted, err
}

// ConvertJSONToMap convertion data from json to Map
func ConvertJSONToMap(data []byte) (converted map[string]interface{}, err error) {
	err = json.Unmarshal(data, &converted)
	return
}

// ConvertJSONToArrayMap convertion data from json to array of Map
func ConvertJSONToArrayMap(data []byte) (converted []map[string]interface{}, err error) {
	err = json.Unmarshal(data, &converted)
	return
}
