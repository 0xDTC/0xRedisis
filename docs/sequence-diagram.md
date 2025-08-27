# Sequence Diagrams - 0xRedisis Redis Exploitation Tool

## 1. Application Startup and Connection

```mermaid
sequenceDiagram
    participant U as User
    participant M as Main App
    participant P as Arg Parser
    participant C as Redis Client
    participant R as Redis Server
    participant Menu as Menu System
    
    U->>M: go run main.go <host> <port> [password]
    M->>P: Parse command arguments
    P->>P: Validate host, port, password
    
    alt Invalid Arguments
        P->>M: Return error
        M->>U: Display help and exit
    else Valid Arguments
        P->>M: Return config
        M->>C: NewRedisClient(config)
        C->>C: Initialize connection pool
        
        M->>C: Connect()
        C->>R: TCP connection attempt
        
        alt Connection Failed
            R->>C: Connection refused/timeout
            C->>M: Return connection error
            M->>U: Display error and exit
        else Connection Success
            R->>C: TCP connection established
            
            opt Authentication Required
                C->>R: AUTH <password>
                R->>C: +OK or -ERR
                
                alt Auth Failed
                    C->>M: Authentication error
                    M->>U: Display auth error and exit
                end
            end
            
            C->>M: Connection successful
            M->>Menu: ShowMainMenu(client, handlers)
            Menu->>U: Display main menu
        end
    end
```

## 2. Basic Reconnaissance Flow

```mermaid
sequenceDiagram
    participant U as User
    participant Menu as Menu System
    participant R as Recon Module
    participant C as Redis Client
    participant RS as Redis Server
    participant D as Display
    
    U->>Menu: Select option 1 (Reconnaissance)
    Menu->>R: reconnaissance.Reconnaissance(client)
    
    R->>U: Display reconnaissance menu
    U->>R: Select recon type
    
    alt Server Information
        R->>C: SendCommand("INFO")
        C->>RS: INFO command
        RS->>C: Server info response
        C->>R: Parsed server info
        R->>D: Format server information
        D->>U: Display formatted results
    
    else Configuration Analysis  
        R->>C: SendCommand("CONFIG", "GET", "*")
        C->>RS: CONFIG GET * command
        RS->>C: Configuration data
        C->>R: Config key-value pairs
        R->>R: Analyze dangerous configs
        R->>D: Format security analysis
        D->>U: Display config analysis
        
    else Memory Analysis
        R->>C: SendCommand("INFO", "memory")
        C->>RS: INFO memory command
        RS->>C: Memory statistics
        C->>R: Memory data
        R->>R: Calculate usage patterns
        R->>D: Format memory report
        D->>U: Display memory analysis
        
    else Keyspace Analysis
        R->>C: SendCommand("INFO", "keyspace")
        C->>RS: INFO keyspace command
        RS->>C: Database statistics
        C->>R: Keyspace data
        
        loop For each database
            R->>C: SendCommand("SELECT", db)
            C->>RS: SELECT database
            R->>C: SendCommand("DBSIZE")
            C->>RS: DBSIZE command  
            RS->>C: Key count
        end
        
        R->>D: Format keyspace report
        D->>U: Display database analysis
    end
    
    R->>U: Reconnaissance completed
    U->>Menu: Return to main menu
```

## 3. Web Shell Injection Sequence

```mermaid
sequenceDiagram
    participant U as User
    participant W as WebShell Module
    participant C as Redis Client
    participant RS as Redis Server
    participant L as Listener
    participant FS as File System
    
    U->>W: Select web shell injection
    W->>U: Display shell options
    U->>W: Choose shell type and directory
    
    W->>W: Generate PHP payload
    W->>W: Detect web-accessible paths
    
    loop For each potential directory
        W->>C: SendCommand("CONFIG", "SET", "dir", directory)
        C->>RS: CONFIG SET dir <path>
        
        alt Directory writable
            RS->>C: +OK
            W->>C: SendCommand("CONFIG", "SET", "dbfilename", "shell.php")
            C->>RS: CONFIG SET dbfilename shell.php
            RS->>C: +OK
            
            W->>C: SendCommand("SET", "webshell", payload)
            C->>RS: SET webshell <php_code>
            RS->>C: +OK
            
            W->>C: SendCommand("SAVE")
            C->>RS: SAVE command
            
            alt Save successful
                RS->>C: +OK
                RS->>FS: Write shell.php to directory
                W->>W: Verify shell accessibility
                
                opt Reverse Shell Requested
                    W->>L: StartListener(port)
                    L->>L: Bind to port and listen
                    W->>W: Generate reverse shell payload
                    W->>C: SendCommand("SET", "reverse", reverse_payload)
                    C->>RS: SET reverse <reverse_shell>
                    W->>C: SendCommand("SAVE")
                    C->>RS: SAVE command
                    RS->>FS: Write reverse shell
                    
                    Note over L: Waiting for connection...
                    FS->>L: Incoming reverse connection
                    L->>U: Shell session established
                end
                
                W->>U: Web shell deployed successfully
                break
            else Save failed
                RS->>C: -ERR
                W->>W: Try next directory
            end
        else Directory not writable
            RS->>C: -ERR
            W->>W: Try next directory
        end
    end
    
    W->>U: Return results
```

