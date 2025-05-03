package main

import (
	"log"

	"github.com/andlabs/ui"
	"github.com/kalvinarts/steam-deck-usbip/server-ui/src/gui"
)

func main() {
	err := ui.Main(func() {
		window := &gui.MainWindow{}
		if err := window.Create(); err != nil {
			ui.MsgBoxError(nil, "Error", err.Error())
			return
		}
		window.Show()
	})
	if err != nil {
		log.Fatal(err)
	}
}
