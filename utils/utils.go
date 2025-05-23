package utils

type Device struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	MAC      string `json:"mac"`
}

type ServiceInfo struct {
	Service string
	Version string
}

var CommonPorts = map[int]string{
	21:    "FTP",
	22:    "SSH",
	23:    "Telnet",
	25:    "SMTP",
	53:    "DNS",
	80:    "HTTP",
	110:   "POP3",
	119:   "NNTP",
	123:   "NTP",
	135:   "Microsoft RPC",
	139:   "NetBIOS Session Service",
	143:   "IMAP",
	161:   "SNMP",
	194:   "IRC",
	443:   "HTTPS",
	445:   "Microsoft SMB",
	465:   "SMTPS",
	514:   "Syslog",
	587:   "SMTP (Mail Submission)",
	993:   "IMAPS",
	995:   "POP3S",
	1433:  "Microsoft SQL Server",
	1521:  "Oracle Database",
	2049:  "NFS",
	2181:  "Zookeeper",
	3306:  "MySQL",
	3389:  "RDP (Remote Desktop Protocol)",
	5432:  "PostgreSQL",
	5672:  "RabbitMQ",
	5900:  "VNC",
	6379:  "Redis",
	8080:  "HTTP-Proxy",
	8443:  "HTTPS-Alt",
	9200:  "Elasticsearch",
	11211: "Memcached",
	27017: "MongoDB",
}

// Map of common UDP ports and their services
var UdpServices = map[int]string{
	53:    "DNS (Domain Name System)",
	67:    "DHCP (Dynamic Host Configuration Protocol) - Server",
	68:    "DHCP (Dynamic Host Configuration Protocol) - Client",
	69:    "TFTP (Trivial File Transfer Protocol)",
	123:   "NTP (Network Time Protocol)",
	161:   "SNMP (Simple Network Management Protocol)",
	162:   "SNMP Trap",
	500:   "ISAKMP (Internet Security Association and Key Management Protocol)",
	514:   "Syslog",
	520:   "RIP (Routing Information Protocol)",
	1812:  "RADIUS (Authentication)",
	1813:  "RADIUS (Accounting)",
	2049:  "NFS (Network File System)",
	5353:  "mDNS (Multicast DNS)",
	33434: "Traceroute",
	137:   "NetBIOS Name Service",
	138:   "NetBIOS Datagram Service",
	1434:  "Microsoft SQL Monitor",
	1900:  "SSDP (Simple Service Discovery Protocol)",
	4500:  "IPSec NAT Traversal",
	5355:  "LLMNR (Link-Local Multicast Name Resolution)",
	5683:  "CoAP (Constrained Application Protocol)",
	64738: "Mumble (Voice Chat)",
}
