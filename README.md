# mister-macropads

A bridge to connect macropads to MiSTer FPGA systems, enabling enhanced control and interaction with your MiSTer setup through dedicated macropads.


## Overview

This project provides daemon services that connect various macropads to MiSTer FPGA systems, allowing users to control their MiSTer experience through dedicated hardware interfaces.


## Features

- **Multiple Backend Support**: Extensible architecture allowing the inclusion of various macropad types
- **MiSTer Integration**: Native integration with MiSTer FPGA Linux operating system
- **Screen Rendering**: Display support for macropads with screens
- **Configuration Management**: INI-based configuration, similar to other MiSTer-related projects
- **Daemon Process**: Runs as a background service with proper process management, despite the lack of a proper init system on the MiSTer FPGA Linux operating system
- **Startup Integration**: Automatic startup support through MiSTer init scripts


## Supported Devices

### Elgato Stream Deck

> [!WARNING]
> This is not implemented yet!

Several [Elgato Stream Deck] devices are supported via the `streamdeck_on.sh` and `streamdeck_off.sh` scripts. For a complete list of supported devices, please check the documentation of the underlying [streamdeck library](https://rafaelmartins.com/p/streamdeck).

### octokeyz

All [octokeyz](https://rafaelmartins.com/p/octokeyz) device variants are supported via the `octokeyz_on.sh` and `octokeyz_off.sh` scripts.


## Installation

TODO


## Usage

### Starting the Service

Use the MiSTer Scripts menu to run the start script for your macropad:

- `octokeyz_on.sh` for [octokeyz](https://rafaelmartins.com/p/octokeyz)
- `streamdeck_on.sh` for [Elgato Stream Deck](https://www.elgato.com/ww/en/s/explore-stream-deck)

This will enable the init scripts, allowing the daemon to start automatically when you boot the system.

### Stopping the Service

Use the MiSTer Scripts menu to run the stop script for your macropad:

- `octokeyz_off.sh` for [octokeyz](https://rafaelmartins.com/p/octokeyz)
- `streamdeck_off.sh` for [Elgato Stream Deck](https://www.elgato.com/ww/en/s/explore-stream-deck)

This will disable the init scripts, and the daemon will not start automatically when you boot the system.

The script may not exist if the service is disabled.

### Configuration

When a service is started, it automatically creates a sample configuration file if one is not found. Edit the generated `streamdeck.ini` or `octokeyz.ini` file to customize your macropad settings. It is stored in the `/media/fat/` directory.


## License

This source code is governed by a GPL-2-style license that can be found in the [LICENSE](LICENSE) file.
