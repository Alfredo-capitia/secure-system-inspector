package main

import (
	"encoding/hex"
	"net"
	"strconv"
	"strings"
)

type ListeningSocket struct {
	Address string
	Port    uint16
}

func parseTCPLine(line string, ipv6 bool) (ListeningSocket, bool) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return ListeningSocket{}, false
	}

	// 0A significa TCP LISTEN.
	if fields[3] != "0A" {
		return ListeningSocket{}, false
	}

	addressParts := strings.Split(fields[1], ":")
	if len(addressParts) != 2 {
		return ListeningSocket{}, false
	}

	ipBytes, err := hex.DecodeString(addressParts[0])
	if err != nil {
		return ListeningSocket{}, false
	}

	expectedLength := 4
	if ipv6 {
		expectedLength = 16
	}

	if len(ipBytes) != expectedLength {
		return ListeningSocket{}, false
	}

	if ipv6 {
		for group := 0; group < 4; group++ {
			start := group * 4
			ipBytes[start], ipBytes[start+3] = ipBytes[start+3], ipBytes[start]
			ipBytes[start+1], ipBytes[start+2] = ipBytes[start+2], ipBytes[start+1]
		}
	} else {
		for left, right := 0, len(ipBytes)-1; left < right; left, right = left+1, right-1 {
			ipBytes[left], ipBytes[right] = ipBytes[right], ipBytes[left]
		}
	}

	port, err := strconv.ParseUint(addressParts[1], 16, 16)
	if err != nil {
		return ListeningSocket{}, false
	}

	return ListeningSocket{
		Address: net.IP(ipBytes).String(),
		Port:    uint16(port),
	}, true
}
