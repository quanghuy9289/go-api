package main

import (
	"api_new/logger"
	"io/ioutil"
	"os"
	"os/exec"
	// "path"
	"path/filepath"
)

var modulesPath = "modules"
var graphqlPath = "gql"

func main() {
	files, err := ioutil.ReadDir(modulesPath)
	if err != nil {
		logger.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			// Prepare go generate command
			module := f.Name()
			logger.Infof("[FOUND MODULE] Running go generate for module: %s", module)
			cmdGQLGenerate := exec.Command("go", "generate", "./...")
			cmdGQLGenerate.Stdout = os.Stdout
			cmdGQLGenerate.Stderr = os.Stderr
			cmdGQLGenerate.Stdin = os.Stdin
			// cmdGQLGenerate := exec.Command("go", "run", "-v", "github.com/99designs/gqlgen")
			var errPath error
			cmdGQLGenerate.Dir, errPath = filepath.Abs(filepath.Join(modulesPath, module, graphqlPath))
			if errPath != nil {
				logger.Fatalf("Could not generate absolute path for module: %s", module)
			}

			// Run the command
			logger.Infof(cmdGQLGenerate.Dir)
			cmdGQLGenerate.Run()
			// out, err := cmdGQLGenerate.Output()
			// if err != nil {
			// 	logger.Fatalf("Encounter error: %v", err)
			// }
			// logger.Infof("[COMPLETE] Output: Success go generate for module %s %s", module, out)
			logger.Infof("[COMPLETE] go generate for module %s", module)
		}
	}
}
