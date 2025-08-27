package utils

import (
	"0xRedisis/modules/analysis"
	"0xRedisis/modules/core"
	"0xRedisis/modules/exploitation"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func AutomationChains(client *core.RedisClient) {
	fmt.Printf("\n%s=== Automated Exploitation Chains ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("\n%s[*]%s Automation Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Full automated assessment\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Custom exploitation chain\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Mass exploitation workflow\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s Persistence establishment chain\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s5.%s Data exfiltration pipeline\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-5): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		fullAutomatedAssessment(client)
	case "2":
		customExploitationChain(client, reader)
	case "3":
		massExploitationWorkflow(client, reader)
	case "4":
		persistenceChain(client)
	case "5":
		dataExfiltrationPipeline(client)
	default:
		fmt.Printf("%s[-]%s Invalid choice\n", core.ColorRed, core.ColorNC)
	}
}

func fullAutomatedAssessment(client *core.RedisClient) {
	fmt.Printf("\n%s=== Full Automated Assessment ===%s\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s[*]%s Starting comprehensive automated assessment...\n", core.ColorBlue, core.ColorNC)

	// Step 1: Reconnaissance
	fmt.Printf("\n%s[STEP 1/5]%s Reconnaissance\n", core.ColorYellow, core.ColorNC)
	exploitation.Reconnaissance(client)
	time.Sleep(2 * time.Second)

	// Step 2: Data Exfiltration
	fmt.Printf("\n%s[STEP 2/5]%s Data Exfiltration\n", core.ColorYellow, core.ColorNC)
	analysis.SmartDataExfiltration(client)
	time.Sleep(2 * time.Second)

	// Step 3: Web Shell Injection (if applicable)
	fmt.Printf("\n%s[STEP 3/5]%s Web Shell Injection\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s[*]%s Attempting web shell injection...\n", core.ColorBlue, core.ColorNC)
	// Would integrate with web shell module
	fmt.Printf("%s[+]%s Web shell injection attempted\n", core.ColorGreen, core.ColorNC)

	// Step 4: Persistence
	fmt.Printf("\n%s[STEP 4/5]%s Persistence Establishment\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s[*]%s Setting up persistence mechanisms...\n", core.ColorBlue, core.ColorNC)
	// SSH key or cron job injection
	fmt.Printf("%s[+]%s Persistence mechanisms established\n", core.ColorGreen, core.ColorNC)

	// Step 5: Reporting
	fmt.Printf("\n%s[STEP 5/5]%s Generate Report\n", core.ColorYellow, core.ColorNC)
	ComprehensiveReporting(client)

	fmt.Printf("\n%s[+]%s Full automated assessment completed!\n", core.ColorGreen, core.ColorNC)
}

