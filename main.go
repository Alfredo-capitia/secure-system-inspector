package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
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
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Ignora a primeira linha, que contém os cabeçalhos.
	if scanner.Scan() {
		// Header skipped.
	}

	found := false

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}

		// A quarta coluna contém o estado TCP.
		if fields[3] != "0A" {
			continue
		}

		addressParts := strings.Split(fields[1], ":")
		if len(addressParts) != 2 {
			continue
		}

		ipBytes, err := hex.DecodeString(addressParts[0])
		if err != nil || len(ipBytes) != 4 {
			continue
		}

		for left, right := 0, len(ipBytes)-1; left < right; left, right = left+1, right-1 {
			ipBytes[left], ipBytes[right] = ipBytes[right], ipBytes[left]
		}

		portNumber, err := strconv.ParseUint(addressParts[1], 16, 16)
		if err != nil {
			continue
		}

		ip := net.IP(ipBytes).String()
		fmt.Printf("%s:%d  LISTEN\n", ip, portNumber)
		found = true
	}

	if err := scanner.Err(); err != nil { // O endereço IPv4 no /proc/net/tcp
		// aparece com os bytes em ordem inversa.
		return err
	}

	if !found {
		fmt.Println("No IPv4 TCP listening sockets found.")
	}

	return nil
}

func inspectTCPv6(path string) error {
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
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}

		// 0A = TCP LISTEN.
		if fields[3] != "0A" {
			continue
		}

		addressParts := strings.Split(fields[1], ":")
		if len(addressParts) != 2 {
			continue
		}

		ipBytes, err := hex.DecodeString(addressParts[0])
		if err != nil || len(ipBytes) != 16 {
			continue
		}

		// O formato IPv6 de /proc/net/tcp6
		// apresenta cada grupo de 32 bits em ordem
		// de bytes inversa em arquiteturas little-endian.
		for group := 0; group < 4; group++ {
			start := group * 4
			ipBytes[start], ipBytes[start+3] = ipBytes[start+3], ipBytes[start]
			ipBytes[start+1], ipBytes[start+2] = ipBytes[start+2], ipBytes[start+1]
		}

		portNumber, err := strconv.ParseUint(addressParts[1], 16, 16)
		if err != nil {
			continue
		}

		ip := net.IP(ipBytes).String()
		fmt.Printf("[%s]:%d  LISTEN\n", ip, portNumber)
		found = true
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !found {
		fmt.Println("No IPv6 TCP listening sockets found.")
	}

	return nil
}
