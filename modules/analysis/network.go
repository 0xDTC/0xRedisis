package analysis

import (
	"0xRedisis/modules/core"
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RedisInstance struct {
	Host     string
	Port     int
	Version  string
	Role     string
	Password bool
	Slaves   int
	Keys     int
	Memory   string
}

func NetworkAnalysis(client *core.RedisClient) {
	fmt.Printf("\n%s=== Network & Cluster Analysis ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("\n%s[*]%s Network Analysis Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Redis cluster discovery\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Network range Redis scanning\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Sentinel discovery and analysis\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s Multi-instance exploitation\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s5.%s Network topology mapping\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-5): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		clusterDiscovery(client, reader)
	case "2":
		networkRangeScan(client, reader)
	case "3":
		sentinelDiscovery(client, reader)
	case "4":
		multiInstanceExploit(client, reader)
	case "5":
		networkTopologyMapping(client, reader)
	default:
		fmt.Printf("%s[-]%s Invalid choice\n", core.ColorRed, core.ColorNC)
	}
}

func clusterDiscovery(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Redis Cluster Discovery ===%s\n", core.ColorBlue, core.ColorNC)

	// Check if current instance is part of a cluster
	fmt.Printf("%s[*]%s Checking cluster configuration...\n", core.ColorBlue, core.ColorNC)

	clusterInfo, err := client.SendCommand("CLUSTER", "INFO")
	if err != nil {
		fmt.Printf("%s[-]%s Cluster commands not available: %s\n", core.ColorRed, core.ColorNC, err.Error())
		fmt.Printf("%s[*]%s This Redis instance is not part of a cluster\n", core.ColorBlue, core.ColorNC)
		return
	}

	fmt.Printf("%s[+]%s Cluster information found:\n", core.ColorGreen, core.ColorNC)
	fmt.Printf("%s%s%s\n", core.ColorYellow, clusterInfo, core.ColorNC)

	// Get cluster nodes
	nodes, err := client.SendCommand("CLUSTER", "NODES")
	if err != nil {
		fmt.Printf("%s[-]%s Failed to get cluster nodes: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fmt.Printf("\n%s[*]%s Parsing cluster nodes...\n", core.ColorBlue, core.ColorNC)
	clusterNodes := parseClusterNodes(nodes)

	fmt.Printf("%s[+]%s Found %d cluster nodes:\n", core.ColorGreen, core.ColorNC, len(clusterNodes))

	for i, node := range clusterNodes {
		fmt.Printf("%s[%d]%s %s:%d (%s) - %s\n",
			core.ColorYellow, i+1, core.ColorNC, node.Host, node.Port, node.Role,
			getNodeStatus(node))
	}

	// Offer to connect to other cluster nodes
	fmt.Print("\nConnect to another cluster node? (y/n): ")
	connect, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(connect)) == "y" {
		selectAndConnectToNode(clusterNodes, reader)
	}
}

func networkRangeScan(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Network Range Redis Scanning ===%s\n", core.ColorBlue, core.ColorNC)

	// Get target range
	fmt.Print("Enter IP range (e.g., 192.168.1.0/24 or 192.168.1.1-192.168.1.100): ")
	ipRange, _ := reader.ReadString('\n')
	ipRange = strings.TrimSpace(ipRange)

	fmt.Print("Enter port range (e.g., 6379 or 6379-6389): ")
	portRange, _ := reader.ReadString('\n')
	portRange = strings.TrimSpace(portRange)

	fmt.Print("Enter number of threads (default 50): ")
	threadsStr, _ := reader.ReadString('\n')
	threads := 50
	if t := strings.TrimSpace(threadsStr); t != "" {
		if parsed, err := strconv.Atoi(t); err == nil {
			threads = parsed
		}
	}

	// Parse IP range
	ips := parseIPRange(ipRange)
	if len(ips) == 0 {
		fmt.Printf("%s[-]%s Invalid IP range\n", core.ColorRed, core.ColorNC)
		return
	}

	// Parse port range
	ports := parsePortRange(portRange)
	if len(ports) == 0 {
		fmt.Printf("%s[-]%s Invalid port range\n", core.ColorRed, core.ColorNC)
		return
	}

	fmt.Printf("%s[*]%s Scanning %d IPs across %d ports with %d threads...\n",
		core.ColorBlue, core.ColorNC, len(ips), len(ports), threads)

	// Perform concurrent scan
	instances := concurrentRedisScan(ips, ports, threads)

	fmt.Printf("%s[+]%s Found %d Redis instances:\n", core.ColorGreen, core.ColorNC, len(instances))

	for i, instance := range instances {
		fmt.Printf("%s[%d]%s %s:%d - %s (%s) - %d keys\n",
			core.ColorYellow, i+1, core.ColorNC, instance.Host, instance.Port,
			instance.Version, instance.Role, instance.Keys)
	}

	// Save results
	if len(instances) > 0 {
		saveResults := askYesNo(reader, "Save scan results to file?")
		if saveResults {
			saveScanResults(instances)
		}
	}
}

func sentinelDiscovery(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Redis Sentinel Discovery ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("%s[*]%s Checking if current instance is a Sentinel...\n", core.ColorBlue, core.ColorNC)

	// Try Sentinel commands
	masters, err := client.SendCommand("SENTINEL", "MASTERS")
	if err != nil {
		fmt.Printf("%s[-]%s Not a Sentinel instance: %s\n", core.ColorRed, core.ColorNC, err.Error())

		// Try to discover Sentinels in network
		fmt.Print("Enter potential Sentinel IP (or press Enter to skip): ")
		sentinelIP, _ := reader.ReadString('\n')
		sentinelIP = strings.TrimSpace(sentinelIP)

		if sentinelIP != "" {
			discoverSentinelsInNetwork(sentinelIP, reader)
		}
		return
	}

	fmt.Printf("%s[+]%s This is a Redis Sentinel instance!\n", core.ColorGreen, core.ColorNC)
	fmt.Printf("%s[*]%s Masters information:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s%s%s\n", core.ColorYellow, masters, core.ColorNC)

	// Get sentinel info
	info, _ := client.SendCommand("INFO", "sentinel")
	if info != "" {
		fmt.Printf("\n%s[*]%s Sentinel info:\n", core.ColorBlue, core.ColorNC)
		fmt.Printf("%s%s%s\n", core.ColorYellow, info, core.ColorNC)
	}

	// Get other sentinels
	fmt.Printf("\n%s[*]%s Discovering other Sentinels...\n", core.ColorBlue, core.ColorNC)

	// Parse master names and get sentinels for each
	masterNames := parseMasterNames(masters)
	for _, masterName := range masterNames {
		sentinels, err := client.SendCommand("SENTINEL", "SENTINELS", masterName)
		if err == nil {
			fmt.Printf("%s[+]%s Sentinels for master '%s':\n", core.ColorGreen, core.ColorNC, masterName)
			fmt.Printf("%s%s%s\n", core.ColorYellow, sentinels, core.ColorNC)
		}
	}
}

func multiInstanceExploit(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Multi-Instance Exploitation ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Print("Load Redis instances from scan file? (y/n): ")
	loadFile, _ := reader.ReadString('\n')

	var instances []RedisInstance

	if strings.ToLower(strings.TrimSpace(loadFile)) == "y" {
		instances = loadScanResults(reader)
	} else {
		// Manual entry
		instances = manualInstanceEntry(reader)
	}

	if len(instances) == 0 {
		fmt.Printf("%s[-]%s No instances to exploit\n", core.ColorRed, core.ColorNC)
		return
	}

	fmt.Printf("\n%s[*]%s Exploitation Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Mass reconnaissance\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Mass web shell injection\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Mass SSH key injection\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s Mass data exfiltration\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter choice (1-4): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		massReconnaissance(instances)
	case "2":
		massWebShellInjection(instances, reader)
	case "3":
		massSSHKeyInjection(instances, reader)
	case "4":
		massDataExfiltration(instances, reader)
	}
}

func networkTopologyMapping(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Network Topology Mapping ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("%s[*]%s Starting comprehensive network topology discovery...\n", core.ColorBlue, core.ColorNC)

	// Start from current instance
	rootInstance := RedisInstance{
		Host: client.Config.Host,
		Port: client.Config.Port,
	}

	// Discover connected instances
	topology := make(map[string][]RedisInstance)
	visited := make(map[string]bool)

	discoverTopologyRecursive(rootInstance, topology, visited, 0, 3) // Max depth 3

	// Display topology
	fmt.Printf("\n%s[+]%s Network Topology Discovered:\n", core.ColorGreen, core.ColorNC)
	displayTopology(topology)

	// Generate topology report
	generateTopologyReport(topology)
}

// Helper functions

func parseClusterNodes(nodes string) []RedisInstance {
	var instances []RedisInstance
	lines := strings.Split(nodes, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			hostPort := strings.Split(fields[1], ":")
			if len(hostPort) >= 2 {
				port, _ := strconv.Atoi(hostPort[1])
				role := "unknown"
				if len(fields) >= 3 {
					role = fields[2]
				}

				instances = append(instances, RedisInstance{
					Host: hostPort[0],
					Port: port,
					Role: role,
				})
			}
		}
	}

	return instances
}

