package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const cliUsage = `Usage: wat <command> [options]

Commands:
  stats           Print system stats (CPU, RAM, Disk, Uptime)
  ports <host>    Scan common ports on a host
  disk  [path]    Rank files and folders by size (default: .)
  config          Manage secure configuration and API keys
  ai    <prompt>  Query local AI assistant (requires OpenAI/Anthropic key)
  logs            Print local application logs
  run   <script>  Run system script (docker, disk, network, services)

Run 'wat <command> -h' for command-specific options.
With no command, the GUI is launched.
`

func isCLIMode() bool {
	if len(os.Args) < 2 {
		return false
	}
	switch os.Args[1] {
	case "stats", "ports", "disk", "config", "ai", "logs", "run", "help", "-h", "--help":
		return true
	}
	return false
}

func runCLI() {
	if len(os.Args) < 2 {
		fmt.Print(cliUsage)
		os.Exit(0)
	}
	app := NewApp()
	switch os.Args[1] {
	case "stats":
		cmdStats(app)
	case "ports":
		cmdPorts(app)
	case "disk":
		cmdDisk(app)
	case "config":
		cmdConfig(app)
	case "ai":
		cmdAI(app)
	case "logs":
		cmdLogs(app)
	case "run":
		cmdRun(app)
	default:
		fmt.Print(cliUsage)
		os.Exit(0)
	}
}

func cmdStats(app *App) {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: wat stats\n\nPrint CPU, memory, disk, and uptime.\n")
	}
	fs.Parse(os.Args[2:])
	app.startup(context.Background())
	time.Sleep(500 * time.Millisecond)
	stats, err := app.GetSystemStats()
	if err != nil {
		fatalf("stats: %v\n", err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "METRIC\tVALUE")
	fmt.Fprintln(w, "------\t-----")
	fmt.Fprintf(w, "CPU\t%.1f%%\n", stats.CPUPercent)
	fmt.Fprintf(w, "Memory\t%.1f / %.1f GB (%.1f%%)\n", stats.MemoryUsed, stats.MemoryTotal, stats.MemoryPercent)
	fmt.Fprintf(w, "Disk\t%.1f / %.1f GB (%.1f%%)\n", stats.DiskUsed, stats.DiskTotal, stats.DiskPercent)
	fmt.Fprintf(w, "Uptime\t%s\n", stats.Uptime)
	w.Flush()
}

func cmdPorts(app *App) {
	fs := flag.NewFlagSet("ports", flag.ExitOnError)
	portsOpt := fs.String("p", "common", "Ports to scan (e.g. 'common', '80,443', '1-1024')")
	allOpt := fs.Bool("all", false, "Scan all 65535 ports")
	timeoutOpt := fs.Int("t", 500, "Timeout in milliseconds for each port attempt")
	concurrencyOpt := fs.Int("c", 1024, "Number of concurrent workers")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: wat ports [options] <host>\n\nScan ports on a host. Defaults to common ports.\n\nOptions:\n")
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[2:])
	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(1)
	}
	host := fs.Arg(0)
	portsSpec := *portsOpt
	if *allOpt {
		portsSpec = "all"
	}
	ports, err := parsePorts(portsSpec)
	if err != nil {
		fatalf("ports: %v\n", err)
	}
	fmt.Printf("Scanning %d ports on %s...\n", len(ports), host)
	open, err := app.ScanPorts(host, ports, *timeoutOpt, *concurrencyOpt)
	if err != nil {
		fatalf("ports: %v\n", err)
	}
	if len(open) == 0 {
		fmt.Println("No open ports found.")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PORT\tSTATUS")
	fmt.Fprintln(w, "----\t------")
	for _, p := range open {
		fmt.Fprintf(w, "%d\topen\n", p)
	}
	w.Flush()
}