## 4. SSH Key Injection Flow

```mermaid
sequenceDiagram
    participant U as User
    participant S as SSH Module  
    participant C as Redis Client
    participant RS as Redis Server
    participant K as Key Generator
    participant FS as Target File System
    
    U->>S: Select SSH key injection
    S->>K: Generate RSA key pair
    K->>K: Create 2048-bit RSA keys
    K->>S: Return public/private keys
    
    S->>U: Display key information
    S->>S: Prepare authorized_keys content
    
    loop For each potential SSH path
        Note over S: Trying /root/.ssh/authorized_keys, /home/*/.ssh/authorized_keys
        
        S->>C: SendCommand("CONFIG", "SET", "dir", ssh_directory)
        C->>RS: CONFIG SET dir <ssh_path>
        
        alt Directory accessible
            RS->>C: +OK
            S->>C: SendCommand("CONFIG", "SET", "dbfilename", "authorized_keys")
            C->>RS: CONFIG SET dbfilename authorized_keys
            RS->>C: +OK
            
            opt Backup existing keys
                S->>C: SendCommand("SET", "backup", "")
                S->>C: SendCommand("SAVE")
                C->>RS: SAVE (create backup)
            end
            
            S->>C: SendCommand("SET", "sshkey", public_key)
            C->>RS: SET sshkey <public_key>
            RS->>C: +OK
            
            S->>C: SendCommand("SAVE")
            C->>RS: SAVE command
            
            alt Save successful
                RS->>C: +OK
                RS->>FS: Write authorized_keys
                
                S->>S: Set file permissions (chmod 600)
                S->>C: SendCommand("EVAL", chmod_lua_script)
                C->>RS: EVAL script for chmod
                RS->>FS: Change file permissions
                
                S->>S: Test SSH connection
                S->>FS: ssh -i private_key user@host
                
                alt SSH connection successful
                    FS->>S: SSH session established
                    S->>U: SSH key injection successful
                    break
                else SSH connection failed
                    S->>S: Try next path
                end
            else Save failed
                RS->>C: -ERR
                S->>S: Try next path
            end
        else Directory not accessible
            RS->>C: -ERR  
            S->>S: Try next path
        end
    end
    
    S->>U: Return injection results
```

## 5. Lua Script Exploitation Flow

```mermaid
sequenceDiagram
    participant U as User
    participant Lua as Lua Module
    participant C as Redis Client
    participant RS as Redis Server
    participant OS as Operating System
    participant L as Listener
    
    U->>Lua: Select Lua exploitation
    Lua->>U: Display exploitation options
    U->>Lua: Choose operation type
    
    alt System Command Execution
        U->>Lua: Enter command to execute
        Lua->>Lua: Wrap command in io.popen() script
        
        Lua->>C: SendCommand("EVAL", lua_script, "0")
        C->>RS: EVAL <script> 0
        RS->>OS: Execute system command
        OS->>RS: Command output
        RS->>C: Script result with output
        C->>Lua: Command execution result
        Lua->>U: Display command output
        
    else File Operations
        U->>Lua: Choose file operation (read/write/delete)
        U->>Lua: Enter file path and content
        
        Lua->>Lua: Generate file operation script
        Lua->>C: SendCommand("EVAL", file_script, "0")
        C->>RS: EVAL <file_script> 0
        RS->>OS: File system operation
        OS->>RS: Operation result
        RS->>C: File operation result
        C->>Lua: File operation status
        Lua->>U: Display operation result
        
    else Reverse Shell
        Lua->>L: Start listener on specified port
        L->>L: Bind and listen on port
        
        Lua->>Lua: Generate reverse shell Lua script
        Lua->>C: SendCommand("EVAL", reverse_script, "0")
        C->>RS: EVAL <reverse_script> 0
        RS->>OS: Execute reverse shell
        OS->>L: Connect to listener
        L->>Lua: Connection established
        Lua->>U: Reverse shell session ready
        
        loop Shell Session
            U->>L: Enter shell command
            L->>OS: Execute command
            OS->>L: Command output
            L->>U: Display output
        end
        
    else Persistent Backdoor
        Lua->>Lua: Create persistent Lua backdoor
        Lua->>C: SendCommand("EVAL", backdoor_script, "0")
        C->>RS: EVAL <backdoor_script> 0
        RS->>OS: Install backdoor mechanism
        OS->>RS: Installation result
        RS->>C: Backdoor installation status
        C->>Lua: Installation confirmation
        Lua->>U: Backdoor installed successfully
    end
    
    Lua->>U: Lua exploitation completed
```

