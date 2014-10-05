package settings

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	Server struct {
		Port        int
		MaxFileSize int `yaml:"max_file_size"`
	}

	NeededVotesToVerify int `yaml:"verify_threshold"`
	PageSize            int `yaml:"pagination_size"`

	CDN struct {
		Address string
	}

	Debug bool
}

var Config config

func Load(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(file, &Config); err != nil {
		return err
	}

	return nil
}
