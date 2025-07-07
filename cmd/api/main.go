package main

import "github.com/oscarsalomon89/scalable-microblogging-platform/cmd/api/modules"

func main() {
	app := modules.NewApp()
	app.Run()
}
