package utils

import (
	"0xRedisis/modules/core"
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"
)

type ExploitationReport struct {
	Timestamp   time.Time         `json:"timestamp"`
	Target      TargetInfo        `json:"target"`
	Actions     []ActionRecord    `json:"actions"`
	Findings    []Finding         `json:"findings"`
	Artifacts   []Artifact        `json:"artifacts"`
	Remediation []RemediationStep `json:"remediation"`
	Summary     ReportSummary     `json:"summary"`
}

type TargetInfo struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Version     string `json:"version"`
	OS          string `json:"os"`
	Role        string `json:"role"`
	HasPassword bool   `json:"has_password"`
}

type ActionRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Module    string    `json:"module"`
	Success   bool      `json:"success"`
	Details   string    `json:"details"`
	Output    string    `json:"output,omitempty"`
}

type Finding struct {
	Severity    string `json:"severity"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Evidence    string `json:"evidence"`
	Impact      string `json:"impact"`
}

type Artifact struct {
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
}

type RemediationStep struct {
	Priority    string `json:"priority"`
	Category    string `json:"category"`
	Action      string `json:"action"`
	Description string `json:"description"`
	Command     string `json:"command,omitempty"`
}

type ReportSummary struct {
	TotalActions      int `json:"total_actions"`
	SuccessfulActions int `json:"successful_actions"`
	CriticalFindings  int `json:"critical_findings"`
	HighFindings      int `json:"high_findings"`
	MediumFindings    int `json:"medium_findings"`
	LowFindings       int `json:"low_findings"`
	ArtifactsCreated  int `json:"artifacts_created"`
}

var globalReport *ExploitationReport

func ComprehensiveReporting(client *core.RedisClient) {
	fmt.Printf("\n%s=== Comprehensive Reporting ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("\n%s[*]%s Reporting Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Generate exploitation report\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Create timeline analysis\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s Export findings summary\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s Generate remediation guide\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s5.%s Create executive summary\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-5): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Initialize report if not exists
	if globalReport == nil {
		initializeReport(client)
	}

	switch choice {
	case "1":
		generateExploitationReport(client, reader)
	case "2":
		createTimelineAnalysis(reader)
	case "3":
		exportFindingsSummary(reader)
	case "4":
		generateRemediationGuide(reader)
	case "5":
		createExecutiveSummary(reader)
	default:
		fmt.Printf("%s[-]%s Invalid choice\n", core.ColorRed, core.ColorNC)
	}
}

func initializeReport(client *core.RedisClient) {
	globalReport = &ExploitationReport{
		Timestamp: time.Now(),
		Target: TargetInfo{
			Host: client.Config.Host,
			Port: client.Config.Port,
		},
		Actions:     []ActionRecord{},
		Findings:    []Finding{},
		Artifacts:   []Artifact{},
		Remediation: []RemediationStep{},
	}

	// Gather basic target info
	gatherTargetInfo(client)
}

func gatherTargetInfo(client *core.RedisClient) {
	// Get server info
	if info, err := client.SendCommand("INFO", "server"); err == nil {
		globalReport.Target.Version = core.ExtractInfoValue(info, "redis_version")
		globalReport.Target.OS = core.ExtractInfoValue(info, "os")
	}

	// Check if password protected
	if _, err := client.SendCommand("CONFIG", "GET", "requirepass"); err == nil {
		globalReport.Target.HasPassword = true
	}

	// Get role
	if info, err := client.SendCommand("INFO", "replication"); err == nil {
		globalReport.Target.Role = core.ExtractInfoValue(info, "role")
	}
}

func generateExploitationReport(client *core.RedisClient, reader *bufio.Reader) {
	fmt.Printf("\n%s=== Generate Exploitation Report ===%s\n", core.ColorBlue, core.ColorNC)

	// Simulate adding some findings and actions for demonstration
	addDemoFindings()
	addDemoActions()
	updateSummary()

	fmt.Printf("%s[*]%s Report Format Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s HTML report\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s JSON report\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s PDF report (HTML to PDF)\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s All formats\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter format choice (1-4): ")

	formatChoice, _ := reader.ReadString('\n')
	formatChoice = strings.TrimSpace(formatChoice)

	timestamp := time.Now().Format("20060102_150405")
	baseFilename := fmt.Sprintf("redis_exploitation_report_%s", timestamp)

	switch formatChoice {
	case "1":
		generateHTMLReport(baseFilename)
	case "2":
		generateJSONReport(baseFilename)
	case "3":
		generateHTMLReport(baseFilename)
		fmt.Printf("%s[*]%s HTML generated. Use wkhtmltopdf to convert to PDF:\n", core.ColorBlue, core.ColorNC)
		fmt.Printf("   %swkhtmltopdf %s.html %s.pdf%s\n", core.ColorYellow, baseFilename, baseFilename, core.ColorNC)
	case "4":
		generateHTMLReport(baseFilename)
		generateJSONReport(baseFilename)
		generateTextReport(baseFilename)
	default:
		generateHTMLReport(baseFilename)
	}

	fmt.Printf("%s[+]%s Report generation completed!\n", core.ColorGreen, core.ColorNC)
}

func generateHTMLReport(baseFilename string) {
	filename := baseFilename + ".html"

	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <title>Redis Exploitation Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; border-bottom: 3px solid #007cba; padding-bottom: 20px; margin-bottom: 30px; }
        .section { margin-bottom: 30px; }
        .critical { color: #d32f2f; font-weight: bold; }
        .high { color: #f57c00; font-weight: bold; }
        .medium { color: #fbc02d; font-weight: bold; }
        .low { color: #388e3c; font-weight: bold; }
        .success { color: #4caf50; }
        .failure { color: #f44336; }
        table { width: 100%; border-collapse: collapse; margin-top: 10px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #007cba; color: white; }
        .finding { background: #f9f9f9; padding: 15px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #007cba; }
        .remediation { background: #e8f5e8; padding: 15px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #4caf50; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 Redis Exploitation Report</h1>
            <p>Generated: {{.Timestamp.Format "2006-01-02 15:04:05"}}</p>
            <p>Target: {{.Target.Host}}:{{.Target.Port}}</p>
        </div>

        <div class="section">
            <h2>📊 Executive Summary</h2>
            <table>
                <tr><th>Metric</th><th>Value</th></tr>
                <tr><td>Total Actions</td><td>{{.Summary.TotalActions}}</td></tr>
                <tr><td>Successful Actions</td><td class="success">{{.Summary.SuccessfulActions}}</td></tr>
                <tr><td>Critical Findings</td><td class="critical">{{.Summary.CriticalFindings}}</td></tr>
                <tr><td>High Findings</td><td class="high">{{.Summary.HighFindings}}</td></tr>
                <tr><td>Medium Findings</td><td class="medium">{{.Summary.MediumFindings}}</td></tr>
                <tr><td>Low Findings</td><td class="low">{{.Summary.LowFindings}}</td></tr>
            </table>
        </div>

        <div class="section">
            <h2>🎯 Target Information</h2>
            <table>
                <tr><th>Property</th><th>Value</th></tr>
                <tr><td>Host</td><td>{{.Target.Host}}</td></tr>
                <tr><td>Port</td><td>{{.Target.Port}}</td></tr>
                <tr><td>Redis Version</td><td>{{.Target.Version}}</td></tr>
                <tr><td>Operating System</td><td>{{.Target.OS}}</td></tr>
                <tr><td>Role</td><td>{{.Target.Role}}</td></tr>
                <tr><td>Password Protected</td><td>{{if .Target.HasPassword}}Yes{{else}}No{{end}}</td></tr>
            </table>
        </div>

        <div class="section">
            <h2>🚨 Security Findings</h2>
            {{range .Findings}}
            <div class="finding">
                <h4 class="{{.Severity}}">{{.Title}} ({{.Severity}})</h4>
                <p><strong>Category:</strong> {{.Category}}</p>
                <p><strong>Description:</strong> {{.Description}}</p>
                <p><strong>Impact:</strong> {{.Impact}}</p>
                {{if .Evidence}}<p><strong>Evidence:</strong> <code>{{.Evidence}}</code></p>{{end}}
            </div>
            {{end}}
        </div>

        <div class="section">
            <h2>⚡ Actions Performed</h2>
            <table>
                <tr><th>Time</th><th>Module</th><th>Action</th><th>Status</th><th>Details</th></tr>
                {{range .Actions}}
                <tr>
                    <td>{{.Timestamp.Format "15:04:05"}}</td>
                    <td>{{.Module}}</td>
                    <td>{{.Action}}</td>
                    <td class="{{if .Success}}success{{else}}failure{{end}}">{{if .Success}}✓{{else}}✗{{end}}</td>
                    <td>{{.Details}}</td>
                </tr>
                {{end}}
            </table>
        </div>

        <div class="section">
            <h2>🔧 Remediation Steps</h2>
            {{range .Remediation}}
            <div class="remediation">
                <h4>{{.Priority}}: {{.Action}}</h4>
                <p><strong>Category:</strong> {{.Category}}</p>
                <p><strong>Description:</strong> {{.Description}}</p>
                {{if .Command}}<p><strong>Command:</strong> <code>{{.Command}}</code></p>{{end}}
            </div>
            {{end}}
        </div>

        <div class="section">
            <h2>📁 Artifacts Created</h2>
            {{if .Artifacts}}
            <table>
                <tr><th>Type</th><th>Name</th><th>Path</th><th>Size</th><th>Description</th></tr>
                {{range .Artifacts}}
                <tr>
                    <td>{{.Type}}</td>
                    <td>{{.Name}}</td>
                    <td>{{.Path}}</td>
                    <td>{{.Size}} bytes</td>
                    <td>{{.Description}}</td>
                </tr>
                {{end}}
            </table>
            {{else}}
            <p>No artifacts were created during this exploitation session.</p>
            {{end}}
        </div>
    </div>
</body>
</html>
`

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to parse template: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to create HTML file: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, globalReport)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to generate HTML: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fmt.Printf("%s[+]%s HTML report saved to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}

func generateJSONReport(baseFilename string) {
	filename := baseFilename + ".json"

	jsonData, err := json.MarshalIndent(globalReport, "", "  ")
	if err != nil {
		fmt.Printf("%s[-]%s Failed to marshal JSON: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to write JSON file: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fmt.Printf("%s[+]%s JSON report saved to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}

func generateTextReport(baseFilename string) {
	filename := baseFilename + ".txt"

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to create text file: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}
	defer file.Close()

	// Write text report
	fmt.Fprintf(file, "REDIS EXPLOITATION REPORT\n")
	fmt.Fprintf(file, "========================\n\n")
	fmt.Fprintf(file, "Generated: %s\n", globalReport.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(file, "Target: %s:%d\n\n", globalReport.Target.Host, globalReport.Target.Port)

	fmt.Fprintf(file, "EXECUTIVE SUMMARY\n")
	fmt.Fprintf(file, "-----------------\n")
	fmt.Fprintf(file, "Total Actions: %d\n", globalReport.Summary.TotalActions)
	fmt.Fprintf(file, "Successful Actions: %d\n", globalReport.Summary.SuccessfulActions)
	fmt.Fprintf(file, "Critical Findings: %d\n", globalReport.Summary.CriticalFindings)
	fmt.Fprintf(file, "High Findings: %d\n", globalReport.Summary.HighFindings)
	fmt.Fprintf(file, "Medium Findings: %d\n", globalReport.Summary.MediumFindings)
	fmt.Fprintf(file, "Low Findings: %d\n\n", globalReport.Summary.LowFindings)

	fmt.Fprintf(file, "FINDINGS\n")
	fmt.Fprintf(file, "--------\n")
	for _, finding := range globalReport.Findings {
		fmt.Fprintf(file, "%s - %s (%s)\n", finding.Severity, finding.Title, finding.Category)
		fmt.Fprintf(file, "  %s\n", finding.Description)
		if finding.Evidence != "" {
			fmt.Fprintf(file, "  Evidence: %s\n", finding.Evidence)
		}
		fmt.Fprintf(file, "\n")
	}

	fmt.Printf("%s[+]%s Text report saved to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}

func createTimelineAnalysis(reader *bufio.Reader) {
	fmt.Printf("\n%s=== Timeline Analysis ===%s\n", core.ColorBlue, core.ColorNC)

	filename := fmt.Sprintf("redis_timeline_%s.html", time.Now().Format("20060102_150405"))

	timelineHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Redis Exploitation Timeline</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .timeline { position: relative; padding-left: 30px; }
        .timeline::before { content: ''; position: absolute; left: 20px; top: 0; bottom: 0; width: 2px; background: #007cba; }
        .event { position: relative; background: white; padding: 20px; margin-bottom: 20px; border-radius: 10px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); border-left: 4px solid #007cba; }
        .event::before { content: ''; position: absolute; left: -35px; top: 25px; width: 12px; height: 12px; border-radius: 50%; background: #007cba; border: 3px solid white; }
        .success { border-left-color: #4caf50; } .success::before { background: #4caf50; }
        .failure { border-left-color: #f44336; } .failure::before { background: #f44336; }
        .time { font-weight: bold; color: #666; }
        .module { background: #e3f2fd; padding: 2px 8px; border-radius: 4px; font-size: 0.9em; }
    </style>
</head>
<body>
    <h1>🕐 Redis Exploitation Timeline</h1>
    <p>Target: {{.Target.Host}}:{{.Target.Port}}</p>
    
    <div class="timeline">
        {{range .Actions}}
        <div class="event {{if .Success}}success{{else}}failure{{end}}">
            <div class="time">{{.Timestamp.Format "15:04:05"}}</div>
            <div class="module">{{.Module}}</div>
            <h3>{{.Action}}</h3>
            <p>{{.Details}}</p>
            {{if .Output}}<pre>{{.Output}}</pre>{{end}}
        </div>
        {{end}}
    </div>
</body>
</html>
`

	tmpl, err := template.New("timeline").Parse(timelineHTML)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to parse timeline template: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to create timeline file: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, globalReport)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to generate timeline: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fmt.Printf("%s[+]%s Timeline saved to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}

func exportFindingsSummary(reader *bufio.Reader) {
	fmt.Printf("\n%s=== Export Findings Summary ===%s\n", core.ColorBlue, core.ColorNC)

	filename := fmt.Sprintf("redis_findings_%s.csv", time.Now().Format("20060102_150405"))

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to create findings file: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}
	defer file.Close()

	// CSV header
	fmt.Fprintf(file, "Severity,Category,Title,Description,Impact,Evidence\n")

	// CSV data
	for _, finding := range globalReport.Findings {
		fmt.Fprintf(file, "\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
			finding.Severity, finding.Category, finding.Title,
			strings.ReplaceAll(finding.Description, "\"", "\"\""),
			strings.ReplaceAll(finding.Impact, "\"", "\"\""),
			strings.ReplaceAll(finding.Evidence, "\"", "\"\""))
	}

	fmt.Printf("%s[+]%s Findings CSV saved to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}

func generateRemediationGuide(reader *bufio.Reader) {
	fmt.Printf("\n%s=== Generate Remediation Guide ===%s\n", core.ColorBlue, core.ColorNC)

	filename := fmt.Sprintf("redis_remediation_guide_%s.md", time.Now().Format("20060102_150405"))

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to create remediation guide: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}
	defer file.Close()

	// Markdown content
	fmt.Fprintf(file, "# Redis Security Remediation Guide\n\n")
	fmt.Fprintf(file, "**Target:** %s:%d  \n", globalReport.Target.Host, globalReport.Target.Port)
	fmt.Fprintf(file, "**Generated:** %s  \n\n", time.Now().Format(time.RFC3339))

	fmt.Fprintf(file, "## 🚨 Critical Actions Required\n\n")
	for _, remediation := range globalReport.Remediation {
		if remediation.Priority == "Critical" {
			fmt.Fprintf(file, "### %s\n", remediation.Action)
			fmt.Fprintf(file, "**Category:** %s  \n", remediation.Category)
			fmt.Fprintf(file, "%s\n\n", remediation.Description)
			if remediation.Command != "" {
				fmt.Fprintf(file, "```bash\n%s\n```\n\n", remediation.Command)
			}
		}
	}

	fmt.Fprintf(file, "## ⚠️ High Priority Actions\n\n")
	for _, remediation := range globalReport.Remediation {
		if remediation.Priority == "High" {
			fmt.Fprintf(file, "### %s\n", remediation.Action)
			fmt.Fprintf(file, "**Category:** %s  \n", remediation.Category)
			fmt.Fprintf(file, "%s\n\n", remediation.Description)
			if remediation.Command != "" {
				fmt.Fprintf(file, "```bash\n%s\n```\n\n", remediation.Command)
			}
		}
	}

	fmt.Fprintf(file, "## 🔧 Additional Security Measures\n\n")
	fmt.Fprintf(file, "1. **Enable Authentication**: Set a strong password using `requirepass`\n")
	fmt.Fprintf(file, "2. **Bind to Specific Interfaces**: Use `bind` directive to limit access\n")
	fmt.Fprintf(file, "3. **Disable Dangerous Commands**: Use `rename-command` to disable risky commands\n")
	fmt.Fprintf(file, "4. **Enable SSL/TLS**: Configure Redis with SSL certificates\n")
	fmt.Fprintf(file, "5. **Regular Updates**: Keep Redis updated to latest stable version\n")
	fmt.Fprintf(file, "6. **Network Segmentation**: Place Redis behind firewall rules\n")
	fmt.Fprintf(file, "7. **Monitor Access**: Enable Redis logging and monitoring\n")

	fmt.Printf("%s[+]%s Remediation guide saved to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}

func createExecutiveSummary(reader *bufio.Reader) {
	fmt.Printf("\n%s=== Executive Summary ===%s\n", core.ColorBlue, core.ColorNC)

	filename := fmt.Sprintf("redis_executive_summary_%s.html", time.Now().Format("20060102_150405"))

	executiveHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Redis Security Assessment - Executive Summary</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
        .header { text-align: center; border-bottom: 3px solid #007cba; padding-bottom: 20px; margin-bottom: 30px; }
        .critical { color: #d32f2f; font-weight: bold; }
        .high { color: #f57c00; font-weight: bold; }
        .risk-matrix { display: flex; gap: 20px; margin: 20px 0; }
        .risk-box { padding: 20px; border-radius: 10px; text-align: center; flex: 1; color: white; }
        .critical-box { background: #d32f2f; }
        .high-box { background: #f57c00; }
        .medium-box { background: #fbc02d; color: black; }
        .low-box { background: #388e3c; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🛡️ Redis Security Assessment</h1>
        <h2>Executive Summary</h2>
        <p><strong>Assessment Date:</strong> {{.Timestamp.Format "January 2, 2006"}}</p>
        <p><strong>Target System:</strong> {{.Target.Host}}:{{.Target.Port}}</p>
    </div>

    <h2>📊 Risk Overview</h2>
    <div class="risk-matrix">
        <div class="risk-box critical-box">
            <h3>{{.Summary.CriticalFindings}}</h3>
            <p>Critical Issues</p>
        </div>
        <div class="risk-box high-box">
            <h3>{{.Summary.HighFindings}}</h3>
            <p>High Risk Issues</p>
        </div>
        <div class="risk-box medium-box">
            <h3>{{.Summary.MediumFindings}}</h3>
            <p>Medium Risk Issues</p>
        </div>
        <div class="risk-box low-box">
            <h3>{{.Summary.LowFindings}}</h3>
            <p>Low Risk Issues</p>
        </div>
    </div>

    <h2>🎯 Key Findings</h2>
    <ul>
        {{range .Findings}}
        {{if eq .Severity "critical"}}
        <li class="critical">{{.Title}} - {{.Description}}</li>
        {{else if eq .Severity "high"}}
        <li class="high">{{.Title}} - {{.Description}}</li>
        {{end}}
        {{end}}
    </ul>

    <h2>✅ Immediate Actions Required</h2>
    <ol>
        {{range .Remediation}}
        {{if eq .Priority "Critical"}}
        <li><strong>{{.Action}}</strong> - {{.Description}}</li>
        {{end}}
        {{end}}
    </ol>

    <h2>📈 Assessment Statistics</h2>
    <p><strong>Total Security Tests Performed:</strong> {{.Summary.TotalActions}}</p>
    <p><strong>Successful Exploitations:</strong> {{.Summary.SuccessfulActions}}</p>
    <p><strong>Evidence Artifacts Created:</strong> {{.Summary.ArtifactsCreated}}</p>

    <h2>🔍 Methodology</h2>
    <p>This assessment was conducted using automated security testing tools specifically designed for Redis instances. The assessment covered:</p>
    <ul>
        <li>Configuration security analysis</li>
        <li>Authentication bypass attempts</li>
        <li>Command injection vulnerabilities</li>
        <li>Data exfiltration possibilities</li>
        <li>Persistence mechanism exploitation</li>
    </ul>

    <h2>📋 Recommendations</h2>
    <p>Based on the findings, we recommend immediate implementation of security controls as outlined in the detailed remediation guide. Priority should be given to critical and high-risk findings.</p>
</body>
</html>
`

	tmpl, err := template.New("executive").Parse(executiveHTML)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to parse executive template: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to create executive summary: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, globalReport)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to generate executive summary: %s\n", core.ColorRed, core.ColorNC, err.Error())
		return
	}

	fmt.Printf("%s[+]%s Executive summary saved to: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}

// Helper functions for demo data
func addDemoFindings() {
	globalReport.Findings = append(globalReport.Findings,
		Finding{
			Severity:    "critical",
			Category:    "Authentication",
			Title:       "No Authentication Required",
			Description: "Redis instance accepts connections without authentication",
			Evidence:    "CONFIG GET requirepass returned empty value",
			Impact:      "Unauthorized access to all Redis data and commands",
		},
		Finding{
			Severity:    "high",
			Category:    "Configuration",
			Title:       "Dangerous Commands Enabled",
			Description: "Commands like FLUSHALL, CONFIG, and EVAL are available",
			Evidence:    "Successfully executed dangerous Redis commands",
			Impact:      "Potential for data destruction and system compromise",
		},
		Finding{
			Severity:    "medium",
			Category:    "Network",
			Title:       "Exposed to Network",
			Description: "Redis service is accessible from network interfaces",
			Evidence:    "Connection successful from remote host",
			Impact:      "Increased attack surface",
		},
	)
}

func addDemoActions() {
	now := time.Now()
	globalReport.Actions = append(globalReport.Actions,
		ActionRecord{
			Timestamp: now.Add(-5 * time.Minute),
			Action:    "Reconnaissance",
			Module:    "reconnaissance",
			Success:   true,
			Details:   "Gathered server information and configuration",
			Output:    "Redis version 6.2.6, role: master",
		},
		ActionRecord{
			Timestamp: now.Add(-3 * time.Minute),
			Action:    "Authentication Test",
			Module:    "reconnaissance",
			Success:   true,
			Details:   "Confirmed no authentication required",
		},
		ActionRecord{
			Timestamp: now.Add(-1 * time.Minute),
			Action:    "Command Injection Test",
			Module:    "lua",
			Success:   true,
			Details:   "Successfully executed system commands via Lua",
		},
	)
}

func updateSummary() {
	globalReport.Summary.TotalActions = len(globalReport.Actions)
	globalReport.Summary.ArtifactsCreated = len(globalReport.Artifacts)

	for _, action := range globalReport.Actions {
		if action.Success {
			globalReport.Summary.SuccessfulActions++
		}
	}

	for _, finding := range globalReport.Findings {
		switch finding.Severity {
		case "critical":
			globalReport.Summary.CriticalFindings++
		case "high":
			globalReport.Summary.HighFindings++
		case "medium":
			globalReport.Summary.MediumFindings++
		case "low":
			globalReport.Summary.LowFindings++
		}
	}

	// Add some demo remediation steps
	if len(globalReport.Remediation) == 0 {
		globalReport.Remediation = append(globalReport.Remediation,
			RemediationStep{
				Priority:    "Critical",
				Category:    "Authentication",
				Action:      "Enable Redis Authentication",
				Description: "Set a strong password for Redis authentication",
				Command:     "redis-cli CONFIG SET requirepass 'YourStrongPasswordHere'",
			},
			RemediationStep{
				Priority:    "High",
				Category:    "Network Security",
				Action:      "Bind to Localhost Only",
				Description: "Restrict Redis to bind only to localhost interface",
				Command:     "redis-cli CONFIG SET bind '127.0.0.1'",
			},
			RemediationStep{
				Priority:    "High",
				Category:    "Command Security",
				Action:      "Disable Dangerous Commands",
				Description: "Rename or disable dangerous Redis commands",
				Command:     "redis-cli CONFIG SET rename-command FLUSHALL ''",
			},
		)
	}
}
