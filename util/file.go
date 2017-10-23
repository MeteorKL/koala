package util

import (
	"encoding/json"
	"io/ioutil"
)

// ReadJSONFile 读取放在文件中的json数据
func ReadJSONFile(filename string) (map[string]string, error) {
	var data = map[string]string{}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func WriteJSONFile(filename string, data map[string]string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0666)
}
