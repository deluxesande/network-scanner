package tcp

import "testing"

func TestIdentifyService(t *testing.T) {
	cases := []struct {
		name        string
		port        int
		banner      string
		wantService string
		wantVersion string
	}{
		{"ssh", 22, "SSH-2.0-OpenSSH_8.9p1 Ubuntu-3ubuntu0.13\r\n", "SSH", "OpenSSH_8.9p1"},
		{"http server header", 80, "HTTP/1.1 200 OK\r\nServer: nginx/1.24.0\r\n\r\n", "HTTP", "nginx/1.24.0"},
		{"https is still https", 443, "HTTP/1.1 400 Bad Request\r\nServer: Apache/2.4.29\r\n", "HTTPS", "Apache/2.4.29"},
		{"no banner", 3306, "", "MySQL", "Unknown"},
	}

	for _, tc := range cases {
		service, version := identifyService(tc.port, tc.banner)
		if service != tc.wantService || version != tc.wantVersion {
			t.Errorf("%s: got (%q, %q), want (%q, %q)", tc.name, service, version, tc.wantService, tc.wantVersion)
		}
	}
}
