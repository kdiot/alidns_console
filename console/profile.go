package console

import (
	"encoding/json"
	"io"
	"os"
)

type Profile struct {
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	DomainName      string `json:"DomainName"`
}

func (profile *Profile) Load(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(data), profile)
	if err != nil {
		return err
	}

	return nil
}
