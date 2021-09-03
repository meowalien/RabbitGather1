package util

import (
	"encoding/json"
	"io"
	"os"
)

// ParseFileJsonConfig read the json config file into the given struct
func ParseFileJsonConfig(sst interface{}, filePath string) (err error) {
	configFile, err := os.Open(filePath)
	if err != nil {
		return
	}
	return ParseReaderJson(sst, configFile)
}

// ParseReaderJson read the json from io.ReadCloser into the given struct
func ParseReaderJson(st interface{}, reader io.ReadCloser) error {
	defer func(rawbody io.ReadCloser) {
		err := rawbody.Close()
		if err != nil {
			panic("cannot close the reader")
		}
	}(reader)
	body := json.NewDecoder(reader)
	body.DisallowUnknownFields()
	err := body.Decode(st)
	if err != nil {
		return err
	}
	return nil
}
