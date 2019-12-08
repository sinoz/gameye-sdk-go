package helper

import (
	"encoding/json"
	"io/ioutil"
)

// ReadFileAsJSON reads a file and converts it to a json structure
func ReadFileAsJSON(path string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
