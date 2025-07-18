#!/bin/bash

# Source common functions
source "$(dirname "$0")/common.sh"

# Reconnaissance function
function reconnaissance() {
  echo -e "${BLUE}[*]${NC} Performing Redis reconnaissance..."
  
  echo -e "\n${YELLOW}=== Server Information ===${NC}"
  redis_cmd INFO server | grep -E "(redis_version|redis_git|redis_mode|os|arch_bits|process_id|tcp_port)"
  
  echo -e "\n${YELLOW}=== Memory Usage ===${NC}"
  redis_cmd INFO memory | grep -E "(used_memory_human|used_memory_peak_human|maxmemory_human)"
  
  echo -e "\n${YELLOW}=== Configuration ===${NC}"
  redis_cmd CONFIG GET "*" | grep -E "(dir|dbfilename|requirepass|save|maxmemory|bind)"
  
  echo -e "\n${YELLOW}=== Database Info ===${NC}"
  redis_cmd INFO keyspace
  
  echo -e "\n${YELLOW}=== Sample Keys ===${NC}"
  local keys=$(redis_cmd KEYS "*" | head -10)
  if [ -n "$keys" ]; then
    echo "$keys"
  else
    echo "No keys found"
  fi
  
  echo -e "\n${YELLOW}=== Dangerous Commands Check ===${NC}"
  local dangerous_cmds=("FLUSHDB" "FLUSHALL" "CONFIG" "EVAL" "SCRIPT")
  for cmd in "${dangerous_cmds[@]}"; do
    local result=$(redis_cmd COMMAND INFO "$cmd" 2>/dev/null)
    if [ $? -eq 0 ]; then
      echo -e "${GREEN}[+]${NC} $cmd command is available"
    else
      echo -e "${RED}[-]${NC} $cmd command is disabled"
    fi
  done
}