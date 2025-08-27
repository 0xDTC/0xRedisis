# Data Flow Diagram (DFD) - 0xRedisis Redis Exploitation Tool

## Level 0 - Context Diagram

```mermaid
graph TD
    A[Security Tester/User] --> B[0xRedisis Tool]
    B --> C[Target Redis Server]
    B --> D[Report Files]
    B --> E[Reverse Shell Sessions]
    
    C --> B
    D --> A
    E --> A
```

## Level 1 - System Overview

```mermaid
graph TB
    subgraph "User Interface Layer"
        A[Command Line Interface]
        B[Menu System]
        C[User Input Handler]
    end
    
    subgraph "Core Framework"
        D[Main Application Controller]
        E[Redis Client Engine]
        F[Connection Manager]
        G[Authentication Handler]
    end
    
    subgraph "Exploitation Engine"
        H[Module Dispatcher]
        I[Command Executor]
        J[Payload Generator]
        K[Reverse Shell Listener]
    end
    
    subgraph "Data Processing"
        L[Response Parser]
        M[Data Analyzer]
        N[Pattern Matcher]
        O[Result Formatter]
    end
    
    subgraph "Output Systems"
        P[Console Display]
        Q[File Export]
        R[Report Generator]
        S[Logging System]
    end
    
    subgraph "External Systems"
        T[Target Redis Server]
        U[File System]
        V[Network Interfaces]
        W[System Shell]
    end
    
    %% Data Flow Connections
    A --> B
    B --> C
    C --> D
    D --> H
    H --> I
    I --> E
    E --> F
    F --> G
    G --> T
    
    T --> L
    L --> M
    M --> N
    N --> O
    O --> P
    O --> Q
    O --> R
    
    J --> K
    K --> V
    I --> J
    Q --> U
    R --> U
    S --> U
    
    I --> W
    W --> I
```

## Level 2 - Detailed Module Flow

### Exploitation Module Data Flow

```mermaid
graph LR
    subgraph "Reconnaissance Module"
        A1[Server Info Request] --> A2[CONFIG GET *]
        A2 --> A3[INFO Command]
        A3 --> A4[CLIENT LIST]
        A4 --> A5[Memory Analysis]
        A5 --> A6[Keyspace Stats]
    end
    
    subgraph "Web Shell Module"
        B1[Directory Detection] --> B2[Payload Selection]
        B2 --> B3[Shell Generation]
        B3 --> B4[File Write Command]
        B4 --> B5[Persistence Check]
    end
    
    subgraph "SSH Key Module"
        C1[Key Generation] --> C2[Path Enumeration]
        C2 --> C3[authorized_keys Write]
        C3 --> C4[Permission Setting]
        C4 --> C5[Connection Test]
    end
    
    subgraph "Database Module"
        D1[Database Scan] --> D2[Key Enumeration]
        D2 --> D3[Data Extraction]
        D3 --> D4[Format Conversion]
        D4 --> D5[Export Generation]
    end
    
    subgraph "Lua Script Module"
        E1[Script Selection] --> E2[Command Wrapping]
        E2 --> E3[EVAL Execution]
        E3 --> E4[Output Parsing]
        E4 --> E5[Shell Handling]
    end
```

### Analysis Module Data Flow

```mermaid
graph TD
    subgraph "Smart Data Exfiltration"
        F1[Pattern Library] --> F2[Key Scanning]
        F2 --> F3[Content Analysis]
        F3 --> F4[Sensitivity Scoring]
        F4 --> F5[Classification]
        F5 --> F6[Export Preparation]
    end
    
    subgraph "Network Analysis"
        G1[Port Scanning] --> G2[Service Detection]
        G2 --> G3[Cluster Discovery]
        G3 --> G4[Topology Mapping]
        G4 --> G5[Multi-target Planning]
    end
    
    subgraph "Reporting System"
        H1[Activity Logging] --> H2[Evidence Collection]
        H2 --> H3[Timeline Generation]
        H3 --> H4[Risk Assessment]
        H4 --> H5[Report Templates]
        H5 --> H6[Multi-format Export]
    end
```

## Level 3 - Core Component Interactions

### Redis Client Communication Flow

```mermaid
sequenceDiagram
    participant UI as User Interface
    participant MC as Main Controller
    participant RC as Redis Client
    participant CP as Command Parser
    participant NW as Network Layer
    participant RS as Redis Server
    
    UI->>MC: User Command
    MC->>RC: Execute Operation
    RC->>CP: Parse Command
    CP->>NW: RESP Protocol Message
    NW->>RS: TCP Connection
    RS->>NW: Response Data
    NW->>CP: Raw Response
    CP->>RC: Parsed Result
    RC->>MC: Structured Data
    MC->>UI: Formatted Output
```

