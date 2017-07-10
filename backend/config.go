package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config holds configuration options stored in a yaml file
type Config struct {
	Port           string
	Repository     string
	Logfile        string
	CORSEnabled    bool
	CORSOrigin     string
	Database       string // file path for BoltDB file
	Static         string
	HugoConfigFile string              `yaml:"hugo_config_file"`
	HugoBin        string              `yaml:"hugo_bin"`
	PrivateKeyPath string              `yaml:"private_key_path"`
	PublicKeyPath  string              `yaml:"public_key_path"`
	FileCategories map[string][]string `yaml:"file_categories"`
}

func loadConfig(path *string) (Config, error) {

	// get the current working directory
	wd, err := os.Getwd()

	// if the file isn't found or can't be accessed
	if err != nil {
		return config, err
	}

	// read the config file
	configFile, err := ioutil.ReadFile(filepath.Join(wd, *path))

	if err != nil {
		return config, err
	}

	// convert the YAML into a Config struct
	c := &Config{}
	yaml.Unmarshal(configFile, c)

	return *c, err

}
