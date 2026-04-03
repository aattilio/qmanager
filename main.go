package main

import (
	"qmanager/src/interface/gui"
)

func main() {
	userInterface := gui.NewQManagerUI()
	userInterface.Run()
}