### Automated Chain Execution Flow

```mermaid
graph TD
    subgraph "Chain Controller"
        I1[Chain Selection] --> I2[Module Sequence]
        I2 --> I3[Dependency Check]
        I3 --> I4[Execution Planning]
    end
    
    subgraph "Step Executor"
        J1[Module Loading] --> J2[Parameter Setup]
        J2 --> J3[Execution Wrapper]
        J3 --> J4[Success Validation]
        J4 --> J5[Error Handling]
    end
    
    subgraph "State Manager"
        K1[Chain State] --> K2[Module Results]
        K2 --> K3[Progress Tracking]
        K3 --> K4[Rollback Capability]
    end
    
    I4 --> J1
    J5 --> K1
    K4 --> I2
```

## Data Store Interactions

### File System Operations

```mermaid
graph LR
    subgraph "Input Files"
        L1[Target Lists]
        L2[Payload Templates]
        L3[Configuration Files]
    end
    
    subgraph "Processing"
        M1[File Reader]
        M2[Template Engine]
        M3[Config Parser]
    end
    
    subgraph "Output Files"
        N1[Exploitation Reports]
        N2[Extracted Data]
        N3[Generated Payloads]
        N4[Activity Logs]
        N5[Shell Scripts]
    end
    
    L1 --> M1
    L2 --> M2
    L3 --> M3
    
    M1 --> N1
    M2 --> N3
    M2 --> N5
    M3 --> N4
    M1 --> N2
```

### Network Communications

```mermaid
graph TB
    subgraph "Outbound Connections"
        O1[Redis Protocol<br/>Port 6379]
        O2[HTTP Requests<br/>Port 80/443]
        O3[Custom Listeners<br/>Dynamic Ports]
    end
    
    subgraph "Inbound Connections"
        P1[Reverse Shells]
        P2[Data Exfiltration]
        P3[Health Checks]
    end
    
    subgraph "Protocol Handlers"
        Q1[RESP Parser]
        Q2[HTTP Client]
        Q3[TCP Server]
        Q4[Shell Handler]
    end
    
    O1 --> Q1
    O2 --> Q2
    O3 --> Q3
    
    P1 --> Q4
    P2 --> Q3
    P3 --> Q2
```

## Security Data Flow

### Authentication and Authorization

```mermaid
graph TD
    subgraph "Credential Handling"
        R1[Command Line Args] --> R2[Password Parsing]
        R2 --> R3[Auth Validation]
        R3 --> R4[Session Management]
    end
    
    subgraph "Secure Operations"
        S1[Memory Protection] --> S2[Credential Cleanup]
        S2 --> S3[Secure Logging]
        S3 --> S4[Audit Trail]
    end
    
    subgraph "Access Control"
        T1[Permission Checks] --> T2[Command Filtering]
        T2 --> T3[Rate Limiting]
        T3 --> T4[Connection Limits]
    end
    
    R4 --> S1
    S4 --> T1
```

## Performance and Monitoring Flow

### Resource Management

```mermaid
graph LR
    subgraph "Resource Monitoring"
        U1[Memory Usage] --> U2[Connection Pool]
        U2 --> U3[Thread Management]
        U3 --> U4[Timeout Handling]
    end
    
    subgraph "Performance Optimization"
        V1[Concurrent Operations] --> V2[Batch Processing]
        V2 --> V3[Result Caching]
        V3 --> V4[Lazy Loading]
    end
    
    subgraph "Error Recovery"
        W1[Connection Recovery] --> W2[Retry Logic]
        W2 --> W3[Graceful Degradation]
        W3 --> W4[Cleanup Operations]
    end
    
    U4 --> V1
    V4 --> W1
```

## Summary

This DFD illustrates the complete data flow architecture of the 0xRedisis tool, showing how data moves through the system from user input to final output. The multi-level approach provides both high-level system understanding and detailed component interactions, ensuring comprehensive coverage of all data processing aspects.

Key characteristics:
- **Modular Design**: Clear separation of concerns across different processing layers
- **Bidirectional Flow**: Data flows both to and from the target Redis server
- **Security Focus**: Proper handling of sensitive authentication data
- **Scalability**: Support for concurrent operations and batch processing
- **Reliability**: Error handling and recovery mechanisms throughout the flow
- **Extensibility**: Pluggable module architecture for easy feature additions