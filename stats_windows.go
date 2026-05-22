//go:build windows

package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

var (
	modkernel32              = windows.NewLazySystemDLL("kernel32.dll")
	procGetSystemTimes       = modkernel32.NewProc("GetSystemTimes")
	procGlobalMemoryStatusEx = modkernel32.NewProc("GlobalMemoryStatusEx")
	procGetTickCount64       = modkernel32.NewProc("GetTickCount64")
)

type filetime struct {
	LowDateTime  uint32
	HighDateTime uint32
}
type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

func filetimeToUint64(ft filetime) uint64 {
	return (uint64(ft.HighDateTime) << 32) | uint64(ft.LowDateTime)
}
func getCPUTimes() (idle, total uint64, err error) {
	var idleTime, kernelTime, userTime filetime
	r1, _, err := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)
	if r1 == 0 {
		return 0, 0, err
	}
	idle = filetimeToUint64(idleTime)
	kernel := filetimeToUint64(kernelTime)
	user := filetimeToUint64(userTime)
	total = kernel + user
	return idle, total, nil
}
func getMemoryStats() (used, total, percent float64, err error) {
	var memInfo memoryStatusEx
	memInfo.Length = uint32(unsafe.Sizeof(memInfo))
	r1, _, err := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if r1 == 0 {
		return 0, 0, 0, err
	}
	totalGB := float64(memInfo.TotalPhys) / (1024 * 1024 * 1024)
	freeGB := float64(memInfo.AvailPhys) / (1024 * 1024 * 1024)
	usedGB := totalGB - freeGB
	percent = float64(memInfo.MemoryLoad)
	return usedGB, totalGB, percent, nil
}
func getDiskStats() (used, total, percent float64, err error) {
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	err = windows.GetDiskFreeSpaceEx(windows.StringToUTF16Ptr("C:\\"), &freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		return 0, 0, 0, err
	}
	usedBytes := totalNumberOfBytes - totalNumberOfFreeBytes
	totalGB := float64(totalNumberOfBytes) / (1024 * 1024 * 1024)
	usedGB := float64(usedBytes) / (1024 * 1024 * 1024)
	percent = (usedGB / totalGB) * 100.0
	return usedGB, totalGB, percent, nil
}
func getUptime() string {
	r1, _, _ := procGetTickCount64.Call()
	uptimeSec := float64(r1) / 1000.0
	days := int(uptimeSec) / (24 * 3600)
	hours := (int(uptimeSec) % (24 * 3600)) / 3600
	minutes := (int(uptimeSec) % 3600) / 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
