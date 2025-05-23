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
