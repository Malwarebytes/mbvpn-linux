# Malwarebytes VPN for Linux

[![Version](https://img.shields.io/github/v/release/Malwarebytes/mbvpn?style=flat-square)](https://github.com/Malwarebytes/mbvpn/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Malwarebytes/mbvpn)](https://goreportcard.com/report/github.com/Malwarebytes/mbvpn)
[![License](https://img.shields.io/github/license/Malwarebytes/mbvpn?style=flat-square)](LICENSE.md)

## IMPORTANT

The tool is in experimental mode, so it is important to know which parts of the Linux system are affected:

- Configuration files in user's config directory:
  - Session info: `~/.config/mbvpn/config.yml`
  - Machine ID: `~/.config/mbvpn/machine-id`
  - WireGuard configurations: `~/.config/mbvpn/servers/*.conf`
- `logout` command removes the configuration files but keeps WireGuard interfaces in the system

## Installation

### Prerequisites

WireGuard tools are required for MBVPN to function properly. Install them using your distribution's package manager:

- **Ubuntu/Debian**: `sudo apt install wireguard wireguard-tools`
- **Fedora**: `sudo dnf install wireguard-tools`
- **CentOS/RHEL**: `sudo yum install wireguard-tools`
- **Arch Linux**: `sudo pacman -S wireguard-tools`
- **openSUSE**: `sudo zypper install wireguard-tools`

### Installation via Go Package Manager (Recommended)

Install MBVPN directly using Go's package manager:

```bash
go install github.com/Malwarebytes/mbvpn@latest
```

> **Note**: Ensure that your `$GOPATH/bin` (typically `~/go/bin`) is in your `$PATH` to run the `mbvpn` command.

### Alternative: Download from Releases

1. Download the latest release for your architecture from the [Releases](https://github.com/Malwarebytes/mbvpn/releases) page
2. Extract the binary and place it in a directory in your `$PATH` (e.g., `/usr/local/bin`)

## Build Configuration

MBVPN can be built with different configurations using the following matrix:

| Environment | Build Type | Command               | Use Case                                          |
|-------------|------------|----------------------|--------------------------------------------------|
| Staging     | Debug      | `make build-st-debug`  | Development and testing against staging           |
| Staging     | Release    | `make build-st-release`| Pre-production testing against staging            |
| Production  | Debug      | `make build-debug`| Troubleshooting against production               |
| Production  | Release    | `make build-release`| Final production builds                          |

Default commands:
- `make build-st`: Creates a staging release build
- `make build-release`: Creates a production release build

### Build Types

- **Debug builds**: Include additional debugging symbols and provide more detailed logs
- **Release builds**: Optimized builds with minimal debug information

You can check your build configuration with the command:
```
mbvpn version
```

## Testing

### Integration and E2E Tests

```bash
# Run integration tests
go test -tags=integration ./test/integration/...

# Run specific integration test
go test -tags=integration ./test/integration/ -run TestRegisterDevice

# Run e2e tests
go test -tags=e2e ./test/e2e/...
```

For security reasons, the tests require a license key to be provided through an environment variable. Before running tests, you must set:

```bash
export MBVPN_TEST_LICENSE_KEY="your-license-key-here"
```

Tests will fail if this environment variable is not set.

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

### Reducing Sudo Password Prompts

The `connect` (alias `c`), `disconnect` (alias `d`), and `status` commands require root privileges to manage WireGuard interfaces because they internally call WireGuard tools using sudo. To reduce the frequency of sudo password prompts, you need to configure sudo to allow passwordless execution of the WireGuard commands:

1. Create a dedicated sudoers file:

```bash
sudo visudo -f /etc/sudoers.d/wireguard
```

2. Add the following content (replace `username` with your actual username):

```
# Allow user to run WireGuard commands without password
username ALL=(ALL) NOPASSWD: /usr/bin/wg, /usr/bin/wg-quick
```

3. Set the correct permissions:

```bash
sudo chmod 440 /etc/sudoers.d/wireguard
```

4. Test your configuration by running one of the mbvpn commands that require WireGuard tools:

```bash
mbvpn status
```

You should no longer be prompted for a password when running `mbvpn status`, `mbvpn connect`, or `mbvpn disconnect` commands.

This allows you to run the specified commands without entering a password, while still requiring the `sudo` prefix.

## Debug Mode and Logging

MBVPN provides several options for controlling logging and debugging output:

### Debug and Trace Flags

MBVPN provides two levels of enhanced logging:

#### Debug Mode

Run any command in debug mode to see detailed logs:

```bash
mbvpn --debug <command>
```

The `--debug` flag enables:
- Detailed error messages with stack traces
- More verbose logging of operations
- Display of underlying API calls and responses
- WireGuard configuration details

#### Trace Mode

For maximum verbosity, use trace mode which includes full HTTP request/response details:

```bash
mbvpn --trace <command>
```

The `--trace` flag enables:
- All debug mode features
- Complete HTTP request and response dumps
- Headers and payload contents of API communications
- Detailed internal operation tracing
- Useful for diagnosing connectivity issues or inspecting API behavior

### Debug Builds

Debug builds (created with `make build-*-debug` commands) provide enhanced debugging capabilities:
- Additional diagnostic information
- More comprehensive logging of operations
- Better visibility into system interactions

For troubleshooting complex issues, use debug builds with the debug flag:
```bash
# Build a debug version first
make build-debug

# Run with debug flag
mbvpn --debug <command>
```

### Log Levels

Debug logs are handled by the logrus library with the following behavior:
- By default, only INFO level and above logs are shown (warnings, errors)
- When using the `--debug` flag, DEBUG level logs are enabled, showing detailed operation information
- When using the `--trace` flag, TRACE level logs are enabled, showing complete HTTP request/response details
- HTTP requests, API responses, and internal operations are logged at DEBUG level
- Network traces (full request/response details) are logged at TRACE level
- Error details are always logged but only displayed to users in debug or trace mode

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

---

This repository is provided as-is and isn't bound to Malwarebytes' SLA.
