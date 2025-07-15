// Package main implements the main function
package main

import (
	"os"

	api "github.com/diogovalentte/homarr-iframes/src"
	"github.com/diogovalentte/homarr-iframes/src/config"
)

func init() {
	// You can set the path to use an .env file below.
	// It can be an absolute path or relative to this file (main.go)
	//TODO: remove Path on Commit
	filePath := ".env.test"

	if err := config.SetConfigs(filePath); err != nil {
		panic(err)
	}
}

func main() {
	router := api.SetupRouter()
	router.SetTrustedProxies(nil)

	router.Run(":" + os.Getenv("PORT"))
}
