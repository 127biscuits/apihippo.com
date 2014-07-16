package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type config struct {
	Port int

	NeededVotesToVerify int `json:"verify_threshold"`
}

var Config config

func Load(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &Config); err != nil {
		errMessage := fmt.Sprintf("Can't unmarshal: %s. Error: %v", path, err)
		return errors.New(errMessage)
	}

	return nil
}
