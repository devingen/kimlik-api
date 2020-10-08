package util

import (
	"encoding/json"
)

func ConvertMapToStruct(objectMap map[string]interface{}, value interface{}) error {

	bytes, encodeErr := json.Marshal(objectMap)
	if encodeErr != nil {
		return encodeErr
	}

	decodeErr := json.Unmarshal(bytes, value)
	if decodeErr != nil {
		return decodeErr
	}
	return nil
}

//
//func ConvertStructToMap(input interface{}) (output map[string]interface{}, err *utils.Error) {
//
//	bytes, encodeErr := json.Marshal(input)
//	if encodeErr != nil {
//		err = &utils.Error{Code: http.StatusInternalServerError, Message: encodeErr.Error()}
//		return
//	}
//
//	decodeErr := json.Unmarshal(bytes, &output)
//	if decodeErr != nil {
//		err = &utils.Error{Code: http.StatusInternalServerError, Message: decodeErr.Error()}
//		return
//	}
//	return
//}
//
//func ConvertStructToMapIgnoreError(input interface{}) (output map[string]interface{}) {
//	output, _ = ConvertStructToMap(input)
//	return
//}
//
//func ConvertStructArrayToMapArray(input []interface{}) (output []map[string]interface{}, err *utils.Error) {
//
//	bytes, encodeErr := json.Marshal(input)
//	if encodeErr != nil {
//		err = &utils.Error{Code: http.StatusInternalServerError, Message: encodeErr.Error()}
//		return
//	}
//
//	decodeErr := json.Unmarshal(bytes, &output)
//	if decodeErr != nil {
//		err = &utils.Error{Code: http.StatusInternalServerError, Message: decodeErr.Error()}
//		return
//	}
//	return
//}
//
//func ConvertStructArrayToMapArrayIgnoreError(input []interface{}) (output []map[string]interface{}) {
//	output, _ = ConvertStructArrayToMapArray(input)
//	return
//}
