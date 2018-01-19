package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Configuration struct {
	Hostname          string `json:"hostname"`
	Port              string `json:"port"`
	EnableFileSharing bool   `json:"enable_file_sharing"`
	EnableProfile     bool   `json:"enable_profile"`
}

func Init() {
	config := Configuration{
		Hostname:          "localhost",
		Port:              "4000",
		EnableFileSharing: false,
		EnableProfile:     false,
	}
	configJSON, _ := json.Marshal(config)
	err := ioutil.WriteFile("config.json", configJSON, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Generated configuration file `config.json`.")
}

func Load(path_to_config string) (Configuration, error) {
	raw, err := ioutil.ReadFile(path_to_config)
	var config Configuration
	if err != nil {
		fmt.Println(err)
		return config, errors.New(err.Error())
	}
	json.Unmarshal(raw, &config)

	return config, nil
}
