#!/bin/bash

# Source common functions
source "$(dirname "$0")/common.sh"

# Dump the database
function dump_database() {
  echo -e "${BLUE}[*]${NC} Dumping the database..."
  
  # Get current database info
  local config_result=$(redis_cmd CONFIG GET dir 2>/dev/null)
  local db_dir=$(echo "$config_result" | tail -n1 | tr -d '\r')
  
  local dbfile_result=$(redis_cmd CONFIG GET dbfilename 2>/dev/null)
  local db_filename=$(echo "$dbfile_result" | tail -n1 | tr -d '\r')
  
  echo -e "${YELLOW}Current database location:${NC} $db_dir/$db_filename"
  
  if redis_cmd SAVE >/dev/null 2>&1; then
    echo -e "${GREEN}[+]${NC} Database dumped successfully!"
    echo -e "${GREEN}[+]${NC} Location: $db_dir/$db_filename"
    echo -e "${YELLOW}Note:${NC} Use scp, wget, or file inclusion to retrieve the dump file"
  else
    echo -e "${RED}[-]${NC} Failed to dump database"
    return 1
  fi
}