package importer

import "fmt"

// =========================================================================
// headers array defines order and names of columns which this code is considered to get as input
var headers = [...]string{"ip_address", "country_code", "country", "city", "latitude", "longitude", "mystery_value"}

// map to easy get column index in csv line
var headersReverseMap map[string]int

func getRecordValue(columnName string) (int, error) {
	if headersReverseMap == nil {
		headersReverseMap = make(map[string]int)
		for index, headerName := range headers {
			headersReverseMap[headerName] = index
		}
	}

	index, ok := headersReverseMap[columnName]
	if !ok {
		return -1, fmt.Errorf("field with name %s not found", columnName)
	}

	return index, nil
}

func checkFileFormat(headLine []string) bool {
	if len(headLine) != len(headers) {
		return false
	}

	isValid := headLine[0] == headers[0] &&
		headLine[1] == headers[1] &&
		headLine[2] == headers[2] &&
		headLine[3] == headers[3] &&
		headLine[4] == headers[4] &&
		headLine[5] == headers[5] &&
		headLine[6] == headers[6]

	return isValid
}
