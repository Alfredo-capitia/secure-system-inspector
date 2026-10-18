package main

import "testing"

func TestParseTCPLineIPv4(t *testing.T) {
	line := "0: 0100007F:0277 00000000:0000 0A 00000000:00000000"

	socket, ok := parseTCPLine(line, false)

	if !ok {
		t.Fatal("expected a listening socket to be parsed")
	}

	if socket.Address != "127.0.0.1" {
		t.Errorf("expected address 127.0.0.1, got %s", socket.Address)
	}

	if socket.Port != 631 {
		t.Errorf("expected port 631, got %d", socket.Port)
	}
}

func TestParseTCPLineIgnoresEstablished(t *testing.T) {
	line := "0: 0100007F:0277 00000000:0000 01 00000000:00000000"

	_, ok := parseTCPLine(line, false)

	if ok {
		t.Fatal("expected an established connection to be ignored")
	}
}

func TestParseTCPLineRejectsMalformedInput(t *testing.T) {
	line := "this is not a valid TCP socket entry"

	_, ok := parseTCPLine(line, false)

	if ok {
		t.Fatal("expected malformed input to be rejected")
	}
}

func TestParseTCPLineIPv6(t *testing.T) {
	line := "0: 00000000000000000000000001000000:0277 00000000000000000000000000000000:0000 0A 00000000:00000000"

	socket, ok := parseTCPLine(line, true)

	if !ok {
		t.Fatal("expected an IPv6 listening socket to be parsed")
	}

	if socket.Address != "::1" {
		t.Errorf("expected address ::1, got %s", socket.Address)
	}

	if socket.Port != 631 {
		t.Errorf("expected port 631, got %d", socket.Port)
	}
}

func TestParseTCPLineMaximumPort(t *testing.T) {
	line := "0: 0100007F:FFFF 00000000:0000 0A 00000000:00000000"

	socket, ok := parseTCPLine(line, false)

	if !ok {
		t.Fatal("expected port 65535 to be accepted")
	}

	if socket.Port != 65535 {
		t.Errorf("expected port 65535, got %d", socket.Port)
	}
}

func TestParseTCPLineInvalidPort(t *testing.T) {
	line := "0: 0100007F:GGGG 00000000:0000 0A 00000000:00000000"

	_, ok := parseTCPLine(line, false)

	if ok {
		t.Fatal("expected invalid hexadecimal port to be rejected")
	}
}
