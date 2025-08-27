package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type MenuOption struct {
	ID          string
	Title       string
	Description string
	Handler     func(*RedisClient)
}

type Menu struct {
	Title   string
	Options []MenuOption
}

func ShowMainMenu(client *RedisClient, handlers map[string]func(*RedisClient)) {
	for {
		fmt.Printf("\n%s=== Redis CTF Exploitation Tool ===%s\n", ColorBlue, ColorNC)
		fmt.Println("Connected to Redis. Select a category:")

		fmt.Printf("\n%s--- Basic Exploitation ---%s\n", ColorYellow, ColorNC)
		fmt.Printf("%s1.%s Reconnaissance (gather info)\n", ColorYellow, ColorNC)
		fmt.Printf("%s2.%s Inject PHP web shell\n", ColorYellow, ColorNC)
		fmt.Printf("%s3.%s Dump database to disk\n", ColorYellow, ColorNC)
		fmt.Printf("%s4.%s Inject SSH key for access\n", ColorYellow, ColorNC)
		fmt.Printf("%s5.%s Inject cron job\n", ColorYellow, ColorNC)

		fmt.Printf("\n%s--- Advanced Exploitation ---%s\n", ColorYellow, ColorNC)
		fmt.Printf("%s6.%s Lua script exploitation\n", ColorYellow, ColorNC)
		fmt.Printf("%s7.%s Master-slave replication abuse\n", ColorYellow, ColorNC)
		fmt.Printf("%s8.%s Redis module exploitation\n", ColorYellow, ColorNC)

		fmt.Printf("\n%s--- Analysis & Exfiltration ---%s\n", ColorYellow, ColorNC)
		fmt.Printf("%s9.%s Smart data exfiltration\n", ColorYellow, ColorNC)
		fmt.Printf("%s10.%s Network & cluster analysis\n", ColorYellow, ColorNC)
		fmt.Printf("%s11.%s Generate exploitation report\n", ColorYellow, ColorNC)

		fmt.Printf("\n%s--- Automation ---%s\n", ColorYellow, ColorNC)
		fmt.Printf("%s12.%s Auto exploitation chain\n", ColorYellow, ColorNC)
		fmt.Printf("%s13.%s Payload generator\n", ColorYellow, ColorNC)

		fmt.Printf("\n%s--- Listener Management ---%s\n", ColorYellow, ColorNC)
		fmt.Printf("%s14.%s Listener help & status\n", ColorYellow, ColorNC)
		fmt.Printf("%s15.%s Stop all listeners\n", ColorYellow, ColorNC)

		fmt.Printf("\n%s0.%s Exit\n", ColorYellow, ColorNC)
		fmt.Print("\nEnter your choice: ")

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if handler, ok := handlers["reconnaissance"]; ok {
				handler(client)
			}
		case "2":
			if handler, ok := handlers["webshell"]; ok {
				handler(client)
			}
		case "3":
			if handler, ok := handlers["database"]; ok {
				handler(client)
			}
		case "4":
			if handler, ok := handlers["ssh"]; ok {
				handler(client)
			}
		case "5":
			if handler, ok := handlers["cron"]; ok {
				handler(client)
			}
		case "6":
			if handler, ok := handlers["lua"]; ok {
				handler(client)
			}
		case "7":
			if handler, ok := handlers["replication"]; ok {
				handler(client)
			}
		case "8":
			if handler, ok := handlers["modules"]; ok {
				handler(client)
			}
		case "9":
			if handler, ok := handlers["exfiltration"]; ok {
				handler(client)
			}
		case "10":
			if handler, ok := handlers["network"]; ok {
				handler(client)
			}
		case "11":
			if handler, ok := handlers["reporting"]; ok {
				handler(client)
			}
		case "12":
			if handler, ok := handlers["automation"]; ok {
				handler(client)
			}
		case "13":
			if handler, ok := handlers["payloads"]; ok {
				handler(client)
			}
		case "14":
			ShowListenerHelp()
		case "15":
			StopAllListeners()
			fmt.Printf("%s[+]%s All listeners stopped\n", ColorGreen, ColorNC)
		case "0":
			fmt.Printf("%s[*]%s Stopping all listeners and exiting...\n", ColorBlue, ColorNC)
			StopAllListeners()
			return
		default:
			fmt.Printf("%s[-]%s Invalid choice. Please try again.\n", ColorRed, ColorNC)
		}

		// Ask if user wants to continue
		fmt.Print("\nPress Enter to return to main menu or 'q' to quit: ")
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(input)) == "q" {
			fmt.Printf("%s[*]%s Exiting...\n", ColorBlue, ColorNC)
			return
		}
	}
}
