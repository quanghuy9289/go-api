package scripts

import (
	"api_new/logger"
	"fmt"
	// "io/ioutil"
	"os"
	"os/exec"
	// "path"
	// "path/filepath"
)

// RunService run a service with name and environment file
func RunService(mode string, name string, envfile string) {
	// Prepare build command
	logger.Infof("[%s] Start service %s", mode, name)
	cmd := exec.Command("bee", "run", "-runargs", fmt.Sprintf("-envfile=%s", envfile))
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
