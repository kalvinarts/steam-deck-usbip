package gui

import (
	"fmt"

	"github.com/andlabs/ui"
	"github.com/kalvinarts/steam-deck-usbip/server-ui/src/utils"
)

type MainWindow struct {
	window          *ui.Window
	statusLabel     *ui.Label
	ipLabel         *ui.Label
	controllerLabel *ui.Label
	startButton     *ui.Button
	stopButton      *ui.Button
	isRunning       bool
	device          *utils.USBDevice
}

func (mw *MainWindow) Create() error {
	if !utils.CheckRoot() {
		return fmt.Errorf("This application must be run as root")
	}

	mw.window = ui.NewWindow("Controller USB/IP Server", 400, 300, false)
	mw.window.SetMargined(true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	// Status section
	mw.statusLabel = ui.NewLabel("Status: Stopped")
	vbox.Append(mw.statusLabel, false)

	// IP Address section
	mw.ipLabel = ui.NewLabel("IP Address: Checking...")
	vbox.Append(mw.ipLabel, false)

	// Controller status
	mw.controllerLabel = ui.NewLabel("Controller: Not detected")
	vbox.Append(mw.controllerLabel, false)

	// Buttons
	buttonBox := ui.NewHorizontalBox()
	buttonBox.SetPadded(true)

	mw.startButton = ui.NewButton("Start Server")
	mw.stopButton = ui.NewButton("Stop Server")
	mw.stopButton.Disable()

	buttonBox.Append(mw.startButton, false)
	buttonBox.Append(mw.stopButton, false)
	vbox.Append(buttonBox, false)

	mw.window.SetChild(vbox)
	mw.window.OnClosing(func(*ui.Window) bool {
		if mw.isRunning {
			mw.stopServer()
		}
		ui.Quit()
		return true
	})

	// Set up button handlers
	mw.startButton.OnClicked(func(*ui.Button) {
		mw.startServer()
	})
	mw.stopButton.OnClicked(func(*ui.Button) {
		mw.stopServer()
	})

	// Initial IP check
	go mw.updateIP()

	return nil
}

func (mw *MainWindow) Show() {
	mw.window.Show()
}

func (mw *MainWindow) updateIP() {
	ip, err := utils.GetLocalIP()
	if err != nil {
		ui.QueueMain(func() {
			mw.ipLabel.SetText("IP Address: Error detecting IP")
		})
		return
	}
	ui.QueueMain(func() {
		mw.ipLabel.SetText("IP Address: " + ip)
	})
}

func (mw *MainWindow) startServer() {
	if err := utils.LoadKernelModules(); err != nil {
		ui.MsgBoxError(mw.window, "Error", "Failed to load kernel modules: "+err.Error())
		return
	}

	device, err := utils.FindSteamController()
	if err != nil {
		ui.MsgBoxError(mw.window, "Error", "Failed to find Steam Controller: "+err.Error())
		return
	}
	mw.device = device

	if err := utils.BindDevice(device.BusID); err != nil {
		ui.MsgBoxError(mw.window, "Error", "Failed to bind device: "+err.Error())
		return
	}

	if err := utils.StartDaemon(); err != nil {
		ui.MsgBoxError(mw.window, "Error", "Failed to start USBIP daemon: "+err.Error())
		return
	}

	mw.isRunning = true
	mw.statusLabel.SetText("Status: Running")
	mw.controllerLabel.SetText("Controller: Connected (Bus ID: " + device.BusID + ")")
	mw.startButton.Disable()
	mw.stopButton.Enable()

	ip, _ := utils.GetLocalIP()
	ui.MsgBox(mw.window, "Server Started",
		"Server is running!\nIP Address: "+ip+"\nBus ID: "+device.BusID+"\n\n"+
			"On your client machine, run:\n"+
			"sudo usbip list -r "+ip+"\n"+
			"sudo usbip attach -r "+ip+" -b "+device.BusID)
}

func (mw *MainWindow) stopServer() {
	if mw.device != nil {
		if err := utils.UnbindDevice(mw.device.BusID); err != nil {
			ui.MsgBoxError(mw.window, "Error", "Failed to unbind device: "+err.Error())
		}
	}

	if err := utils.StopDaemon(); err != nil {
		ui.MsgBoxError(mw.window, "Error", "Failed to stop USBIP daemon: "+err.Error())
	}

	mw.isRunning = false
	mw.statusLabel.SetText("Status: Stopped")
	mw.controllerLabel.SetText("Controller: Not detected")
	mw.startButton.Enable()
	mw.stopButton.Disable()
}
