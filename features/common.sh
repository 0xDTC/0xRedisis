#!/bin/bash

# Colors for better output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Execute Redis commands with error handling
function redis_cmd() {
  local command="$@"
  local result
  local exit_code

  if [ -z "$PASSWORD" ]; then
    result=$(timeout 10 redis-cli -h "$HOST" -p "$PORT" $command 2>&1)
    exit_code=$?
  else
    result=$(timeout 10 redis-cli -h "$HOST" -p "$PORT" -a "$PASSWORD" $command 2>&1)
    exit_code=$?
  fi

  echo "$result"

  if [ $exit_code -eq 124 ]; then
    echo -e "${RED}Error:${NC} Command timed out after 10 seconds"
    return 1
  fi

  if echo "$result" | grep -q "Connection refused\|timeout\|Network is unreachable"; then
    echo -e "${RED}Error:${NC} Cannot connect to Redis server"
    return 1
  fi

  if echo "$result" | grep -q "NOAUTH"; then
    echo -e "${RED}Error:${NC} Authentication required but no password provided"
    return 1
  fi

  if echo "$result" | grep -q "ERR invalid password"; then
    echo -e "${RED}Error:${NC} Invalid password provided"
    return 1
  fi

  if echo "$result" | grep -q "ERR"; then
    echo -e "${YELLOW}Warning:${NC} Redis returned an error: $result"
    return 1
  fi

  if echo "$result" | grep -q "MISCONF"; then
    echo -e "${YELLOW}Warning:${NC} Redis cannot persist to disk due to configuration issues."
    echo "You may encounter issues with commands requiring persistence."
  fi

  return 0
}

# Connect to Redis using a purely read-only command
function connect_to_redis() {
  echo -e "${BLUE}[*]${NC} Testing Redis connection on ${YELLOW}$HOST:$PORT${NC}..."
  
  local info_result=$(redis_cmd INFO server 2>/dev/null)
  
  if [ $? -eq 0 ]; then
    local version=$(echo "$info_result" | grep "redis_version" | cut -d: -f2 | tr -d '\r')
    local os=$(echo "$info_result" | grep "os:" | cut -d: -f2 | tr -d '\r')
    echo -e "${GREEN}[+]${NC} Connection successful!"
    echo -e "${GREEN}[+]${NC} Redis version: ${YELLOW}$version${NC}"
    echo -e "${GREEN}[+]${NC} OS: ${YELLOW}$os${NC}"
    
    # Check if we can write (test with a harmless command)
    local test_result=$(redis_cmd PING 2>/dev/null)
    if echo "$test_result" | grep -q "PONG"; then
      echo -e "${GREEN}[+]${NC} Redis is responsive"
    fi
  else
    echo -e "${RED}[-]${NC} Failed to connect to Redis on $HOST:$PORT"
    exit 1
  fi
}

# Function to display the help menu
function show_help() {
  echo -e "\n${BLUE}Redis CTF Exploitation Tool${NC}"
  echo -e "${YELLOW}Usage:${NC} $0 <host> <port> [password]"
  echo -e "${YELLOW}Description:${NC} Redis exploitation tool for CTF scenarios"
  echo -e "${YELLOW}Examples:${NC}"
  echo -e "  $0 192.168.1.100 6379"
  echo -e "  $0 target.ctf.com 6379 mypassword"
  echo -e "If no password is provided, the script assumes no authentication is required.\n"
  exit 1
}

# Validate inputs
function validate_inputs() {
  # Check if at least host and port are provided
  if [ $# -lt 2 ]; then
    show_help
  fi

  HOST=$1
  PORT=$2
  PASSWORD=${3:-""}

  # Validate port number
  if ! [[ "$PORT" =~ ^[0-9]+$ ]] || [ "$PORT" -lt 1 ] || [ "$PORT" -gt 65535 ]; then
    echo -e "${RED}Error:${NC} Invalid port number. Must be between 1-65535."
    exit 1
  fi

  # Check if redis-cli is installed
  if ! command -v redis-cli &> /dev/null; then
    echo -e "${RED}Error:${NC} redis-cli is not installed or not in PATH."
    echo -e "${YELLOW}Install with:${NC} sudo apt-get install redis-tools"
    exit 1
  fi
}