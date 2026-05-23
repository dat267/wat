//go:build linux

package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

func getDeviceID(path string) (uint64, error) {
	var st syscall.Stat_t
	if err := syscall.Stat(path, &st); err != nil {
		return 0, err
	}
	return st.Dev, nil
}

func getCPUTimes() (idle, total uint64, err error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 9 {
				return 0, 0, fmt.Errorf("invalid cpu fields")
			}
			var user, nice, system, idleVal, iowait, irq, softirq, steal uint64
			fmt.Sscanf(strings.Join(fields[1:9], " "), "%d %d %d %d %d %d %d %d",
				&user, &nice, &system, &idleVal, &iowait, &irq, &softirq, &steal)
			idle = idleVal + iowait
			total = user + nice + system + idleVal + iowait + irq + softirq + steal
			return idle, total, nil
		}
	}
	return 0, 0, fmt.Errorf("cpu line not found")
}
func getMemoryStats() (used, total, percent float64, err error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, err
	}
	var memTotal, memAvailable float64
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fmt.Sscanf(line, "MemTotal: %f", &memTotal)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fmt.Sscanf(line, "MemAvailable: %f", &memAvailable)
		}
	}
	if memTotal == 0 {
		return 0, 0, 0, fmt.Errorf("could not parse memory total")
	}
	memTotalGB := memTotal / (1024 * 1024)
	memAvailableGB := memAvailable / (1024 * 1024)
	memUsedGB := memTotalGB - memAvailableGB
	percent = (memUsedGB / memTotalGB) * 100.0
	return memUsedGB, memTotalGB, percent, nil
}
func getDiskStats() (used, total, percent float64, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs("/", &stat)
	if err != nil {
		return 0, 0, 0, err
	}
	totalBytes := stat.Blocks * uint64(stat.Bsize)
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	usedBytes := totalBytes - freeBytes
	totalGB := float64(totalBytes) / (1024 * 1024 * 1024)
	usedGB := float64(usedBytes) / (1024 * 1024 * 1024)
	percent = (usedGB / totalGB) * 100.0
	return usedGB, totalGB, percent, nil
}
func getUptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "Unknown"
	}
	var uptimeSec float64
	fmt.Sscanf(string(data), "%f", &uptimeSec)
	days := int(uptimeSec) / (24 * 3600)
	hours := (int(uptimeSec) % (24 * 3600)) / 3600
	minutes := (int(uptimeSec) % 3600) / 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
