package netutil

// PortScanResult 端口扫描结果
type PortScanResult struct {
	Port         int    `json:"port"`
	Status       string `json:"status"`
	Service      string `json:"service,omitempty"`
	Description  string `json:"description,omitempty"`
	ResponseTime *int   `json:"responseTime,omitempty"`
}

// CountOpenPorts 计算开放端口数量
func CountOpenPorts(results []PortScanResult) int {
	count := 0
	for _, result := range results {
		if result.Status == "open" {
			count++
		}
	}
	return count
}

// FilterOpenPorts 过滤出开放的端口
func FilterOpenPorts(results []PortScanResult) []PortScanResult {
	var openPorts []PortScanResult
	for _, result := range results {
		if result.Status == "open" {
			openPorts = append(openPorts, result)
		}
	}
	return openPorts
}

// GetPortService 根据端口号获取常见服务名称
func GetPortService(port int) string {
	commonPorts := map[int]string{
		20:   "FTP Data",
		21:   "FTP Control",
		22:   "SSH",
		23:   "Telnet",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		110:  "POP3",
		143:  "IMAP",
		443:  "HTTPS",
		993:  "IMAPS",
		995:  "POP3S",
		3306: "MySQL",
		5432: "PostgreSQL",
		6379: "Redis",
		8080: "HTTP Alt",
		8443: "HTTPS Alt",
		9200: "Elasticsearch",
	}
	
	if service, exists := commonPorts[port]; exists {
		return service
	}
	return ""
}