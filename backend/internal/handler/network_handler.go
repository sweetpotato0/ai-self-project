package handler

import (
	"fmt"
	"net/http"
	"time"

	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/netutil"
	"gin-web-framework/pkg/response"

	"github.com/gin-gonic/gin"
)

// NetworkHandler 网络工具处理器
type NetworkHandler struct {
	logger       logger.LoggerInterface
	toolsService service.ToolsServiceInterface
}

// NewNetworkHandler 创建网络工具处理器
func NewNetworkHandler(logger logger.LoggerInterface, toolsService service.ToolsServiceInterface) *NetworkHandler {
	return &NetworkHandler{
		logger:       logger,
		toolsService: toolsService,
	}
}

// PortScanRequest 端口扫描请求
type PortScanRequest struct {
	Target     string `json:"target" binding:"required"`
	Ports      []int  `json:"ports" binding:"required"`
	Timeout    int    `json:"timeout"`
	Concurrent int    `json:"concurrent"`
}

// PortScanResponse 端口扫描响应
type PortScanResponse struct {
	Target  string                     `json:"target"`
	Results []netutil.PortScanResult `json:"results"`
}

// DNSRequest DNS查询请求
type DNSRequest struct {
	Domain string `json:"domain" binding:"required"`
	Type   string `json:"type"` // 查询类型：A, AAAA, CNAME, MX, NS, TXT, SOA, PTR
}

// PortScan 端口扫描
func (h *NetworkHandler) PortScan(c *gin.Context) {
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
	results := h.toolsService.ScanPorts(req.Target, req.Ports, time.Duration(req.Timeout)*time.Millisecond, req.Concurrent)

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

// DNSLookup DNS查询
func (h *NetworkHandler) DNSLookup(c *gin.Context) {
	var req DNSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid DNS lookup request", "error", err)
		response.Error(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	// 设置默认查询类型
	if req.Type == "" {
		req.Type = "A"
	}

	// 验证域名格式
	if !netutil.IsValidDomain(req.Domain) {
		h.logger.Error("Invalid domain", "domain", req.Domain)
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("无效的域名格式：%s", req.Domain))
		return
	}

	// 执行DNS查询
	dnsResponse := h.toolsService.LookupDNS(req.Domain, req.Type)

	h.logger.WithFields(map[string]any{
		"domain": req.Domain,
		"type":   req.Type,
		"status": dnsResponse.Status,
	}).Info("DNS lookup completed")

	response.Success(c, dnsResponse)
}