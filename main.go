package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Println("================================")
	fmt.Println("   Secure System Inspector")
	fmt.Println("================================")

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		return
	}

	fmt.Println("Hostname:", hostname)
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("CPU Cores:", runtime.NumCPU())
	fmt.Println("Go Version:", runtime.Version())

	fmt.Println("\nTCP LISTENING PORTS (IPv4)")
	fmt.Println("--------------------------------")

	if err := inspectTCPv4("/proc/net/tcp"); err != nil {
		fmt.Println("Error inspecting TCP sockets:", err)
	}
	fmt.Println("\nTCP LISTENING PORTS (IPv6)")
	fmt.Println("--------------------------------")

	if err := inspectTCPv6("/proc/net/tcp6"); err != nil {
		fmt.Println("Error inspecting IPv6 sockets:", err)
	}
}

func inspectTCPv4(path string) error {
	return inspectTCP(path, false)
}

func inspectTCPv6(path string) error {
	return inspectTCP(path, true)
}

func inspectTCP(path string, ipv6 bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Ignora o cabeçalho.
	if scanner.Scan() {
		// Header skipped.
	}

	found := false

	for scanner.Scan() {
		socket, ok := parseTCPLine(scanner.Text(), ipv6)
		if !ok {
			continue
		}

		if ipv6 {
			fmt.Printf("[%s]:%d  LISTEN\n", socket.Address, socket.Port)
		} else {
			fmt.Printf("%s:%d  LISTEN\n", socket.Address, socket.Port)
		}

		found = true
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !found {
		if ipv6 {
			fmt.Println("No IPv6 TCP listening sockets found.")
		} else {
			fmt.Println("No IPv4 TCP listening sockets found.")
		}
	}

	return nil
}
