#!/bin/bash

# Cleanup function to handle script interruption
cleanup() {
    echo -e "\nDisconnecting USB devices..."
    PORT=$(usbip port | grep "Port" | cut -d' ' -f2 | sed -e 's/[^0-9]*//g')
    if [ ! -z "$PORT" ]; then
        usbip detach -p "$PORT"
        echo "USB device detached"
    fi
    exit 0
}

# Set up trap for termination signals
trap cleanup SIGINT SIGTERM SIGKILL

# Check if IP address is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <steam-deck-ip>"
    exit 1
fi

STEAM_DECK_IP=$1

# Ensure script is run as root
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root"
    exit 1
fi

# Load required kernel module
modprobe vhci-hcd

# List available devices on Steam Deck
echo "Listing available devices on Steam Deck..."
DEVICES=$(usbip list -r $STEAM_DECK_IP)

if [ $? -ne 0 ]; then
    echo "Failed to connect to Steam Deck at $STEAM_DECK_IP"
    echo "Make sure the Steam Deck is running the server script and is accessible on the network"
    exit 1
fi

# Find Steam Controller bus ID
STEAM_BUSID=$(echo "$DEVICES" | grep "28de:" | cut -d: -f1)

if [ -z "$STEAM_BUSID" ]; then
    echo "No Steam Controller found on Steam Deck"
    exit 1
fi

echo "Found Steam Controller with bus ID: $STEAM_BUSID"
echo "Attaching USB device..."

# Attach the USB device
usbip attach -r $STEAM_DECK_IP -b $STEAM_BUSID

if [ $? -eq 0 ]; then
    echo "Successfully attached Steam Controller"
    echo "Press Ctrl+C to detach and exit"
    
    # Wait loop
    while true; do
        sleep 1
    done
else 
    echo "Failed to attach USB device"
    exit 1
fi