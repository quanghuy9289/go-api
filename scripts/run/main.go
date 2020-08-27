package main

import (
	"api_new/logger"
	// "io/ioutil"
	"os"
	"os/exec"
	// "path"
	// "path/filepath"
)

var buildPath = "./build/api"

func main() {
	// Prepare build command
	logger.Infof("[PROD MODE] Start")
	cmd := exec.Command("./build/api")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the command
	logger.Infof(cmd.Dir)
	err := cmd.Run()
	if err != nil {
		logger.Fatalf("Encounter error: %v", err)
	}
	// out, err := cmdGQLGenerate.Output()
	// if err != nil {
	// 	logger.Fatalf("Encounter error: %v", err)
	// }
	// logger.Infof("[COMPLETE] Output: %s", out)
}
