#!/bin/bash

# Ensure script is run as root
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root"
    exit 1
fi

# Get local IP address (prefer wlan0 or eth0)
LOCAL_IP=$(ip -4 addr show scope global | grep inet | awk '{print $2}' | cut -d/ -f1 | head -n 1)
if [ -z "$LOCAL_IP" ]; then
    echo "Could not determine local IP address"
    exit 1
fi

# Load required kernel modules
modprobe usbip_core
modprobe usbip_host

# Find Steam Controller USB device
STEAM_BUSID=$(usbip list -l | grep -B 1 "28de:" | grep "busid" | awk '{print $3}')

if [ -z "$STEAM_BUSID" ]; then
    echo "Steam Controller not found"
    exit 1
fi

echo "Steam Controller found successfully"

# Bind the Steam Controller to usbip
echo "Binding Steam Controller..."
usbip bind -b $STEAM_BUSID

# Start usbip daemon
if ! pgrep usbipd > /dev/null; then
    echo "Starting usbipd..."
    usbipd -D
fi

echo "Server is ready. Steam Deck IP address: $LOCAL_IP"
echo "On your client machine, run:"
echo "  sudo usbip list -r $LOCAL_IP"
echo "Then connect with:"
echo "  sudo usbip attach -r $LOCAL_IP -b $STEAM_BUSID"
echo "Press Ctrl+C to stop the server"

# Keep the script running and handle cleanup on exit
trap cleanup SIGINT SIGTERM SIGKILL
cleanup() {
    echo "Unbinding USB device..."
    usbip unbind -b $STEAM_BUSID
    echo "Stopping usbipd..."
    killall usbipd
    exit 0
}

# Wait indefinitely
while true; do
    sleep 1
done