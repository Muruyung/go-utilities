package converter

import "time"

// ConvertTypeToStdMysqlType convertion data type to standard type mysql
func ConvertTypeToStdMysqlType(dataMap map[string]interface{}) map[string]interface{} {
	for key := range dataMap {
		switch val := dataMap[key].(type) {
		case bool:
			if val {
				dataMap[key] = 1
			} else {
				dataMap[key] = 0
			}
		case time.Time:
			dataMap[key] = ConvertDateToString(val)
		}
	}
	return dataMap
}

// ConvertBooleanToInt convertion data type boolean to int
func ConvertBooleanToInt(data bool) int8 {
	if data {
		return 1
	}

	return 0
}
