package converter

import (
	"encoding/json"
)

// ConvertInterfaceToMap convertion data from interface/struct to map
func ConvertInterfaceToMap(data interface{}) (converted map[string]interface{}, err error) {
	dataJSON, err := ConvertInterfaceToJSON(data)
	if err != nil {
		return
	}

	converted, err = ConvertJSONToMap(dataJSON)
	return
}

// ConvertMapToInterface convertion data from map to interface/struct
func ConvertMapToInterface(data map[string]interface{}, converted interface{}) (interface{}, error) {
	dataJSON, err := ConvertInterfaceToJSON(data)
	if err != nil {
		return nil, err
	}

	return ConvertJSONToInterface(dataJSON, converted)
}

// ConvertInterfaceToJSON convertion data from interface to json
func ConvertInterfaceToJSON(data interface{}) (converted []byte, err error) {
	converted, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return
}
