package main

import "github.com/sofiukl/io-service/core"

func main() {
	app := core.App{}
	app.Initialize()
	app.Run(":" + app.Config.ServerPort)
}
