package main

import (
	"flag"

	"api_new/server"
)

func main() {
	// Implement command flags
	// -port custom flag to specify which port the service should run on
	portPtr := flag.Int("port", -1, "port number to listen on")
	// -envFile customer flag to specify which env file to load
	envfilePtr := flag.String("envfile", "", "optional environment file to load and overwrite the current ENV")
	flag.Parse() // Parse the flag from command

	// If no value, reset to nil
	if portPtr != nil && *portPtr < 0 {
		portPtr = nil
	}

	if envfilePtr != nil && len(*envfilePtr) == 0 {
		envfilePtr = nil
	}

	// Run the server
	server.Run(portPtr, envfilePtr)
}
