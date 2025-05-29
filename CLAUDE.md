# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MBVPN is a command-line tool for Linux that interfaces with Malwarebytes VPN services. It uses WireGuard technology to establish secure VPN connections and communicates with Malwarebytes' backend services (Holocron) for authentication, server information, and connection management.

## Build Commands

**IMPORTANT**: All build, test, and Go commands must be run inside the Docker container using the docker-compose wrapper.

```bash
# Build for staging environment
docker-compose run --rm mbvpn make build-st

# Build for production environment  
docker-compose run --rm mbvpn make build-prod

# Build staging debug version
docker-compose run --rm mbvpn make build-st-debug

# Build production debug version
docker-compose run --rm mbvpn make build-prod-debug

# Install the binary (production)
docker-compose run --rm mbvpn make install

# Install staging version
docker-compose run --rm mbvpn make install-st

# Clean build artifacts
docker-compose run --rm mbvpn make clean

# Restart Docker environment
make docker-restart
```

## Testing Commands

**IMPORTANT**: All test and Go commands must be run inside the Docker container using the docker-compose wrapper.

```bash
# Run all tests
docker-compose run --rm mbvpn make test

# Run unit tests only
docker-compose run --rm mbvpn make test-unit

# Run integration tests
docker-compose run --rm mbvpn make test-integration

# Run specific integration test
docker-compose run --rm mbvpn go test -tags=integration ./test/integration/ -run TestRegisterDevice

# Run e2e tests
docker-compose run --rm mbvpn make test-e2e

# Run installation tests
make test-install

# Run installation tests for specific distros
make test-install-ubuntu
make test-install-fedora
make test-install-centos
make test-install-arch
make test-install-opensuse
```

## Core Architecture

The codebase follows a modular architecture with clear separation of concerns:

### Command Layer (cmd/)

- Uses Cobra library for CLI functionality
- Commands include: login, logout, connect, disconnect, status, servers, countries, cities

### Configuration (pkg/config/)

- Manages credentials and settings in `~/.config/mbvpn/config.yml`
- Handles machine identification with `ConfigFileMachineIdProvider` via `~/.config/mbvpn/machine-id`
- Exposes `ConfigProvider` interface for storage operations

### Remote API (pkg/remote/)

- `Holocron` interface communicates with Malwarebytes backend services
- GraphQL API calls for device registration, activation/deactivation, and VPN operations
- Manages server locations, network details, and key registration

### Server Management (pkg/servers/)

- `ServerStorage` interface handles server information storage and retrieval
- Stores server data in `~/.config/mbvpn/servers.json`
- Provides methods to find servers by name, city, or country code

### Session Management (pkg/session/)

- Handles user authentication via license key or MB code
- Controls session lifecycle (login/logout)
- Validates active sessions

### VPN Operations (pkg/vpn/)

- Manages WireGuard configuration generation
- Creates interface configurations in `~/.config/mbvpn/servers/`
- Handles connection and disconnection operations
- Provides status information

## Key Workflows

1. **Authentication Flow**:
   - Register device with backend → Get installation token → Activate device with license
   - Storage of session information in config.yml

2. **Connection Flow**:
   - Load server list → Select server → Register public key with backend → Generate WireGuard config → Establish connection

3. **Directory Structure**:
   - `~/.config/mbvpn/config.yml` - Session information
   - `~/.config/mbvpn/machine-id` - Unique device identifier
   - `~/.config/mbvpn/servers.json` - VPN server information
   - `~/.config/mbvpn/servers/*.conf` - WireGuard configurations

## Critical Dependencies

1. WireGuard tools (`wg` and `wg-quick`) - Required for VPN functionality
2. Root/sudo access - Required for network interface operations

## Environment Variables

- `MBVPN_HOLOCRON_URL_ST` - Staging environment API URL (set in .env file)
- `MBVPN_HOLOCRON_URL_PROD` - Production environment API URL (set in .env file)

## Commit Messages

**IMPORTANT**: Do NOT include any references to Claude Code, AI assistance, or co-authorship information in commit messages. Write clean, professional commit messages that focus solely on the changes made.

- Use conventional commit format when possible
- Focus on what was changed and why
- Do not mention Claude Code, AI tools, or add Co-Authored-By tags
- Keep messages concise and descriptive
