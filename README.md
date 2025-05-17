# Malwarebytes VPN for Linux

## IMPORTANT

The tool is in experimental mode, so it is important to know which parts of the Linux system are affected:

- Configuration files in user's config directory:
  - Session info: `~/.config/mbvpn/config.yml`
  - Machine ID: `~/.config/mbvpn/machine-id`
  - WireGuard configurations: `~/.config/mbvpn/servers/*.conf`
- `logout` command removes the configuration files but keeps WireGuard interfaces in the system

## Installation

### Option 1: Recommended Installation

The recommended installation method is to clone the repository and run the installation script:

1. Clone the repository:
   ```bash
   git clone https://github.com/Malwarebytes/mbvpn.git
   cd mbvpn
   ```

2. Run the installation script:
   ```bash
   ./install.sh
   ```

The script will:
1. Check system requirements
2. Install WireGuard dependencies if needed
3. Build MBVPN from source
4. Install the binary to `/usr/local/bin/mbvpn`

> Note: The script requires sudo privileges to install dependencies and place the binary in system directories.

### Option 2: From Releases

1. Download and unpack the latest release for your architecture from the "Releases" page.
2. Update the `$PATH` variable to have access to the unpacked binary.

#### Prerequisites for Manual Installation

WireGuard tools are required for MBVPN to function properly. If installing manually, you can install them with your package manager:

- For Ubuntu/Debian: `sudo apt install wireguard wireguard-tools`
- For Fedora: `sudo dnf install wireguard-tools`
- For CentOS/RHEL: `sudo yum install wireguard-tools`
- For Arch Linux: `sudo pacman -S wireguard-tools`
- For openSUSE: `sudo zypper install wireguard-tools`

1. Download and unpack the latest release for your architecture from the "Releases" page.
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

## Testing

### Integration and E2E Tests

```bash
# Run integration tests
go test -tags=integration ./test/integration/...

# Run specific integration test
go test -tags=integration ./test/integration/ -run TestRegisterDevice

# Run e2e tests
go test -tags=integration ./test/e2e/...
```

For security reasons, the tests require a license key to be provided through an environment variable. Before running tests, you must set:

```bash
export MBVPN_TEST_LICENSE_KEY="your-license-key-here"
```

Tests will fail if this environment variable is not set.

### Installation Script Testing

The installation script tests verify proper functioning on multiple Linux distributions using Docker containers:

```bash
# Test on all supported Linux distributions
make test-install

# Test on a specific distribution
make test-install-ubuntu
make test-install-fedora
make test-install-centos
make test-install-arch
make test-install-opensuse
```

Each test:
1. Builds a Docker container for the target distribution
2. Tests the installation script inside the container
3. Verifies WireGuard tools installation
4. Confirms proper binary installation

> Note: Running installation tests requires Docker to be installed on your system.

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

The `connect` (alias `c`), `disconnect` (alias `d`), and `status` commands require root privileges to manage WireGuard interfaces. To reduce the frequency of sudo password prompts, you can configure sudo to allow passwordless execution of these specific commands:

1. Create a dedicated sudoers file:

```bash
sudo visudo -f /etc/sudoers.d/mbvpn
```

2. Add the following content (replace `username` with your actual username and adjust the path if necessary):

```
# Allow user to run specific mbvpn commands without password
username ALL=(ALL) NOPASSWD: /usr/local/bin/mbvpn connect, /usr/local/bin/mbvpn c, /usr/local/bin/mbvpn disconnect, /usr/local/bin/mbvpn d, /usr/local/bin/mbvpn status
```

3. Set the correct permissions:

```bash
sudo chmod 440 /etc/sudoers.d/mbvpn
```

4. Test your configuration:

```bash
sudo -l | grep mbvpn
```

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
make build-prod-debug

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