## 6. Smart Data Exfiltration Sequence

```mermaid
sequenceDiagram
    participant U as User
    participant E as Exfiltration Module
    participant C as Redis Client
    participant RS as Redis Server
    participant P as Pattern Matcher
    participant Export as Export System
    
    U->>E: Select smart data exfiltration
    E->>U: Display exfiltration options
    U->>E: Choose scan type and patterns
    
    E->>E: Initialize pattern library
    E->>E: Load sensitivity patterns
    
    loop For each database
        E->>C: SendCommand("SELECT", database_num)
        C->>RS: SELECT <db>
        RS->>C: +OK
        
        E->>C: SendCommand("KEYS", "*")
        C->>RS: KEYS *
        RS->>C: List of all keys
        
        loop For each key
            E->>C: SendCommand("TYPE", key)
            C->>RS: TYPE <key>
            RS->>C: Key type (string/list/hash/etc)
            
            alt String type
                E->>C: SendCommand("GET", key)
                C->>RS: GET <key>
                RS->>C: Key value
                
            else Hash type
                E->>C: SendCommand("HGETALL", key)
                C->>RS: HGETALL <key>
                RS->>C: Hash field-value pairs
                
            else List type
                E->>C: SendCommand("LRANGE", key, "0", "-1")
                C->>RS: LRANGE <key> 0 -1
                RS->>C: List elements
                
            else Set type
                E->>C: SendCommand("SMEMBERS", key)
                C->>RS: SMEMBERS <key>
                RS->>C: Set members
            end
            
            C->>E: Raw key data
            E->>P: Analyze data for patterns
            
            P->>P: Check for passwords
            P->>P: Check for API keys
            P->>P: Check for PII data
            P->>P: Check for financial data
            P->>P: Check for network info
            
            alt Sensitive data found
                P->>E: Pattern matches with sensitivity scores
                E->>E: Classify and categorize findings
                E->>Export: Queue for export
            else No sensitive data
                P->>E: No matches found
            end
        end
    end
    
    E->>Export: Process all findings
    Export->>Export: Generate reports by category
    Export->>Export: Apply export format (JSON/CSV/TXT)
    
    alt Real-time monitoring requested
        E->>E: Start monitoring loop
        loop Monitoring active
            E->>C: SendCommand("MONITOR")
            C->>RS: MONITOR command
            RS->>C: Real-time key changes
            C->>E: New/modified keys
            E->>P: Analyze new data
            
            alt New sensitive data found
                P->>E: Pattern matches
                E->>U: Alert user of new findings
                E->>Export: Add to ongoing report
            end
        end
    end
    
    Export->>U: Display exfiltration results
    Export->>U: Provide download links for reports
    E->>U: Data exfiltration completed
```

## 7. Automated Chain Execution Flow

