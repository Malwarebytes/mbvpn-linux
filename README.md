# Malwarebytes VPN for Linux

[![Version](https://img.shields.io/github/v/release/malwarebytes/mbvpn-linux?style=flat-square)](https://github.com/malwarebytes/mbvpn-linux/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/malwarebytes/mbvpn-linux)](https://goreportcard.com/report/github.com/malwarebytes/mbvpn-linux)
[![License](https://img.shields.io/github/license/malwarebytes/mbvpn-linux?style=flat-square)](LICENSE)

## IMPORTANT

The client is in experimental mode, so it is important to know which parts of the Linux system are affected:

- The root-owned daemon stores per-user credentials, keys, and connection state
  under `/var/lib/mbvpn`.
- The daemon socket is `/run/mbvpn/mbvpnd.sock` and is accessible only to members
  of the `mbvpn` Unix group.
- `logout` removes the calling user's daemon-owned session state.

## Installation

### systemd Linux

```bash
brew install malwarebytes/tap/mbvpn
```

The `mbvpn` client is unprivileged. Install and authorize the daemon once:

```bash
sudo groupadd --system mbvpn
sudo usermod -aG mbvpn "$USER"
sudo install -D -m 0644 "$(brew --prefix mbvpn)/libexec/mbvpnd.service" /usr/lib/systemd/system/mbvpnd.service
sudo systemctl daemon-reload
sudo systemctl enable --now mbvpnd.service
```

Start a new login session after changing group membership. Do not grant file
capabilities to `mbvpn`; the daemon holds the narrowly scoped network privilege.

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

Use `mbvpn disconnect` without specifying a server. It disconnects only the
calling user's recorded mbvpn connection and never enumerates unrelated
WireGuard interfaces.

---

This repository is provided as-is and isn't bound to Malwarebytes' SLA.
