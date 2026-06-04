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
- `logout` command removes the configuration files

## Installation

### Homebrew (Linux)

```bash
brew install malwarebytes/tap/mbvpn
```

To allow mbvpn to create WireGuard interfaces without running as root, grant it the necessary capabilities:

```bash
sudo setcap cap_net_admin,cap_net_raw+eip $(which mbvpn)
```

> Note: You may need to reapply these capabilities after updating the binary.

## Usage

> You can access help information using `mbvpn --help` or `mbvpn <command> --help` if you want to read details about specific commands.

First, authenticate yourself using the `login` command with MB-code:

`mbvpn login --code MB-XXXXXX`

### Managing Servers

View available servers and locations:

- List all countries: `mbvpn countries`
- List cities within countries: `mbvpn cities`
- List all servers: `mbvpn servers`

### Connection Management

Connect to a server:
`mbvpn connect <server>` 

Disconnect from VPN:
`mbvpn disconnect` 

Check connection status:
`mbvpn status`

End your session and free up the license seat:
`mbvpn logout`

## Troubleshooting

### Login, logout and license issues

The `mbvpn logout` command deactivates your device (makes the license seat free) and removes local data for the current session. If you still face issues, delete the following files:
- `~/.config/mbvpn/config.yml`
- `~/.config/mbvpn/machine-id`

### Cannot disconnect, lost internet access

Use `mbvpn disconnect` (without specifying a server). The tool will attempt to disconnect from all WireGuard connections.

---

This repository is provided as-is and isn't bound to Malwarebytes' SLA.
