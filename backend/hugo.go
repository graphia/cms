package main

import (
	"fmt"
	"os"
	"os/exec"
)

func buildStaticSite() ([]byte, error) {

	command := exec.Command(config.HugoBin, "--config", config.HugoConfigFile)
	// FIXME change to hugo dir https://stackoverflow.com/questions/43135919/how-to-run-a-shell-command-in-a-specific-folder-with-golang

	// set the title
	env := os.Environ()
	env = append(env, fmt.Sprintf("HUGO_TITLE=%s", config.SiteTitle))
	command.Env = env

	out, err := command.CombinedOutput()
	if err != nil {
		Error.Println("Couldn't publish", string(out), err.Error())
		return nil, err
	}

	return out, err

}
