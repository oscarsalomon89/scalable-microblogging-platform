package main

import "github.com/oscarsalomon89/go-hexagonal/cmd/api/modules"

func main() {
	app := modules.NewApp()
	app.Run()
}
