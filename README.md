# Malwarebytes VPN for Linux

## IMPORTANT

The tool is in experimental mode, so it is important to know, what parts of the linux system are affected:

- Configuration file in user's config directory for storing session info: `~/.config/mbvpn/config.yml`
- `login` command creates Wireguard interfaces with `mb-` prefix in `/etc/wireguard/`
- `logout` command removes the configuration file only and keeps Wireguard interfaces in the file system

## Prerequisites

Install Wireguard tools using your prefered package manager:
- `wg`
- `wg-quick`

## Installation

1. Download and unpack the latest release for your architecture from "Releases".
2. Update the `$PATH` variable to have an access to the unpacked binary.

## Usage

> You can access help information using `mbvpn --help` or `mbvpn <command> --help` if you want to read details about the specific command.

First, you should authenticate yourself using `login` command with license key or MB-code:

`mbvpn login --code MB-XXXXXX`

or

`mbvpn login --key XXXX-XXXX-XXXX-XXXX`

Load and prepare VPN servers with `servers` command:

`mbvpn servers`

Use `mbvpn up <server>` and `mbvpn down <server>` to connect/disconnect to/from the server.

Use `mbvpn logout` to deactivate license and clean the current session.

Use `mbvpn status` to get information about current Wireguard connection.

## Troubleshooting

### Loging, logout and license issues

`mbvpn logout` command deactivates devices (makes license seat free) and removes local data for the current session. If you still face issues, delete `~/.config/mbvpn/config.yml` file.

### Cannot disconnect, lost the internet access

Try to use `mbvpn disconnect` without specifiying server. The tool will try to disconnect from the Wireguard connections, those server names start with `mb-`.

If it still doesn't help, use `wg-quick down <server>`.

> After connecting/disconnecting give the tool a few seconds to establish connection.