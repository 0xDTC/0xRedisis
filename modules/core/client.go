package core

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	// Color codes for terminal output
	ColorRed    = "\033[0;31m"
	ColorGreen  = "\033[0;32m"
	ColorYellow = "\033[1;33m"
	ColorBlue   = "\033[0;34m"
	ColorNC     = "\033[0m" // No Color
)

type Config struct {
	Host     string
	Port     int
	Password string
	Timeout  time.Duration
}

type RedisClient struct {
	Config Config
	conn   net.Conn
}

func NewRedisClient(config Config) *RedisClient {
	return &RedisClient{Config: config}
}

func (c *RedisClient) Connect() error {
	fmt.Printf("%s[*]%s Testing Redis connection on %s%s:%d%s...\n",
		ColorBlue, ColorNC, ColorYellow, c.Config.Host, c.Config.Port, ColorNC)

	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", c.Config.Host, c.Config.Port),
		c.Config.Timeout)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	c.conn = conn

	// Test connection with INFO command
	info, err := c.SendCommand("INFO", "server")
	if err != nil {
		return fmt.Errorf("failed to get server info: %v", err)
	}

	// Parse version and OS from INFO response
	version := ExtractInfoValue(info, "redis_version")
	os := ExtractInfoValue(info, "os")

	fmt.Printf("%s[+]%s Connection successful!\n", ColorGreen, ColorNC)
	fmt.Printf("%s[+]%s Redis version: %s%s%s\n", ColorGreen, ColorNC, ColorYellow, version, ColorNC)
	fmt.Printf("%s[+]%s OS: %s%s%s\n", ColorGreen, ColorNC, ColorYellow, os, ColorNC)

	// Test if Redis is responsive
	pong, err := c.SendCommand("PING")
	if err == nil && strings.Contains(pong, "PONG") {
		fmt.Printf("%s[+]%s Redis is responsive\n", ColorGreen, ColorNC)
	}

	return nil
}

func (c *RedisClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *RedisClient) SendCommand(args ...string) (string, error) {
	if c.conn == nil {
		return "", fmt.Errorf("not connected to Redis")
	}

	// Authenticate if password is provided
	if c.Config.Password != "" {
		authCmd := fmt.Sprintf("*2\r\n$4\r\nAUTH\r\n$%d\r\n%s\r\n",
			len(c.Config.Password), c.Config.Password)
		if _, err := c.conn.Write([]byte(authCmd)); err != nil {
			return "", fmt.Errorf("auth failed: %v", err)
		}

		// Read auth response
		c.conn.SetReadDeadline(time.Now().Add(c.Config.Timeout))
		reader := bufio.NewReader(c.conn)
		authResp, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("auth response error: %v", err)
		}
		if !strings.Contains(authResp, "+OK") {
			return "", fmt.Errorf("authentication failed")
		}
	}

	// Build Redis command in RESP protocol format
	cmd := fmt.Sprintf("*%d\r\n", len(args))
	for _, arg := range args {
		cmd += fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
	}

	// Send command
	c.conn.SetWriteDeadline(time.Now().Add(c.Config.Timeout))
	if _, err := c.conn.Write([]byte(cmd)); err != nil {
		return "", fmt.Errorf("write error: %v", err)
	}

	// Read response
	c.conn.SetReadDeadline(time.Now().Add(c.Config.Timeout))
	reader := bufio.NewReader(c.conn)
	response := ""

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		response += line

		// Simple response parsing - break on complete response
		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") ||
			strings.HasPrefix(line, ":") {
			break
		}

		// For bulk strings, read the specified number of bytes
		if strings.HasPrefix(line, "$") {
			lengthStr := strings.TrimSpace(line[1:])
			if length, err := strconv.Atoi(lengthStr); err == nil && length > 0 {
				data := make([]byte, length+2) // +2 for \r\n
				reader.Read(data)
				response += string(data[:length])
				break
			}
		}
	}

	if strings.HasPrefix(response, "-ERR") {
		return "", fmt.Errorf("Redis error: %s", response)
	}

	return response, nil
}

func ExtractInfoValue(info, key string) string {
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, key+":") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "unknown"
}

func ParseConfigValue(response, key string) string {
	lines := strings.Split(response, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == key && i+1 < len(lines) {
			return strings.TrimSpace(lines[i+1])
		}
	}
	return "unknown"
}