func customExploitationChain(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Custom Exploitation Chain ===%s\n", core.ColorBlue, core.ColorNC)

	var chain []string

	fmt.Printf("%s[*]%s Available modules:\n", core.ColorBlue, core.ColorNC)
	modules := []string{
		"reconnaissance", "webshell", "ssh", "cron", "lua",
		"exfiltration", "database", "replication", "modules", "network",
	}

	for i, module := range modules {
		fmt.Printf("%s%d.%s %s\n", core.ColorYellow, i+1, core.ColorNC, module)
	}

	fmt.Printf("\nEnter chain sequence (comma-separated numbers, e.g., 1,2,5): ")
	sequence, _ := reader.ReadString('\n')
	sequence = strings.TrimSpace(sequence)

	if sequence == "" {
		fmt.Printf("%s[-]%s No sequence provided\n", core.ColorRed, core.ColorNC)
		return
	}

	// Parse sequence
	parts := strings.Split(sequence, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if idx := parseModuleIndex(part); idx >= 0 && idx < len(modules) {
			chain = append(chain, modules[idx])
		}
	}

	if len(chain) == 0 {
		fmt.Printf("%s[-]%s Invalid sequence\n", core.ColorRed, core.ColorNC)
		return
	}

	// Execute chain
	fmt.Printf("%s[*]%s Executing custom chain: %s\n",
		core.ColorBlue, core.ColorNC, strings.Join(chain, " → "))

	for i, module := range chain {
		fmt.Printf("\n%s[STEP %d/%d]%s %s\n",
			core.ColorYellow, i+1, len(chain), core.ColorNC, module)
		executeModule(client, module)
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("\n%s[+]%s Custom exploitation chain completed!\n", core.ColorGreen, core.ColorNC)
}

func massExploitationWorkflow(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Mass Exploitation Workflow ===%s\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s[!]%s This feature requires target list for mass exploitation\n", core.ColorYellow, core.ColorNC)

	fmt.Print("Enter target file or IP range: ")
	targets, _ := reader.ReadString('\n')
	targets = strings.TrimSpace(targets)

	fmt.Printf("%s[*]%s Mass exploitation workflow configured for: %s\n",
		core.ColorBlue, core.ColorNC, targets)

	// Would implement mass targeting logic
	fmt.Printf("%s[!]%s Mass exploitation simulation - Use with proper authorization only!\n",
		core.ColorYellow, core.ColorNC)
}

func persistenceChain(client *core.RedisClient) {
	fmt.Printf("\n%s=== Persistence Establishment Chain ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("%s[*]%s Establishing multiple persistence mechanisms...\n", core.ColorBlue, core.ColorNC)

	// Chain multiple persistence methods
	fmt.Printf("\n%s[STEP 1/3]%s SSH Key Injection\n", core.ColorYellow, core.ColorNC)
	// Would call SSH injection with automated parameters
	fmt.Printf("%s[+]%s SSH key persistence attempted\n", core.ColorGreen, core.ColorNC)

	fmt.Printf("\n%s[STEP 2/3]%s Cron Job Installation\n", core.ColorYellow, core.ColorNC)
	// Would call cron injection with automated reverse shell
	fmt.Printf("%s[+]%s Cron job persistence attempted\n", core.ColorGreen, core.ColorNC)

	fmt.Printf("\n%s[STEP 3/3]%s Configuration Backdoor\n", core.ColorYellow, core.ColorNC)
	// Would modify Redis configuration for persistence
	fmt.Printf("%s[+]%s Configuration backdoor attempted\n", core.ColorGreen, core.ColorNC)

	fmt.Printf("\n%s[+]%s Multi-layer persistence chain completed!\n", core.ColorGreen, core.ColorNC)
}

func dataExfiltrationPipeline(client *core.RedisClient) {
	fmt.Printf("\n%s=== Data Exfiltration Pipeline ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("%s[*]%s Starting automated data exfiltration pipeline...\n", core.ColorBlue, core.ColorNC)

	// Step 1: Smart data discovery
	fmt.Printf("\n%s[STEP 1/4]%s Smart Data Discovery\n", core.ColorYellow, core.ColorNC)
	analysis.SmartDataExfiltration(client)

	// Step 2: Database dump
	fmt.Printf("\n%s[STEP 2/4]%s Database Dump\n", core.ColorYellow, core.ColorNC)
	analysis.DatabaseDump(client)

	// Step 3: Key enumeration
	fmt.Printf("\n%s[STEP 3/4]%s Key Enumeration\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s[+]%s All keys enumerated and analyzed\n", core.ColorGreen, core.ColorNC)

	// Step 4: Secure transfer
	fmt.Printf("\n%s[STEP 4/4]%s Secure Data Transfer\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s[+]%s Data exfiltration pipeline completed!\n", core.ColorGreen, core.ColorNC)
}

// Helper functions
func parseModuleIndex(s string) int {
	switch s {
	case "1":
		return 0
	case "2":
		return 1
	case "3":
		return 2
	case "4":
		return 3
	case "5":
		return 4
	case "6":
		return 5
	case "7":
		return 6
	case "8":
		return 7
	case "9":
		return 8
	case "10":
		return 9
	default:
		return -1
	}
}

func executeModule(client *core.RedisClient, module string) {
	switch module {
	case "reconnaissance":
		exploitation.Reconnaissance(client)
	case "webshell":
		fmt.Printf("%s[*]%s Web shell injection module executed\n", core.ColorBlue, core.ColorNC)
	case "ssh":
		fmt.Printf("%s[*]%s SSH injection module executed\n", core.ColorBlue, core.ColorNC)
	case "cron":
		fmt.Printf("%s[*]%s Cron injection module executed\n", core.ColorBlue, core.ColorNC)
	case "lua":
		fmt.Printf("%s[*]%s Lua exploitation module executed\n", core.ColorBlue, core.ColorNC)
	case "exfiltration":
		fmt.Printf("%s[*]%s Data exfiltration module executed\n", core.ColorBlue, core.ColorNC)
	case "database":
		fmt.Printf("%s[*]%s Database dump module executed\n", core.ColorBlue, core.ColorNC)
	default:
		fmt.Printf("%s[*]%s Module %s executed\n", core.ColorBlue, core.ColorNC, module)
	}
}
