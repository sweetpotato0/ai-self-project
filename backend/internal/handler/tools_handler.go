package handler

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/netutil"
	"gin-web-framework/pkg/response"

	"github.com/gin-gonic/gin"
)

// ToolsHandler 工具相关的处理器
type ToolsHandler struct {
	logger logger.LoggerInterface
}

// NewToolsHandler 创建新的工具处理器
func NewToolsHandler(logger logger.LoggerInterface) *ToolsHandler {
	return &ToolsHandler{
		logger: logger,
	}
}

// PortScanRequest 端口扫描请求
type PortScanRequest struct {
	Target     string `json:"target" binding:"required"`
	Ports      []int  `json:"ports" binding:"required"`
	Timeout    int    `json:"timeout"`
	Concurrent int    `json:"concurrent"`
}

type PortScanResult = netutil.PortScanResult

// PortScanResponse 端口扫描响应
type PortScanResponse struct {
	Target  string           `json:"target"`
	Results []PortScanResult `json:"results"`
}

// 常用端口服务映射
var commonPorts = map[int]struct {
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

// PortScan 端口扫描
func (h *ToolsHandler) PortScan(c *gin.Context) {
	var req PortScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid port scan request", "error", err)
		response.Error(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	// 设置默认值
	if req.Timeout == 0 {
		req.Timeout = 3000 // 3秒超时
	}
	if req.Concurrent == 0 {
		req.Concurrent = 20 // 默认并发20个
	}

	// 验证目标主机
	if !netutil.IsValidTarget(req.Target) {
		h.logger.Error("Invalid target", "target", req.Target)
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("无效的目标地址：%s。请输入完整的IP地址（如192.168.1.1）或有效域名", req.Target))
		return
	}

	// 验证端口数量限制
	if len(req.Ports) > 1000 {
		h.logger.Error("Too many ports to scan", "count", len(req.Ports))
		response.Error(c, http.StatusBadRequest, "Too many ports (max 1000)")
		return
	}

	// 执行端口扫描
	results := h.scanPorts(req.Target, req.Ports, time.Duration(req.Timeout)*time.Millisecond, req.Concurrent)

	scanResponse := PortScanResponse{
		Target:  req.Target,
		Results: results,
	}

	h.logger.WithFields(map[string]any{
		"target":        req.Target,
		"ports_scanned": len(req.Ports),
		"open_ports":    netutil.CountOpenPorts(results),
	}).Info("Port scan completed")
	response.Success(c, scanResponse)
}

// scanPorts 执行实际的端口扫描
func (h *ToolsHandler) scanPorts(target string, ports []int, timeout time.Duration, concurrent int) []PortScanResult {
	results := make([]PortScanResult, len(ports))
	sem := make(chan struct{}, concurrent) // 并发控制
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, port := range ports {
		wg.Add(1)
		go func(index, p int) {
			defer wg.Done()
			sem <- struct{}{}        // 获取信号量
			defer func() { <-sem }() // 释放信号量

			result := h.scanSinglePort(target, p, timeout)

			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, port)
	}

	wg.Wait()
	return results
}

// scanSinglePort 扫描单个端口
func (h *ToolsHandler) scanSinglePort(target string, port int, timeout time.Duration) PortScanResult {
	start := time.Now()

	address := net.JoinHostPort(target, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, timeout)

	result := PortScanResult{
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