func parseIPRange(ipRange string) []string {
	var ips []string

	if strings.Contains(ipRange, "/") {
		// CIDR notation
		_, network, err := net.ParseCIDR(ipRange)
		if err != nil {
			return ips
		}

		// Generate IPs in network (simplified)
		ip := network.IP
		for network.Contains(ip) {
			ips = append(ips, ip.String())
			ip = nextIP(ip)
			if len(ips) > 254 { // Safety limit
				break
			}
		}
	} else if strings.Contains(ipRange, "-") {
		// Range notation (simplified)
		parts := strings.Split(ipRange, "-")
		if len(parts) == 2 {
			// Simple implementation for same subnet
			ips = append(ips, parts[0], parts[1])
		}
	} else {
		// Single IP
		ips = append(ips, ipRange)
	}

	return ips
}

func parsePortRange(portRange string) []int {
	var ports []int

	if strings.Contains(portRange, "-") {
		parts := strings.Split(portRange, "-")
		if len(parts) == 2 {
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			for i := start; i <= end; i++ {
				ports = append(ports, i)
			}
		}
	} else {
		port, _ := strconv.Atoi(portRange)
		ports = append(ports, port)
	}

	return ports
}

func concurrentRedisScan(ips []string, ports []int, threads int) []RedisInstance {
	var instances []RedisInstance
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Create work channel
	work := make(chan string, len(ips)*len(ports))

	// Add all IP:port combinations to work channel
	for _, ip := range ips {
		for _, port := range ports {
			work <- fmt.Sprintf("%s:%d", ip, port)
		}
	}
	close(work)

	// Start worker goroutines
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for target := range work {
				parts := strings.Split(target, ":")
				if len(parts) != 2 {
					continue
				}

				ip := parts[0]
				port, _ := strconv.Atoi(parts[1])

				if instance := scanRedisInstance(ip, port); instance != nil {
					mutex.Lock()
					instances = append(instances, *instance)
					mutex.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	return instances
}

func scanRedisInstance(host string, port int) *RedisInstance {
	// Try to connect with timeout
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 2*time.Second)
	if err != nil {
		return nil
	}
	defer conn.Close()

	// Try Redis PING command
	conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))

	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return nil
	}

	response := string(buffer[:n])
	if !strings.Contains(response, "PONG") {
		return nil
	}

	// Get basic info
	instance := &RedisInstance{
		Host: host,
		Port: port,
	}

	// Try to get more detailed info
	config := core.Config{
		Host:    host,
		Port:    port,
		Timeout: 2 * time.Second,
	}

	tempClient := core.NewRedisClient(config)
	if err := tempClient.Connect(); err == nil {
		if info, err := tempClient.SendCommand("INFO", "server"); err == nil {
			instance.Version = core.ExtractInfoValue(info, "redis_version")
			instance.Role = core.ExtractInfoValue(info, "role")
		}

		if dbsize, err := tempClient.SendCommand("DBSIZE"); err == nil {
			instance.Keys, _ = strconv.Atoi(strings.TrimSpace(dbsize))
		}

		tempClient.Close()
	}

	return instance
}

