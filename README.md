# Malwarebytes VPN for Linux

[![Version](https://img.shields.io/github/v/release/malwarebytes/mbvpn-linux?style=flat-square)](https://github.com/malwarebytes/mbvpn-linux/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/malwarebytes/mbvpn-linux)](https://goreportcard.com/report/github.com/malwarebytes/mbvpn-linux)
[![License](https://img.shields.io/github/license/malwarebytes/mbvpn-linux?style=flat-square)](LICENSE)

## IMPORTANT

The client is in experimental mode, so it is important to know which parts of the Linux system are affected:

- Configuration files in user's config directory:
  - Session info: `~/.config/mbvpn/config.yml`
  - Machine ID: `~/.config/mbvpn/machine-id`
  - WireGuard configurations: `~/.config/mbvpn/servers/*.conf`
- `logout` command removes the configuration files but keeps WireGuard interfaces in the system

## Installation

### Prerequisites

WireGuard tools are required for MBVPN to function properly. Install them using your distribution's package manager. For example:

- **Ubuntu/Debian**: `sudo apt install wireguard wireguard-tools`
- **Fedora**: `sudo dnf install wireguard-tools`
- **Arch Linux**: `sudo pacman -S wireguard-tools`

### Installation via Go Package Manager (Recommended)

Install MBVPN directly using Go's package manager:

```bash
go install github.com/malwarebytes/mbvpn-linux/cmd/mbvpn@latest
```

### Alternative: Download from Releases

1. Download the latest release for your architecture from the [Releases](https://github.com/Malwarebytes/mbvpn-linux/releases) page
2. Extract the binary and place it in a directory in your `$PATH` (e.g., `/usr/local/bin`)

## Usage

> You can access help information using `mbvpn --help` or `mbvpn <command> --help` if you want to read details about specific commands.

First, authenticate yourself using the `login` command with license key or MB-code:

`mbvpn login --code MB-XXXXXX`

or

`mbvpn login --key XXXX-XXXX-XXXX-XXXX`

### Managing Servers

View available servers and locations:

- List all countries: `mbvpn countries`
- List cities within countries: `mbvpn cities`
- List all servers: `mbvpn servers`

### Connection Management

Connect to a server:
`sudo mbvpn connect <server>` (or the shorthand `sudo mbvpn c <server>`)

Disconnect from VPN:
`sudo mbvpn disconnect` (or the shorthand `sudo mbvpn d`)

Check connection status:
`sudo mbvpn status`

End your session:
`mbvpn logout`

## Troubleshooting

### Login, logout and license issues

The `mbvpn logout` command deactivates your device (makes the license seat free) and removes local data for the current session. If you still face issues, delete the following files:
- `~/.config/mbvpn/config.yml`
- `~/.config/mbvpn/machine-id`

### Cannot disconnect, lost internet access

Use `sudo mbvpn disconnect` (without specifying a server). The tool will attempt to disconnect from all WireGuard connections.

If that doesn't help, manually disconnect using WireGuard directly:
1. Find your connection: `wg show`
2. Disconnect: `sudo wg-quick down <interface_name>`

> After connecting/disconnecting, give the tool a few seconds to establish or terminate the connection.

---

This repository is provided as-is and isn't bound to Malwarebytes' SLA.
