package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config holds configuration options stored in a yaml file
type Config struct {
	Port               string
	Repository         string
	Logfile            string
	CORSEnabled        bool
	CORSOrigin         string
	Database           string // file path for BoltDB file
	Static             string
	HugoConfigFile     string              `yaml:"hugo_config_file"`
	HugoBin            string              `yaml:"hugo_bin"`
	PrivateKeyPath     string              `yaml:"private_key_path"`
	PublicKeyPath      string              `yaml:"public_key_path"`
	FileCategories     map[string][]string `yaml:"file_categories"`
	MediaTypes         map[string]string   `yaml:"media_types"`
	TranslationEnabled bool                `yaml:"translation_enabled"`
	DefaultLanguage    string              `yaml:"default_language"`
	EnabledLanguages   []string            `yaml:"enabled_languages"`
	AllLanguages       map[string]struct {
		Name string `yaml:"name"`
		Flag string `yaml:"flag"`
	} `yaml:"all_languages"`
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

	// simply check the supplied slice for the presence of string
	contains := func(s []string, e string) bool {
		for _, a := range s {
			if a == e {
				return true
			}
		}
		return false
	}

	codes := make([]string, len(c.AllLanguages))

	// if translation is enabled, ensure that our language settings
	// are valid
	if c.TranslationEnabled {

		// check languages have been set up correctly

		// first get a list of valid codes
		for l := range c.AllLanguages {
			codes = append(codes, l)
		}

		// make sure a default is set
		if c.DefaultLanguage == "" {
			return *c, fmt.Errorf("Translation enabled but no default language specified")
		}

		// throw an error if there are no enabled languages
		if len(c.EnabledLanguages) == 0 {
			return *c, fmt.Errorf("Translation enabled but no languages enabled")
		}

		// only a warning if there's just one - things will (kind of) work 🙄
		if len(c.EnabledLanguages) == 1 {
			Warning.Println("Translation is turned on but only one language is enabled")
		}

		// make sure the default language code exists
		defaultFound := contains(codes, c.DefaultLanguage)
		if !defaultFound {
			return *c, fmt.Errorf("Default language '%s' not found", c.DefaultLanguage)
		}

		// and make sure each of the enabled language codes do too
		for _, el := range c.EnabledLanguages {
			if !contains(codes, el) {
				return *c, fmt.Errorf("Language code '%s' not found", el)
			}
		}
	}

	return *c, err

}