func nextIP(ip net.IP) net.IP {
	next := make(net.IP, len(ip))
	copy(next, ip)
	for i := len(next) - 1; i >= 0; i-- {
		next[i]++
		if next[i] != 0 {
			break
		}
	}
	return next
}

func getNodeStatus(node RedisInstance) string {
	if strings.Contains(node.Role, "master") {
		return "Master"
	} else if strings.Contains(node.Role, "slave") {
		return "Slave"
	}
	return "Unknown"
}

func selectAndConnectToNode(nodes []RedisInstance, reader *bufio.Reader) {
	fmt.Print("Enter node number to connect to: ")
	nodeStr, _ := reader.ReadString('\n')
	nodeNum, _ := strconv.Atoi(strings.TrimSpace(nodeStr))

	if nodeNum > 0 && nodeNum <= len(nodes) {
		node := nodes[nodeNum-1]
		fmt.Printf("%s[*]%s Attempting connection to %s:%d...\n",
			core.ColorBlue, core.ColorNC, node.Host, node.Port)

		config := core.Config{
			Host:    node.Host,
			Port:    node.Port,
			Timeout: 5 * time.Second,
		}

		newClient := core.NewRedisClient(config)
		if err := newClient.Connect(); err == nil {
			fmt.Printf("%s[+]%s Connected successfully!\n", core.ColorGreen, core.ColorNC)
			// Here you could switch context or open new session
		} else {
			fmt.Printf("%s[-]%s Connection failed: %s\n", core.ColorRed, core.ColorNC, err.Error())
		}
	}
}

