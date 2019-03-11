package main

import (
	// "fmt"
	"github.com/mattermost/mattermost-server/plugin"
)

func main() {

	// dont do this!!  it makes plugin not connect
	// fmt.Println("jason")
	plugin.ClientMain(&Plugin{})
}