```mermaid
sequenceDiagram
    participant U as User
    participant Auto as Automation Engine
    participant Chain as Chain Controller
    participant M1 as Module 1
    participant M2 as Module 2
    participant M3 as Module N
    participant C as Redis Client
    participant RS as Redis Server
    participant Report as Reporting System
    
    U->>Auto: Select automation chain
    Auto->>U: Display chain options
    U->>Auto: Choose chain type
    
    Auto->>Chain: Load chain configuration
    Chain->>Chain: Parse execution sequence
    Chain->>Chain: Validate module dependencies
    
    alt Full Automated Assessment
        Note over Chain: Sequence: Recon → Exfiltration → WebShell → Persistence → Report
        
        Chain->>M1: Execute reconnaissance
        M1->>C: Multiple reconnaissance commands
        C->>RS: INFO, CONFIG GET, etc.
        RS->>C: Server information
        C->>M1: Reconnaissance data
        M1->>Chain: Step 1 completed successfully
        
        Chain->>Auto: Update progress (20%)
        Auto->>U: Step 1/5 completed
        
        Chain->>M2: Execute smart exfiltration
        M2->>C: Database scanning commands
        C->>RS: KEYS, GET, HGETALL, etc.
        RS->>C: Database contents
        C->>M2: Extracted data
        M2->>Chain: Step 2 completed successfully
        
        Chain->>Auto: Update progress (40%)
        Auto->>U: Step 2/5 completed
        
        Chain->>M3: Execute web shell injection
        M3->>C: File write commands
        C->>RS: CONFIG SET, SET, SAVE
        RS->>C: File write results
        C->>M3: Shell deployment status
        M3->>Chain: Step 3 completed successfully
        
        Chain->>Auto: Update progress (60%)
        Auto->>U: Step 3/5 completed
        
        Note over Chain: Continue with remaining modules...
        
    else Custom Chain
        U->>Auto: Select modules and sequence
        Auto->>Chain: Custom sequence configuration
        
        loop For each module in sequence
            Chain->>Chain: Load next module
            Chain->>M1: Execute module with parameters
            M1->>C: Module-specific commands
            C->>RS: Execute operations
            RS->>C: Operation results
            C->>M1: Processed results
            
            alt Module successful
                M1->>Chain: Success status
                Chain->>Chain: Continue to next module
                Chain->>Auto: Update progress
                Auto->>U: Step X completed
            else Module failed
                M1->>Chain: Error status
                Chain->>Chain: Evaluate error handling
                
                alt Retry allowed
                    Chain->>M1: Retry module execution
                else Skip and continue
                    Chain->>Chain: Skip to next module
                    Chain->>Auto: Log warning
                else Stop on failure
                    Chain->>Auto: Chain execution failed
                    Auto->>U: Display error and stop
                end
            end
        end
    end
    
    Chain->>Report: Generate comprehensive report
    Report->>Report: Compile all module results
    Report->>Report: Generate timeline analysis
    Report->>Report: Create risk assessment
    Report->>Report: Add remediation recommendations
    
    Report->>Auto: Final report ready
    Auto->>U: Chain execution completed
    Auto->>U: Display summary and report location
```

## 8. Error Handling and Recovery Flow

```mermaid
sequenceDiagram
    participant U as User
    participant App as Application
    participant C as Redis Client
    participant RS as Redis Server
    participant EH as Error Handler
    participant Log as Logger
    
    App->>C: Execute operation
    C->>RS: Send Redis command
    
    alt Connection Lost
        RS-->>C: Connection timeout/reset
        C->>EH: Connection error
        EH->>Log: Log connection failure
        EH->>EH: Attempt reconnection
        
        loop Retry attempts (max 3)
            EH->>C: Reconnect attempt
            C->>RS: New connection
            
            alt Reconnection successful
                RS->>C: Connection established
                C->>EH: Connection restored
                EH->>App: Retry original operation
                break
            else Reconnection failed
                RS-->>C: Connection failed
                EH->>EH: Wait and retry
            end
        end
        
        alt All retries failed
            EH->>App: Connection permanently failed
            App->>U: Display connection error
            App->>App: Graceful shutdown
        end
        
    else Authentication Error
        RS->>C: -NOAUTH Authentication required
        C->>EH: Authentication error
        EH->>Log: Log auth failure
        EH->>U: Prompt for credentials
        U->>EH: Provide password
        EH->>C: Retry with new credentials
        C->>RS: AUTH <new_password>
        
        alt Auth successful
            RS->>C: +OK
            C->>App: Authentication restored
            App->>App: Continue operation
        else Auth failed again
            RS->>C: -ERR invalid password
            C->>EH: Auth failed
            EH->>App: Authentication permanently failed
            App->>U: Display auth error and exit
        end
        
    else Command Error
        RS->>C: -ERR unknown command or syntax error
        C->>EH: Command error
        EH->>Log: Log command failure
        
        alt Recoverable error
            EH->>App: Command failed, try alternative
            App->>App: Use fallback method
        else Fatal error
            EH->>App: Unrecoverable command error
            App->>U: Display error message
            App->>App: Skip operation
        end
        
    else Resource Exhaustion
        RS->>C: -ERR OOM command not allowed
        C->>EH: Resource error
        EH->>Log: Log resource issue
        EH->>App: Resource exhaustion detected
        App->>App: Reduce operation scope
        App->>C: Retry with smaller dataset
        
    else Normal Operation
        RS->>C: +OK or data response
        C->>App: Successful result
        App->>U: Display operation result
    end
```

This comprehensive set of sequence diagrams covers all major operational flows in the 0xRedisis tool, showing the detailed interactions between components, error handling, and the complete user experience from startup to completion.