func parsePorts(portsStr string) ([]int, error) {
	if portsStr == "all" {
		ports := make([]int, 65535)
		for i := 0; i < 65535; i++ {
			ports[i] = i + 1
		}
		return ports, nil
	}
	if portsStr == "common" {
		return []int{21, 22, 23, 25, 53, 80, 443, 3306, 5432, 6379, 8080, 8443, 27017}, nil
	}
	var ports []int
	seen := make(map[int]bool)
	parts := strings.Split(portsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}
			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start port: %s", rangeParts[0])
			}
			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end port: %s", rangeParts[1])
			}
			if start > end || start < 1 || end > 65535 {
				return nil, fmt.Errorf("port range out of bounds: %s", part)
			}
			for i := start; i <= end; i++ {
				if !seen[i] {
					seen[i] = true
					ports = append(ports, i)
				}
			}
		} else {
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("port out of bounds: %d", port)
			}
			if !seen[port] {
				seen[port] = true
				ports = append(ports, port)
			}
		}
	}
	return ports, nil
}

func cmdDisk(app *App) {
	fs := flag.NewFlagSet("disk", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: wat disk [path]\n\nRank files and folders by size. Defaults to current directory.\n")
	}
	fs.Parse(os.Args[2:])
	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}
	fmt.Printf("Ranking %s...\n\n", path)
	entries, err := app.RankDirectory(path)
	if err != nil {
		fatalf("disk: %v\n", err)
	}
	if len(entries) == 0 {
		fmt.Println("Directory is empty.")
		return
	}
	var total int64
	for _, e := range entries {
		total += e.Size
	}
	for _, e := range entries {
		name := e.Name
		if e.IsDir {
			name += "/"
		}
		pct := 0.0
		if total > 0 {
			pct = float64(e.Size) / float64(total) * 100
		}
		fmt.Printf("  %8s  %5.1f%%  %s\n", humanSize(e.Size), pct, name)
	}
	fmt.Printf("\n  %8s  total\n", humanSize(total))
}

func humanSize(b int64) string {
	const (
		GB = 1024 * 1024 * 1024
		MB = 1024 * 1024
		KB = 1024
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.1f GB", float64(b)/GB)
	case b >= MB:
		return fmt.Sprintf("%.1f MB", float64(b)/MB)
	case b >= KB:
		return fmt.Sprintf("%.1f KB", float64(b)/KB)
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format, args...)
	os.Exit(1)
}

