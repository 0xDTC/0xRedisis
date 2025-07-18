# Redis CTF Exploitation Tool

A modular Redis exploitation tool designed for CTF (Capture The Flag) scenarios and educational purposes.

## Structure

```
0xRedisis/
├── Redisis              # Main executable script
├── Redisis_backup       # Backup of original monolithic script
├── README.md           # This file
└── features/           # Feature modules directory
    ├── common.sh       # Common utilities and functions
    ├── reconnaissance.sh # Redis reconnaissance module
    ├── webshell.sh     # Web shell injection module
    ├── ssh_injection.sh # SSH key injection module
    ├── database_dump.sh # Database dump module
    └── cron_injection.sh # Cron job injection module
```

## Usage

```bash
./Redisis <host> <port> [password]
```

### Examples

```bash
# Connect to Redis without authentication
./Redisis 192.168.1.100 6379

# Connect to Redis with password
./Redisis target.ctf.com 6379 mypassword
```

## Features

### 1. Reconnaissance
- Gathers server information (version, OS, architecture)
- Checks memory usage and configuration
- Lists database keys and keyspace info
- Tests for dangerous commands availability

### 2. Web Shell Injection
- Custom command execution shells
- Interactive PHP shells with cmd parameter
- Reverse shell payloads
- Automatically tries common web directories

### 3. SSH Key Injection
- Uses existing SSH keys
- Generates new key pairs
- Accepts custom public keys
- Tries common SSH directories

### 4. Database Dump
- Saves current database to disk
- Shows dump location
- Provides retrieval instructions

### 5. Cron Job Injection
- Injects persistent cron jobs
- Tries common cron directories
- Executes payloads every minute

## Modular Architecture

Each feature is contained in its own module for easy maintenance:

- **common.sh**: Shared utilities, Redis connection handling, input validation
- **reconnaissance.sh**: Information gathering functions
- **webshell.sh**: Web shell injection logic
- **ssh_injection.sh**: SSH key injection functionality
- **database_dump.sh**: Database dumping operations
- **cron_injection.sh**: Cron job injection methods

## Adding New Features

1. Create a new `.sh` file in the `features/` directory
2. Source `common.sh` at the top of your module
3. Define your functions
4. Add the source line to the main `Redisis` script
5. Add menu option and case statement entry

## Requirements

- **redis-cli**: Redis command-line interface
- **ssh-keygen**: For generating SSH keys (optional)
- **Standard Unix tools**: bash, timeout, grep, etc.

Install on Debian/Ubuntu:
```bash
sudo apt update
sudo apt install redis-tools openssh-client
```

## Security Notice

This tool is designed for educational purposes and authorized penetration testing only. Use only on systems you own or have explicit permission to test.
