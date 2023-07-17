package ddns

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kdiot/alidns-console/utility"
)

type DDNS struct {
	AccessKeyId     *string `json:"AccessKeyId"`
	AccessKeySecret *string `json:"AccessKeySecret"`
	DomainName      *string `json:"DomainName"`
	RR              *string `json:"RR"`
	Type            *string `json:"Type"`
	TTL             *int64  `json:"TTL"`
	Network         *string `json:"Network"`
}

func (d *DDNS) Check() error {
	if d.AccessKeyId == nil || *d.AccessKeyId == "" {
		return errors.New("AccessKeyId cannot be nil or empty")
	}
	if d.AccessKeySecret == nil || *d.AccessKeySecret == "" {
		return errors.New("AccessKeySecret cannot be nil or empty")
	}
	if d.DomainName == nil || *d.DomainName == "" {
		return errors.New("DomainName cannot be nil or empty")
	}
	if d.RR == nil || *d.RR == "" {
		return errors.New("RR cannot be nil or empty")
	}
	if d.Type == nil || (*d.Type != "A" && *d.Type != "AAAA") {
		return errors.New("domain name record type must be 'A' or 'AAAA'")
	}
	return nil
}

type Config struct {
	AccessKeyId     *string          `json:"AccessKeyId"`
	AccessKeySecret *string          `json:"AccessKeySecret"`
	DomainName      *string          `json:"DomainName"`
	LogFile         *string          `json:"LogFile"`
	LogLevel        utility.LogLevel `json:"LogLevel"`
	CheckInterval   time.Duration    `json:"CheckInterval"`
	RetryInterval   time.Duration    `json:"RetryInterval"`
	DomainList      []*DDNS          `json:"DomainList"`
}

func (conf *Config) Load(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open dynamic domain name configuration file, error: %s", err.Error())
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), &conf)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func LoadConfig(fileName string) (*Config, error) {
	conf := &Config{
		CheckInterval: 10,
		RetryInterval: 30,
	}
	if err := conf.Load(fileName); err != nil {
		return nil, err
	} else {
		return conf, nil
	}
}
