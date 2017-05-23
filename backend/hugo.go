package main

import (
	"bytes"
	"os/exec"
)

func buildStaticSite() ([]byte, error) {

	var stderr, stdout bytes.Buffer

	//Debug.Println("args", "--config", config.HugoConfigFile)

	//Debug.Println(fmt.Sprintf("Executing %s with %s", config.HugoBin, args))

	hugo := exec.Command(config.HugoBin, "--config", config.HugoConfigFile)

	hugo.Stdout = &stdout
	hugo.Stderr = &stderr

	//output, err := hugo.Output()

	err := hugo.Run()

	if err != nil {
		//Debug.Println("output:", output)
		Debug.Println("err:", stderr.String())
		return nil, err
	}

	//Debug.Println("hugo build output", output)

	return []byte("ok"), nil
}
