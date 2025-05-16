# Malwarebytes VPN for Linux

## IMPORTANT

The tool is in experimental mode, so it is important to know which parts of the Linux system are affected:

- Configuration files in user's config directory:
  - Session info: `~/.config/mbvpn/config.yml`
  - Machine ID: `~/.config/mbvpn/machine-id`
  - WireGuard configurations: `~/.config/mbvpn/servers/*.conf`
- `logout` command removes the configuration files but keeps WireGuard interfaces in the system

## Prerequisites

Install WireGuard tools using your preferred package manager:
- `wg`
- `wg-quick`

## Installation

1. Download and unpack the latest release for your architecture from "Releases".
2. Update the `$PATH` variable to have access to the unpacked binary.

## Build Configuration

MBVPN can be built with different configurations using the following matrix:

| Environment | Build Type | Command               | Use Case                                          |
|-------------|------------|----------------------|--------------------------------------------------|
| Staging     | Debug      | `make build-st-debug`  | Development and testing against staging           |
| Staging     | Release    | `make build-st-release`| Pre-production testing against staging            |
| Production  | Debug      | `make build-prod-debug`| Troubleshooting against production               |
| Production  | Release    | `make build-prod-release`| Final production builds                          |

Default commands:
- `make build-st`: Creates a staging release build
- `make build-prod`: Creates a production release build

### Build Types

- **Debug builds**: Include additional debugging symbols and provide more detailed logs
- **Release builds**: Optimized builds with minimal debug information

You can check your build configuration with the command:
```
mbvpn version
```

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
`mbvpn connect <server>` (or the shorthand `mbvpn c <server>`)

Disconnect from VPN:
`mbvpn disconnect` (or the shorthand `mbvpn d`)

Check connection status:
`mbvpn status`

End your session:
`mbvpn logout`

## Debug Mode

You can run any command in debug mode to see more detailed logs:

`mbvpn --debug <command>`

## Troubleshooting

### Login, logout and license issues

The `mbvpn logout` command deactivates your device (makes the license seat free) and removes local data for the current session. If you still face issues, delete the following files:
- `~/.config/mbvpn/config.yml`
- `~/.config/mbvpn/machine-id`

### Cannot disconnect, lost internet access

Use `mbvpn disconnect` (without specifying a server). The tool will attempt to disconnect from all WireGuard connections.

If that doesn't help, manually disconnect using WireGuard directly:
1. Find your connection: `wg show`
2. Disconnect: `sudo wg-quick down <interface_name>`

> After connecting/disconnecting, give the tool a few seconds to establish or terminate the connection.