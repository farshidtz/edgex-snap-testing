package utils

import (
	"fmt"
	"net"
	"testing"
	"time"
)

const dialTimeout = 2 * time.Second

// WaitServiceOnline waits for a service to come online by dialing its port(s)
// up to a maximum number
func WaitServiceOnline(t *testing.T, timeout time.Duration, ports ...string) {
	start := time.Now()

	portsMap := make(map[string]bool)
	for _, port := range ports {
		portsMap[port] = false
	}
	closePorts := func() (pts []string) {
		for port, open := range portsMap {
			if !open {
				pts = append(pts, port)
			}
		}
		return pts
	}

	for time.Since(start) < timeout {

		t.Logf("Waiting for service ports %s. Elapsed %v/%v", closePorts(), time.Since(start), timeout)

		for _, port := range closePorts() {

			conn, err := net.DialTimeout("tcp", ":"+port, dialTimeout)
			if conn != nil {
				t.Logf("Service port %s is open.", port)
				portsMap[port] = true
				continue
			}

			if err != nil {
				t.Fatalf("Timed out waiting for service ports: %v. Error:\n%v", closePorts(), err)
			} else if conn == nil {
				t.Fatalf("Timed out waiting for service ports: %v", closePorts())
			}
		}
	}
}

// RequirePortOpen checks if the local port(s) accepts connections
func RequirePortOpen(t *testing.T, host string, ports ...string) {

	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", ":"+port, dialTimeout)
		if err != nil {
			conn.Close()
			t.Errorf("Port %s is not open: %s", port, err)
		}

		if conn == nil {
			t.Errorf("Port %s is not open", port)
		}

		if conn != nil {
			t.Logf("Port %v is open.", port)
			conn.Close()
		}
	}
	if t.Failed() {
		t.FailNow()
	}
}

// checkListenAllInterfaces checks if the port(s) listen on all interfaces
func RequireListenAllInterfaces(t *testing.T, mustListen bool, ports ...string) {
	for _, port := range ports {
		stdout, _ := Exec(t, "sudo lsof -nPi :"+port+" | { grep \\* || true; }")
		isListening := (stdout != "")

		if mustListen && !isListening {
			t.Errorf("Port %v not listening to all interfaces", port)
		} else if !mustListen && isListening {
			t.Errorf("Port %v is listening to all interfaces", port)
		}
	}
	if t.Failed() {
		t.FailNow()
	}
}

// RequireListenLoopback checks if the port(s) listen on the loopback interface
func RequireListenLoopback(t *testing.T, ports ...string) {
	for _, port := range ports {
		stdout, _ := Exec(t, "sudo lsof -nPi :"+port+" | { grep 127.0.0.1  || true; }")
		isListening := stdout != ""

		if !isListening {
			t.Errorf("Port %v not listening on loopback interface", port)
		}
	}
	if t.Failed() {
		t.FailNow()
	}
}

// RequirePortAvailable checks if a port is available (not open) locally
func RequirePortAvailable(t *testing.T, port string) {
	stdout, _ := Exec(t, fmt.Sprintf("sudo lsof -nPi :%s || true", port))
	if stdout != "" {
		t.Fatalf("Port %s is not available", port)
	}
	t.Logf("Port %s is available.", port)
}