func cmdConfig(app *App) {
	fs := flag.NewFlagSet("config", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: wat config <action> [arguments]\n\nActions:\n  set <category> <key> <value>   Set a configuration value (category: 'apiKeys' or 'settings')\n  get <category> <key>           Get a configuration value\n  delete <category> <key>        Delete a configuration value\n  reset                          Reset all configuration values to defaults\n  show                           Show all configuration keys with metadata\n  path                           Print the path to the configuration file\n")
	}
	fs.Parse(os.Args[2:])
	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(1)
	}
	switch fs.Arg(0) {
	case "set":
		if fs.NArg() < 4 {
			fatalf("usage: wat config set <category> <key> <value>\n")
		}
		category := fs.Arg(1)
		if category != "apiKeys" && category != "settings" {
			fatalf("invalid category: %s (must be 'apiKeys' or 'settings')\n", category)
		}
		key := fs.Arg(2)
		val := fs.Arg(3)
		if IsKeyMandatory(key) && strings.TrimSpace(val) == "" {
			fatalf("key '%s' is mandatory and cannot be empty\n", key)
		}
		if err := app.SetConfigValue(category, key, val); err != nil {
			fatalf("failed to set config: %v\n", err)
		}
		fmt.Printf("Successfully set %s.%s\n", category, key)
	case "get":
		if fs.NArg() < 3 {
			fatalf("usage: wat config get <category> <key>\n")
		}
		category := fs.Arg(1)
		if category != "apiKeys" && category != "settings" {
			fatalf("invalid category: %s (must be 'apiKeys' or 'settings')\n", category)
		}
		key := fs.Arg(2)
		cfg, err := app.GetConfig()
		if err != nil {
			fatalf("failed to load config: %v\n", err)
		}
		var val string
		var found bool
		if category == "apiKeys" {
			val, found = cfg.APIKeys[key]
		} else {
			val, found = cfg.Settings[key]
		}
		if !found {
			fatalf("key '%s' not found in category '%s'\n", key, category)
		}
		fmt.Println(val)
	case "delete":
		if fs.NArg() < 3 {
			fatalf("usage: wat config delete <category> <key>\n")
		}
		category := fs.Arg(1)
		if category != "apiKeys" && category != "settings" {
			fatalf("invalid category: %s (must be 'apiKeys' or 'settings')\n", category)
		}
		key := fs.Arg(2)
		if err := app.DeleteConfigValue(category, key); err != nil {
			fatalf("failed to delete config: %v\n", err)
		}
		if IsKeyMandatory(key) {
			fmt.Printf("Successfully reset mandatory config %s.%s to default: '%s'\n", category, key, GetDefaultValue(key))
		} else {
			fmt.Printf("Successfully deleted config %s.%s\n", category, key)
		}
	case "reset":
		if err := app.ResetConfig(); err != nil {
			fatalf("failed to reset config: %v\n", err)
		}
		fmt.Println("Successfully reset all configuration values to defaults.")
	case "show":
		cfg, err := app.GetConfig()
		if err != nil {
			fatalf("failed to load config: %v\n", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CATEGORY\tKEY\tVALUE\tSECRET\tMANDATORY")
		fmt.Fprintln(w, "--------\t---\t-----\t------\t---------")
		for k, v := range cfg.APIKeys {
			secTag := "yes"
			mandTag := "no"
			if IsKeyMandatory(k) {
				mandTag = "yes"
			}
			valTag := "[hidden]"
			if !IsKeySecret("apiKeys", k) {
				secTag = "no"
				valTag = v
			}
			fmt.Fprintf(w, "apiKeys\t%s\t%s\t%s\t%s\n", k, valTag, secTag, mandTag)
		}
		for k, v := range cfg.Settings {
			secTag := "no"
			mandTag := "no"
			if IsKeyMandatory(k) {
				mandTag = "yes"
			}
			valTag := v
			if IsKeySecret("settings", k) {
				secTag = "yes"
				valTag = "[hidden]"
			}
			fmt.Fprintf(w, "settings\t%s\t%s\t%s\t%s\n", k, valTag, secTag, mandTag)
		}
		w.Flush()
	case "path":
		path, err := getConfigPath()
		if err != nil {
			fatalf("failed to get config path: %v\n", err)
		}
		fmt.Println(path)
	default:
		fs.Usage()
		os.Exit(1)
	}
}

func cmdAI(app *App) {
	if len(os.Args) < 3 {
		fmt.Println("usage: wat ai <prompt>")
		os.Exit(1)
	}
	prompt := strings.Join(os.Args[2:], " ")
	cfg, err := loadConfig()
	if err != nil {
		fatalf("failed to load configuration: %v\n", err)
	}
	provider := "anthropic"
	if cfg.APIKeys["anthropic"] == "" && cfg.APIKeys["openai"] != "" {
		provider = "openai"
	}
	fmt.Printf("Querying %s...\n", provider)
	app.startup(context.Background())
	res, err := app.AskAI(provider, "", "You are a helpful shell assistant. Keep your answer brief, direct and informative.", prompt)
	if err != nil {
		fatalf("AI error: %v\n", err)
	}
	fmt.Println(res)
}

func cmdLogs(app *App) {
	path, err := getConfigPath()
	if err != nil {
		fatalf("failed to get config path: %v\n", err)
	}
	logPath := filepath.Join(filepath.Dir(path), "app.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		fatalf("failed to read logs: %v\n", err)
	}
	fmt.Print(string(data))
}

func cmdRun(app *App) {
	if len(os.Args) < 3 {
		fmt.Println("usage: wat run <script> (docker | disk | network | services)")
		os.Exit(1)
	}
	script := os.Args[2]
	app.startup(context.Background())
	out, err := app.ExecuteScript(script)
	if err != nil {
		fatalf("script failed: %v\n", err)
	}
	fmt.Print(out)
}