func askYesNo(reader *bufio.Reader, question string) bool {
	fmt.Printf("%s (y/n): ", question)
	answer, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(answer)) == "y"
}

func saveScanResults(instances []RedisInstance) {
	filename := fmt.Sprintf("redis_scan_results_%s.txt", time.Now().Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to create results file: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "Redis Network Scan Results\n")
	fmt.Fprintf(file, "Timestamp: %s\n\n", time.Now().Format(time.RFC3339))

	for _, instance := range instances {
		fmt.Fprintf(file, "%s:%d | %s | %s | %d keys\n",
			instance.Host, instance.Port, instance.Version, instance.Role, instance.Keys)
	}

	fmt.Printf("%s[+]%s Results saved to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}

func loadScanResults(reader *bufio.Reader) []RedisInstance {
	fmt.Print("Enter scan results filename: ")
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	// Implementation would parse the saved file format
	fmt.Printf("%s[!]%s File loading not implemented in demo\n", core.ColorYellow, core.ColorNC)
	return []RedisInstance{}
}

func manualInstanceEntry(reader *bufio.Reader) []RedisInstance {
	var instances []RedisInstance

	fmt.Print("Enter number of instances to add: ")
	countStr, _ := reader.ReadString('\n')
	count, _ := strconv.Atoi(strings.TrimSpace(countStr))

	for i := 0; i < count; i++ {
		fmt.Printf("Instance %d:\n", i+1)
		fmt.Print("  Host: ")
		host, _ := reader.ReadString('\n')

		fmt.Print("  Port: ")
		portStr, _ := reader.ReadString('\n')
		port, _ := strconv.Atoi(strings.TrimSpace(portStr))

		instances = append(instances, RedisInstance{
			Host: strings.TrimSpace(host),
			Port: port,
		})
	}

	return instances
}

func massReconnaissance(instances []RedisInstance) {
	fmt.Printf("%s[*]%s Performing mass reconnaissance on %d instances...\n",
		core.ColorBlue, core.ColorNC, len(instances))

	for i, instance := range instances {
		fmt.Printf("\r%s[*]%s Progress: %d/%d", core.ColorBlue, core.ColorNC, i+1, len(instances))

		config := core.Config{
			Host:    instance.Host,
			Port:    instance.Port,
			Timeout: 3 * time.Second,
		}

		client := core.NewRedisClient(config)
		if err := client.Connect(); err == nil {
			// Basic recon
			if info, err := client.SendCommand("INFO", "server"); err == nil {
				instance.Version = core.ExtractInfoValue(info, "redis_version")
			}
			client.Close()
		}
	}

	fmt.Printf("\n%s[+]%s Mass reconnaissance completed\n", core.ColorGreen, core.ColorNC)
}

func massWebShellInjection(instances []RedisInstance, reader *bufio.Reader) {
	fmt.Printf("%s[!]%s Mass web shell injection - Use with caution!\n", core.ColorYellow, core.ColorNC)
	fmt.Print("Continue? (y/n): ")

	if !askYesNo(reader, "") {
		return
	}

	fmt.Printf("%s[*]%s Starting mass web shell injection...\n", core.ColorBlue, core.ColorNC)
	// Implementation would iterate through instances and inject shells
	fmt.Printf("%s[!]%s Mass exploitation feature - implementation details omitted\n", core.ColorYellow, core.ColorNC)
}

func massSSHKeyInjection(instances []RedisInstance, reader *bufio.Reader) {
	fmt.Printf("%s[!]%s Mass SSH key injection - Use with caution!\n", core.ColorYellow, core.ColorNC)
	// Similar to web shell injection
	fmt.Printf("%s[!]%s Mass exploitation feature - implementation details omitted\n", core.ColorYellow, core.ColorNC)
}

func massDataExfiltration(instances []RedisInstance, reader *bufio.Reader) {
	fmt.Printf("%s[*]%s Starting mass data exfiltration...\n", core.ColorBlue, core.ColorNC)
	// Would connect to each instance and extract data
	fmt.Printf("%s[!]%s Mass exfiltration feature - implementation details omitted\n", core.ColorYellow, core.ColorNC)
}

func discoverSentinelsInNetwork(sentinelIP string, reader *bufio.Reader) {
	fmt.Printf("%s[*]%s Attempting to connect to potential Sentinel at %s\n",
		core.ColorBlue, core.ColorNC, sentinelIP)

	// Try common Sentinel ports
	sentinelPorts := []int{26379, 26380, 26381}

	for _, port := range sentinelPorts {
		config := core.Config{
			Host:    sentinelIP,
			Port:    port,
			Timeout: 3 * time.Second,
		}

		client := core.NewRedisClient(config)
		if err := client.Connect(); err == nil {
			masters, err := client.SendCommand("SENTINEL", "MASTERS")
			if err == nil {
				fmt.Printf("%s[+]%s Found Sentinel at %s:%d\n", core.ColorGreen, core.ColorNC, sentinelIP, port)
				fmt.Printf("%s%s%s\n", core.ColorYellow, masters, core.ColorNC)
			}
			client.Close()
		}
	}
}

func parseMasterNames(masters string) []string {
	var names []string
	// Simple parsing - would need more robust implementation
	lines := strings.Split(masters, "\n")
	for _, line := range lines {
		if strings.Contains(line, "name") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "name" && i+1 < len(fields) {
					names = append(names, fields[i+1])
				}
			}
		}
	}
	return names
}

