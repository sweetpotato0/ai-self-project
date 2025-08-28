package service

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/netutil"
)

// ToolsService 工具服务接口
type ToolsServiceInterface interface {
	// 网络工具
	ScanPorts(target string, ports []int, timeout time.Duration, concurrent int) []netutil.PortScanResult
	LookupDNS(domain, recordType string) DNSResponse
}

// ToolsService 工具服务实现
type ToolsService struct {
	logger logger.LoggerInterface
}

// NewToolsService 创建工具服务
func NewToolsService(logger logger.LoggerInterface) ToolsServiceInterface {
	return &ToolsService{
		logger: logger,
	}
}

// DNSRecord DNS记录
type DNSRecord struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   int    `json:"ttl,omitempty"`
}

// DNSResponse DNS查询响应
type DNSResponse struct {
	Domain  string      `json:"domain"`
	Type    string      `json:"type"`
	Records []DNSRecord `json:"records"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
}

// ScanPorts 执行端口扫描
func (s *ToolsService) ScanPorts(target string, ports []int, timeout time.Duration, concurrent int) []netutil.PortScanResult {
	results := make([]netutil.PortScanResult, len(ports))
	sem := make(chan struct{}, concurrent) // 并发控制
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 常用端口服务映射
	commonPorts := map[int]struct {
		Service     string
		Description string
	}{
		21:   {"FTP", "File Transfer Protocol"},
		22:   {"SSH", "Secure Shell"},
		23:   {"Telnet", "Telnet Protocol"},
		25:   {"SMTP", "Simple Mail Transfer Protocol"},
		53:   {"DNS", "Domain Name System"},
		80:   {"HTTP", "HyperText Transfer Protocol"},
		110:  {"POP3", "Post Office Protocol v3"},
		143:  {"IMAP", "Internet Message Access Protocol"},
		443:  {"HTTPS", "HTTP over TLS/SSL"},
		993:  {"IMAPS", "IMAP over TLS/SSL"},
		995:  {"POP3S", "POP3 over TLS/SSL"},
		1433: {"MSSQL", "Microsoft SQL Server"},
		3306: {"MySQL", "MySQL Database"},
		3389: {"RDP", "Remote Desktop Protocol"},
		5432: {"PostgreSQL", "PostgreSQL Database"},
		5900: {"VNC", "Virtual Network Computing"},
		8080: {"HTTP-Alt", "HTTP Alternative Port"},
		8443: {"HTTPS-Alt", "HTTPS Alternative Port"},
	}

	for i, port := range ports {
		wg.Add(1)
		go func(index, p int) {
			defer wg.Done()
			sem <- struct{}{}        // 获取信号量
			defer func() { <-sem }() // 释放信号量

			result := s.scanSinglePort(target, p, timeout, commonPorts)

			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, port)
	}

	wg.Wait()
	return results
}

// scanSinglePort 扫描单个端口
func (s *ToolsService) scanSinglePort(target string, port int, timeout time.Duration, commonPorts map[int]struct {
	Service     string
	Description string
}) netutil.PortScanResult {
	start := time.Now()

	address := net.JoinHostPort(target, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, timeout)

	result := netutil.PortScanResult{
		Port: port,
	}

	// 添加服务信息
	if portInfo, exists := commonPorts[port]; exists {
		result.Service = portInfo.Service
		result.Description = portInfo.Description
	}

	if err != nil {
		result.Status = "closed"
		return result
	}

	defer conn.Close()

	responseTime := int(time.Since(start).Milliseconds())
	result.Status = "open"
	result.ResponseTime = &responseTime

	return result
}

// LookupDNS 执行DNS查询
func (s *ToolsService) LookupDNS(domain, recordType string) DNSResponse {
	response := DNSResponse{
		Domain:  domain,
		Type:    recordType,
		Status:  "success",
		Records: []DNSRecord{},
	}

	switch strings.ToUpper(recordType) {
	case "A":
		if ips, err := net.LookupIP(domain); err == nil {
			for _, ip := range ips {
				if ip.To4() != nil { // IPv4地址
					response.Records = append(response.Records, DNSRecord{
						Type:  "A",
						Value: ip.String(),
					})
				}
			}
		} else {
			response.Status = "error"
			response.Message = fmt.Sprintf("A记录查询失败: %v", err)
		}

	case "AAAA":
		if ips, err := net.LookupIP(domain); err == nil {
			for _, ip := range ips {
				if ip.To4() == nil && ip.To16() != nil { // IPv6地址
					response.Records = append(response.Records, DNSRecord{
						Type:  "AAAA",
						Value: ip.String(),
					})
				}
			}
		} else {
			response.Status = "error"
			response.Message = fmt.Sprintf("AAAA记录查询失败: %v", err)
		}

	case "CNAME":
		if cname, err := net.LookupCNAME(domain); err == nil {
			// 移除末尾的点
			cname = strings.TrimSuffix(cname, ".")
			response.Records = append(response.Records, DNSRecord{
				Type:  "CNAME",
				Value: cname,
			})
		} else {
			response.Status = "error"
			response.Message = fmt.Sprintf("CNAME记录查询失败: %v", err)
		}

	case "MX":
		if mxRecords, err := net.LookupMX(domain); err == nil {
			for _, mx := range mxRecords {
				response.Records = append(response.Records, DNSRecord{
					Type:  "MX",
					Value: fmt.Sprintf("%d %s", mx.Pref, strings.TrimSuffix(mx.Host, ".")),
				})
			}
		} else {
			response.Status = "error"
			response.Message = fmt.Sprintf("MX记录查询失败: %v", err)
		}

	case "NS":
		if nsRecords, err := net.LookupNS(domain); err == nil {
			for _, ns := range nsRecords {
				response.Records = append(response.Records, DNSRecord{
					Type:  "NS",
					Value: strings.TrimSuffix(ns.Host, "."),
				})
			}
		} else {
			response.Status = "error"
			response.Message = fmt.Sprintf("NS记录查询失败: %v", err)
		}

	case "TXT":
		if txtRecords, err := net.LookupTXT(domain); err == nil {
			for _, txt := range txtRecords {
				response.Records = append(response.Records, DNSRecord{
					Type:  "TXT",
					Value: txt,
				})
			}
		} else {
			response.Status = "error"
			response.Message = fmt.Sprintf("TXT记录查询失败: %v", err)
		}

	case "PTR":
		// 反向DNS查询
		if names, err := net.LookupAddr(domain); err == nil {
			for _, name := range names {
				response.Records = append(response.Records, DNSRecord{
					Type:  "PTR",
					Value: strings.TrimSuffix(name, "."),
				})
			}
		} else {
			response.Status = "error"
			response.Message = fmt.Sprintf("PTR记录查询失败: %v", err)
		}

	default:
		response.Status = "error"
		response.Message = fmt.Sprintf("不支持的DNS记录类型: %s", recordType)
	}

	return response
}