package main

import (
	"api_new/scripts"
)

var envfile = "config/.dev.service.1.env"

func main() {
	scripts.RunService("DEPLOY DEV MODE", "Service 1", envfile)
}
