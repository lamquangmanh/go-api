package utils

import (
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// HostInfo contains information about the current host system.
// This struct is used to gather and store key system metrics for cluster node registration.
type HostInfo struct {
	Hostname string   // The hostname of the machine
	IPs      []net.IP // List of all IPv4 addresses on active interfaces
	CPU      int      // Number of CPU cores
	MemoryMB int      // Total system memory in megabytes
}

// GetTotalMemoryMB reads the system's total memory from /proc/meminfo and returns it in megabytes.
// This function is Linux-specific and reads the MemTotal field from the meminfo file.
// Sample input:  /proc/meminfo contains "MemTotal: 16384000 kB"
// Sample output: 16000
//
// Returns:
//
//	The total system memory in MB, or 0 if unable to read or parse the file
func GetTotalMemoryMB() int {
	// Read the meminfo file which contains memory statistics
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}

	// Split the file content into lines
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Look for the MemTotal line
		if strings.HasPrefix(line, "MemTotal:") {
			// Parse the line: "MemTotal:  16384000 kB"
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				// Convert from KB to MB
				kb, err := strconv.Atoi(fields[1])
				if err != nil {
					return 0
				}
				return kb / 1024 // Convert kilobytes to megabytes
			}
		}
	}
	return 0
}

// DetectHost gathers information about the current host system.
// It collects hostname, CPU count, total memory, and all IPv4 addresses.
// Sample input:  host with hostname="node-1", 8 CPUs, one IP 192.168.1.10
// Sample output: &HostInfo{Hostname:"node-1", CPU:8, MemoryMB:..., IPs:[192.168.1.10]}
//
// This function is useful for:
//   - Node registration in a cluster
//   - System health monitoring
//   - Resource allocation decisions
//
// Returns:
//   - A pointer to HostInfo containing all gathered system information
//   - An error (currently always nil, but included for future extensibility)
func DetectHost() (*HostInfo, error) {
	// Get the system's hostname
	hostname, _ := os.Hostname()

	// Get the number of logical CPU cores
	cpu := runtime.NumCPU()

	// Get total system memory in MB
	mem := GetTotalMemoryMB()

	// Get all IPv4 addresses from active network interfaces
	ips := GetAllIPs()

	// Return the collected system information
	return &HostInfo{
		Hostname: hostname,
		CPU:      cpu,
		MemoryMB: mem,
		IPs:      ips,
	}, nil
}

// GenerateUUID returns a random RFC-4122 UUID string.
// Sample input:  none
// Sample output: "550e8400-e29b-41d4-a716-446655440000"
func GenerateUUID() string {
	return uuid.New().String()
}
