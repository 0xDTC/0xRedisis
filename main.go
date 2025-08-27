package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"0xRedisis/modules/analysis"
	"0xRedisis/modules/core"
	"0xRedisis/modules/exploitation"
	"0xRedisis/modules/persistence"
	"0xRedisis/modules/utils"
)

func main() {
	if len(os.Args) < 3 {
		showHelp()
		return
	}

	config, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("%s[-]%s %s\n", core.ColorRed, core.ColorNC, err.Error())
		os.Exit(1)
	}

	client := core.NewRedisClient(config)

	if err := client.Connect(); err != nil {
		fmt.Printf("%s[-]%s Failed to connect: %s\n", core.ColorRed, core.ColorNC, err.Error())
		os.Exit(1)
	}
	defer client.Close()

	// Create handlers map for menu system
	handlers := map[string]func(*core.RedisClient){
		"reconnaissance": exploitation.Reconnaissance,
		"webshell":       exploitation.WebShellInjection,
		"database":       analysis.DatabaseDump,
		"ssh":            persistence.SSHKeyInjection,
		"cron":           persistence.CronJobInjection,
		"lua":            exploitation.LuaExploitation,
		"replication":    exploitation.ReplicationAbuse,
		"modules":        exploitation.RedisModuleExploitation,
		"exfiltration":   analysis.SmartDataExfiltration,
		"network":        analysis.NetworkAnalysis,
		"reporting":      utils.ComprehensiveReporting,
		"automation":     utils.AutomationChains,
		"payloads":       utils.PayloadGenerator,
	}

	core.ShowMainMenu(client, handlers)
}

func showHelp() {
	fmt.Printf("\n%sRedis CTF Exploitation Tool v2.0%s\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%sUsage:%s %s <host> <port> [password]\n", core.ColorYellow, core.ColorNC, os.Args[0])
	fmt.Printf("%sDescription:%s Advanced Redis exploitation tool for CTF scenarios and security testing\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%sExamples:%s\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("  %s 192.168.1.100 6379\n", os.Args[0])
	fmt.Printf("  %s target.ctf.com 6379 mypassword\n", os.Args[0])
	fmt.Println("\nIf no password is provided, the script assumes no authentication is required.")

	fmt.Printf("\n%sFeatures:%s\n", core.ColorYellow, core.ColorNC)
	fmt.Println("  • Basic exploitation (webshells, SSH, cron, database dumps)")
	fmt.Println("  • Advanced techniques (Lua scripts, replication abuse, modules)")
	fmt.Println("  • Smart data analysis and exfiltration")
	fmt.Println("  • Automated exploitation chains")
	fmt.Println("  • Comprehensive reporting")
}

func parseArgs(args []string) (core.Config, error) {
	if len(args) < 2 {
		return core.Config{}, fmt.Errorf("insufficient arguments")
	}

	port, err := strconv.Atoi(args[1])
	if err != nil || port < 1 || port > 65535 {
		return core.Config{}, fmt.Errorf("invalid port number. Must be between 1-65535")
	}

	config := core.Config{
		Host:    args[0],
		Port:    port,
		Timeout: 10 * time.Second,
	}

	if len(args) > 2 {
		config.Password = args[2]
	}

	return config, nil
}
