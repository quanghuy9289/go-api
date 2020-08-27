package main

import (
	"api_new/scripts"
)

var envfile = "config/.dev.service.2.env"

func main() {
	scripts.RunService("DEPLOY DEV MODE", "Service 2", envfile)
}
