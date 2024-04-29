# wsproxy

TTY to WebSocket proxy.

Supported devices:
- OHAUS Courier 5000

## Usage

The serial device will be autodiscovered. If two devices is detected, an error will be reported. Device selection is currently not supported.

The websocket endpoint is available at `ws://localhost:23193/ws`. A small debug page is available at http://localhost:23193

### Windows service

The proxy can be registered as a Windows service to run in background and be started automatically.

> [!IMPORTANT]  
> Administrator rights are required to manage services.

```sh
wsproxy.exe service install # Install the service
wsproxy.exe service stop    # Stop the service
wsproxy.exe service start   # Start the service
wsproxy.exe service remove  # Remove the service
```

When a executable is running, the file is locked. You must stop the service before updating.

## Development Mode

### Debug mode

The debug mode allow sending custom value using the debug page.

```sh
./wsproxy server --debug
```

### Simulate mode

The simulate serial mode disable the serial connection, no device is required. The serial port will always be reported as opened.
The debug mode is automaticaly enabled for sending custom value

```sh
./wsproxy server --simulate-serial
```

### Docker image

For easier development environment setup, a Docker image is available.

```yaml
services:
  wsproxy:
    image: ghcr.io/spacefoot/wsproxy:latest
    command: server --debug --simulate-serial
    ports:
      - 23193:23193
    ## To use a real device, remove "--simulate-serial" from command. The device MUST be connected first.
    ## Unplug and plug it back may break the mount and require a recreation of the container.
    ## The dialout group inside the container and on the host may not match, resulting in permission errors.
    ## You can get the host group ID with "getent group dialout | cut -d: -f3" and replace "dialout" with it.
    # user: nonroot:dialout
    # devices:
    #   - /dev/ttyACM0
```
