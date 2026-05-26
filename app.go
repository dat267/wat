package main

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type SystemStats struct {
	CPUPercent    float64 `json:"cpuPercent"`
	MemoryPercent float64 `json:"memoryPercent"`
	MemoryUsed    float64 `json:"memoryUsed"`
	MemoryTotal   float64 `json:"memoryTotal"`
	DiskPercent   float64 `json:"diskPercent"`
	DiskUsed      float64 `json:"diskUsed"`
	DiskTotal     float64 `json:"diskTotal"`
	Uptime        string  `json:"uptime"`
}
type App struct {
	ctx          context.Context
	lastCPUIDle  uint64
	lastCPUTotal uint64
	isGUI        bool
}

func NewApp() *App {
	return &App{}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	initStore()
	initLogging(a)
	startScheduler(a)
	logInfo("wat desktop application starting up", nil)
	if idle, total, err := getCPUTimes(); err == nil {
		a.lastCPUIDle = idle
		a.lastCPUTotal = total
	}
}
func (a *App) GetSystemStats() (SystemStats, error) {
	var stats SystemStats
	currIdle, currTotal, err := getCPUTimes()
	if err == nil {
		idleDiff := currIdle - a.lastCPUIDle
		totalDiff := currTotal - a.lastCPUTotal
		if totalDiff > 0 {
			stats.CPUPercent = (1.0 - float64(idleDiff)/float64(totalDiff)) * 100.0
		}
		a.lastCPUIDle = currIdle
		a.lastCPUTotal = currTotal
	}
	if mu, mt, mp, err := getMemoryStats(); err == nil {
		stats.MemoryUsed = mu
		stats.MemoryTotal = mt
		stats.MemoryPercent = mp
	}
	if du, dt, dp, err := getDiskStats(); err == nil {
		stats.DiskUsed = du
		stats.DiskTotal = dt
		stats.DiskPercent = dp
	}
	stats.Uptime = getUptime()
	return stats, nil
}
func (a *App) ExecuteScript(name string) (string, error) {
	isWin := runtime.GOOS == "windows"
	var cmd *exec.Cmd
	switch name {
	case "docker":
		if isWin {
			cmd = exec.Command("powershell", "-Command", "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'")
		} else {
			cmd = exec.Command("bash", "-c", "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' 2>&1")
		}
	case "disk":
		if isWin {
			cmd = exec.Command("powershell", "-Command", "Get-PSDrive C | Select-Object @{Name='Drive';Expression={$_.Name}}, @{Name='Used(GB)';Expression={[math]::Round($_.Used/1GB,2)}}, @{Name='Free(GB)';Expression={[math]::Round($_.Free/1GB,2)}} | Format-Table -Autosize")
		} else {
			cmd = exec.Command("bash", "-c", "df -h / 2>&1")
		}
	case "network":
		if isWin {
			cmd = exec.Command("powershell", "-Command", "Get-NetIPAddress -AddressFamily IPv4 | Select-Object InterfaceAlias, IPAddress | Format-Table -Autosize")
		} else {
			cmd = exec.Command("bash", "-c", "ip -brief address 2>&1")
		}
	case "services":
		if isWin {
			cmd = exec.Command("powershell", "-Command", "Get-Service | Where-Object {$_.Status -eq 'Running'} | Select-Object -First 15 DisplayName, Status | Format-Table -Autosize")
		} else {
			cmd = exec.Command("bash", "-c", "systemctl list-units --type=service --state=running --no-legend | head -n 15 2>&1")
		}
	default:
		return "", fmt.Errorf("unknown script: %s", name)
	}
	out, _ := cmd.CombinedOutput()
	return string(out), nil
}
func (a *App) PingHost(host string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "3", host)
	} else {
		cmd = exec.Command("ping", "-c", "3", host)
	}
	out, _ := cmd.CombinedOutput()
	return string(out), nil
}
func (a *App) ScanPorts(host string, ports []int, timeoutMs int, concurrency int) ([]int, error) {
	if timeoutMs <= 0 {
		timeoutMs = 500
	}
	if concurrency <= 0 {
		concurrency = 1024
	}
	var openPorts []int
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	timeout := time.Duration(timeoutMs) * time.Millisecond
	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			address := net.JoinHostPort(host, fmt.Sprintf("%d", p))
			conn, err := net.DialTimeout("tcp", address, timeout)
			if err == nil {
				conn.Close()
				mu.Lock()
				openPorts = append(openPorts, p)
				mu.Unlock()
			}
		}(port)
	}
	wg.Wait()
	sort.Ints(openPorts)
	return openPorts, nil
}

type DiskEntry struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"isDir"`
}

func (a *App) RankDirectory(dirPath string) ([]DiskEntry, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var results []DiskEntry
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		size, _ := dirSize(fullPath)
		results = append(results, DiskEntry{
			Name:  entry.Name(),
			Path:  fullPath,
			Size:  size,
			IsDir: entry.IsDir(),
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Size > results[j].Size
	})
	return results, nil
}

func dirSize(path string) (int64, error) {
	parentDev, _ := getDeviceID(filepath.Dir(path))
	selfDev, err := getDeviceID(path)
	if err == nil && parentDev != 0 && selfDev != parentDev {
		return 0, nil
	}
	var total int64
	err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && p != path {
			dev, err := getDeviceID(p)
			if err == nil && selfDev != 0 && dev != selfDev {
				return filepath.SkipDir
			}
		}
		if !d.IsDir() {
			if info, err := d.Info(); err == nil {
				total += info.Size()
			}
		}
		return nil
	})
	return total, err
}

func (a *App) BroadcastLog(msg string) {
	if a.isGUI && a.ctx != nil {
		wailsRuntime.EventsEmit(a.ctx, "log", msg)
	}
}
