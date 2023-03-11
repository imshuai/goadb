package goadb

import (
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixMilli())
}

// 解析设备列表，返回设备序列号列表
func parseDevicesList(devicesList string) []string {
	devices := []string{}
	lines := strings.Split(devicesList, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "*") || strings.HasPrefix(line, "List of devices attached") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 1 && fields[1] == "device" {
			devices = append(devices, fields[0])
		}
	}
	return devices
}

func refreshDevices() error {
	// 运行 adb 命令来获取连接的设备列表
	cmd := exec.Command("adb", "devices")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	// 解析设备列表，获取设备序列号列表
	devices = parseDevicesList(string(out))
	return nil
}

func RandBetween(start, end int) int {
	return rand.Intn(end-start) + start
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
