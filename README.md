# Steam Deck usbip

Scripts to allow you to use your Steam Deck as a controller on Linux by sharing the controller over the network using `usbip`

## Prerequisites

On the Steam Deck:
1. Go to Desktop Mode

2. Open a terminal

3. Set a password for the `deck` user (if not already set):
   ```bash
   passwd deck
   ```

4. Enable writing to the Steam Deck's filesystem:
   ```bash
   sudo steamos-readonly disable
   ```

5. Install `usbip` tools and `git`:
   ```bash
   sudo pacman -S usbip git
   ```

On the client Linux machine:
- Install USB/IP tools (varies by distribution):
  - Arch Linux: `sudo pacman -S usbip`
  - Fedora: `sudo dnf install usbip-utils`
  - Debian/Ubuntu: `sudo apt install usbip`

## Installation (both machines)

1. Clone this repository to your desired location:
   ```bash
   git clone https://github.com/kalvinarts/steam-deck-usbip.git
   ```

## Usage

### On the Steam Deck (Server)

- See [Adding to Steam as a Non-Steam Game](#adding-to-steam-as-a-non-steam-game) for instructions on how to run this script from Steam's Big Picture mode. (recommended method)
- Alternatively, run the server script directly in a terminal (is recommended to have a keyboard connected if you use this method):
   ```bash
   sudo ./server.sh
   ```

To stop the server, you can either:

- Press `Ctrl+C` in the terminal (you will need a keyboard connected)
- Hard reset the Steam Deck using the power button (for now this is the only way to stop the server if you don't have a keyboard connected)

See the [roadmap](#roadmap) for future improvements.

### On the Client Machine

Run the client script passing your Steam Deck's IP address as the first argument:

```bash
sudo ./client.sh <steam-deck-ip>
```

The script will automatically connect to your Steam Deck's controller.

To manually manage the connection, you can use these commands instead of the client script:

1. List available devices:
   ```bash
   sudo usbip list -r <steam-deck-ip>
   ```

2. Attach the Steam Controller:
   ```bash
   sudo usbip attach -r <steam-deck-ip> -b <busid>
   ```

3. To detach when done:
   ```bash
   sudo usbip detach -p <port>
   ```

The Steam Controller should now be available on your client machine.

## Adding to Steam as a Non-Steam Game

To run the script directly from Steam's Big Picture mode:

1. In Desktop Mode, open Steam
2. Click "Add a Game" in the bottom left
3. Select "Add a Non-Steam Game"
4. Click "BROWSE"
5. Navigate to and select `/usr/bin/konsole`
6. Click "Add Selected Programs"
7. Right-click the newly added Konsole in your Steam library
8. Select "Properties"
9. Change the name to "USB/IP Controller" or similar
10. In the "LAUNCH OPTIONS" field, enter:
    ```
    --hold -e sudo /path/to/server.sh
    ```
   Replace `/path/to/server.sh` with the actual path to where you saved the script

Now you can launch the USB/IP server directly from Gaming Mode or Big Picture Mode. You'll need to enter your deck password when prompted.

## Troubleshooting

- Make sure both machines are on the same network
- Check that the required kernel modules are loaded: `usbip_core` and `usbip_host` on the Steam Deck, and `vhci-hcd` on the client (the scripts should handle this automatically)
- Ensure the firewall allows connections on port 3240

## Roadmap

- [ ] Add an installer script for easier setup on the Steam Deck
- [ ] Add an installer script for easier setup on the client machine
- [ ] Add a GUI for so you can start/stop the server without needing a keyboard
