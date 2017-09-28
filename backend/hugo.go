package main

import (
	"os/exec"
)

func buildStaticSite() ([]byte, error) {

	command := exec.Command(config.HugoBin, "--config", config.HugoConfigFile)

	out, err := command.Output()
	if err != nil {
		Error.Println("Couldn't publish", string(out), err.Error())
		return nil, err
	}

	return out, err

}
