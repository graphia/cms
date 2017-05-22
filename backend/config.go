package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config holds configuration options stored in a yaml file
type Config struct {
	Port        string
	Repository  string
	Logfile     string
	CORSEnabled bool
	CORSOrigin  string
	Database    string // file path for BoltDB file
	Static      string
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
