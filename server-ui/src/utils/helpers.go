package utils

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

type USBDevice struct {
	BusID string
	Found bool
}

func CheckRoot() bool {
	return os.Geteuid() == 0
}

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no IP address found")
}

func LoadKernelModules() error {
	modules := []string{"usbip_core", "usbip_host"}
	for _, module := range modules {
		cmd := exec.Command("modprobe", module)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to load %s: %v", module, err)
		}
	}
	return nil
}

func FindSteamController() (*USBDevice, error) {
	cmd := exec.Command("usbip", "list", "-l")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list USB devices: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	device := &USBDevice{Found: false}

	for i, line := range lines {
		if strings.Contains(line, "28de:") && i > 0 {
			// Check previous line for busid
			busLine := lines[i-1]
			if strings.Contains(busLine, "busid") {
				fields := strings.Fields(busLine)
				if len(fields) >= 3 {
					device.BusID = fields[2]
					device.Found = true
					break
				}
			}
		}
	}

	if !device.Found {
		return device, fmt.Errorf("Steam Controller not found")
	}
	return device, nil
}

func BindDevice(busID string) error {
	cmd := exec.Command("usbip", "bind", "-b", busID)
	return cmd.Run()
}

func UnbindDevice(busID string) error {
	cmd := exec.Command("usbip", "unbind", "-b", busID)
	return cmd.Run()
}

func StartDaemon() error {
	// Check if daemon is already running
	if err := exec.Command("pgrep", "usbipd").Run(); err == nil {
		return nil // Daemon is already running
	}

	cmd := exec.Command("usbipd", "-D")
	return cmd.Run()
}

func StopDaemon() error {
	cmd := exec.Command("killall", "usbipd")
	return cmd.Run()
}
