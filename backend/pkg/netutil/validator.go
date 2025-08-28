package netutil

import (
	"net"
	"strconv"
	"strings"
)

// IsValidTarget 验证目标主机是否有效（IP地址或域名）
func IsValidTarget(target string) bool {
	// 检查是否为IP地址
	if ip := net.ParseIP(target); ip != nil {
		// 严格验证IPv4地址格式（必须是完整的四段式）
		if ip.To4() != nil {
			parts := strings.Split(target, ".")
			if len(parts) != 4 {
				return false // IPv4必须有4段
			}
			// 验证每段都是0-255的数字
			for _, part := range parts {
				if len(part) == 0 || (len(part) > 1 && part[0] == '0') {
					return false // 不允许前导零
				}
				num, err := strconv.Atoi(part)
				if err != nil || num < 0 || num > 255 {
					return false
				}
			}
		}
		
		// 不允许扫描本地环回地址和广播地址等
		if ip.IsLoopback() && target != "127.0.0.1" && target != "::1" {
			return false
		}
		return true
	}

	// 检查是否为有效域名
	if len(target) == 0 || len(target) > 253 {
		return false
	}

	// 简单的域名格式验证
	parts := strings.Split(target, ".")
	if len(parts) < 2 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return false
		}
		// 检查字符是否合法（字母、数字、连字符）
		for _, r := range part {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
				return false
			}
		}
	}

	return true
}

// IsValidPort 验证端口号是否有效
func IsValidPort(port int) bool {
	return port > 0 && port <= 65535
}

// IsValidPortRange 验证端口范围是否有效
func IsValidPortRange(startPort, endPort int) bool {
	return IsValidPort(startPort) && IsValidPort(endPort) && startPort <= endPort
}

// IsValidDomain 验证域名格式是否有效
func IsValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// 移除末尾的点（FQDN格式）
	domain = strings.TrimSuffix(domain, ".")

	// 域名不能以点开头或结尾
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}

	// 分割域名各部分
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return false
		}
		
		// 每部分不能以连字符开头或结尾
		if strings.HasPrefix(part, "-") || strings.HasSuffix(part, "-") {
			return false
		}

		// 检查字符是否合法（字母、数字、连字符）
		for _, r := range part {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
				return false
			}
		}
	}

	return true
}