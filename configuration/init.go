package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Hostname        string `json:"hostname"`
	Port            string `json:"port"`
	SharedDirectory string `json:"shared_directory"`
}

const SHARED_DIR = "./shared"

func Init() error {
	config := Configuration{
		Hostname:        "127.0.0.1",
		Port:            "4000",
		SharedDirectory: SHARED_DIR,
	}
	fmt.Println("Creating `config.json` file.")
	fmt.Println("Creating `./shared` directory.")
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	err := ioutil.WriteFile("config.json", configJSON, 0644)
	if err != nil {
		return err
	}
	if _, err := os.Stat(SHARED_DIR); os.IsNotExist(err) {
		err = os.MkdirAll(SHARED_DIR, 0755)
		if err != nil {
			fmt.Println("Failed to create shared directory!")
		}
	}
	return nil
}

func Load(path_to_config string) (Configuration, error) {
	raw, err := ioutil.ReadFile(path_to_config)
	var config Configuration
	if err != nil {
		return config, errors.New(err.Error())
	}
	json.Unmarshal(raw, &config)

	return config, nil
}
