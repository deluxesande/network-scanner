package subnet

import "testing"

// Guards the OS-specific ping output parsing (lowercase ttl= on linux/darwin,
// TTL= on windows). Loopback always answers.
func TestPingWithTTLParsesLocalhost(t *testing.T) {
	alive, os := PingWithTTL("127.0.0.1")
	if !alive || os == "Unknown" {
		t.Fatalf("loopback ping not parsed: alive=%v os=%q", alive, os)
	}
}
