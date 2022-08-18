package util

import (
	"bytes"
	"encoding/json"
)

func PrettyJson(m interface{}) string {
	data, _ := json.MarshalIndent(m, "", "    ")
	// replace Code: xx to Code: xx//text for code
	// dataStr := string(data)
	dataStr := deleteNewLine(data)
	return dataStr
}

func deleteNewLine(data []byte) string {
	intArrayFlag := false
	lineMapFlag := false
	var result bytes.Buffer
	for index, b := range data {
		switch b {
		case '[':
			if nextChar(index, data) == JSONInt {
				intArrayFlag = true
			}
			result.WriteByte(b)
		case ']':
			if intArrayFlag {
				intArrayFlag = false
			}
			result.WriteByte(b)
		case '{':
			if mergeLineInMap(index, data) {
				lineMapFlag = true
			}
			result.WriteByte(b)
		case '}':
			lineMapFlag = false
			result.WriteByte(b)
		case '\n', ' ':
			if !intArrayFlag && !lineMapFlag {
				result.WriteByte(b)
			}
		default:
			result.WriteByte(b)
		}
	}

	return result.String()
}

func mergeLineInMap(index int, data []byte) bool {
	bracketsDepth := 0
	elementCount := 0
	for _, b := range data[index+1:] {
		switch b {
		case '{':
			return false
		case '}':
			return elementCount <= 8
		case '[':
			bracketsDepth++
		case ']':
			bracketsDepth--
		case ':':
			if bracketsDepth == 0 {
				elementCount++
			}
		}
	}
	return false
}

const (
	// JSONArray means array struct [] in json
	JSONArray = iota
	// JSONMap means map struct {} in json
	JSONMap
	// JSONString means string "" in json
	JSONString
	// JSONInt means int in json
	JSONInt
	// JSONEnd means nothing valuable found
	JSONEnd
)

func nextChar(index int, data []byte) int {
	for _, b := range data[index+1:] {
		switch b {
		case '{':
			return JSONMap
		case '[':
			return JSONArray
		case '"':
			return JSONString
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			return JSONInt
		}
	}
	return JSONEnd
}
