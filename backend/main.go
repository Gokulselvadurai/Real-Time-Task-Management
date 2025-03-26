package main

import (
	"my-rest-api/config"
	"my-rest-api/router"
)

func main() {
	config.Init()
	r := router.SetupRouter()
	r.Run(":8080")
}
