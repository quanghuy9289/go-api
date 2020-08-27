package main

import (
	"api_new/logger"
	// "io/ioutil"
	"os"
	"os/exec"
	// "path"
	"path/filepath"
)

var buildPath = "build"
var mainFile = "main.go"
var app = "api"

func main() {
	outputPath := filepath.Join(buildPath, app)

	// Prepare build command
	logger.Infof("[BEGIN] Build app: %s", app)
	cmd := exec.Command("go", "build", "-o", outputPath, mainFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the command
	err := cmd.Run()
	if err != nil {
		logger.Fatalf("Encounter error: %v", err)
	}
	logger.Infof("[COMPLETE] Built app %s to %s", app, outputPath)
}