func discoverTopologyRecursive(instance RedisInstance, topology map[string][]RedisInstance, visited map[string]bool, depth, maxDepth int) {
	if depth >= maxDepth {
		return
	}

	key := fmt.Sprintf("%s:%d", instance.Host, instance.Port)
	if visited[key] {
		return
	}
	visited[key] = true

	// Connect and discover related instances
	config := core.Config{
		Host:    instance.Host,
		Port:    instance.Port,
		Timeout: 3 * time.Second,
	}

	client := core.NewRedisClient(config)
	if err := client.Connect(); err != nil {
		return
	}
	defer client.Close()

	var relatedInstances []RedisInstance

	// Try cluster discovery
	if nodes, err := client.SendCommand("CLUSTER", "NODES"); err == nil {
		clusterNodes := parseClusterNodes(nodes)
		relatedInstances = append(relatedInstances, clusterNodes...)
	}

	// Try replication info
	if info, err := client.SendCommand("INFO", "replication"); err == nil {
		// Parse master/slave relationships
		lines := strings.Split(info, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "master_host:") {
				masterHost := strings.TrimSpace(strings.Split(line, ":")[1])
				if masterPort, err := client.SendCommand("CONFIG", "GET", "port"); err == nil {
					port, _ := strconv.Atoi(core.ParseConfigValue(masterPort, "port"))
					relatedInstances = append(relatedInstances, RedisInstance{
						Host: masterHost,
						Port: port,
						Role: "master",
					})
				}
			}
		}
	}

	topology[key] = relatedInstances

	// Recursively discover
	for _, related := range relatedInstances {
		discoverTopologyRecursive(related, topology, visited, depth+1, maxDepth)
	}
}

func displayTopology(topology map[string][]RedisInstance) {
	for root, connections := range topology {
		fmt.Printf("%s[ROOT]%s %s%s%s\n", core.ColorGreen, core.ColorNC, core.ColorYellow, root, core.ColorNC)
		for _, conn := range connections {
			fmt.Printf("  └── %s%s:%d%s (%s)\n", core.ColorYellow, conn.Host, conn.Port, core.ColorNC, conn.Role)
		}
		fmt.Println()
	}
}

func generateTopologyReport(topology map[string][]RedisInstance) {
	filename := fmt.Sprintf("redis_topology_%s.txt", time.Now().Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "Redis Network Topology Report\n")
	fmt.Fprintf(file, "Generated: %s\n\n", time.Now().Format(time.RFC3339))

	for root, connections := range topology {
		fmt.Fprintf(file, "ROOT: %s\n", root)
		for _, conn := range connections {
			fmt.Fprintf(file, "  Connected to: %s:%d (%s)\n", conn.Host, conn.Port, conn.Role)
		}
		fmt.Fprintf(file, "\n")
	}

	fmt.Printf("%s[+]%s Topology report saved to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}
