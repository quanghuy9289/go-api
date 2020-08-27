package main

import (
	"api_new/scripts"
)

var envfile = "config/.dev.env"

func main() {
	scripts.RunService("DEV MODE", "Development service", envfile)
}
