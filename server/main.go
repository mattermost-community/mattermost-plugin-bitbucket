package main

import (
	// "fmt"
	"github.com/mattermost/mattermost-server/plugin"
)

func main() {
	// fmt.Println("jason")
	plugin.ClientMain(&Plugin{})
}
