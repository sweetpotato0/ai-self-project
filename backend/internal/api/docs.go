package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// APIDoc API文档结构
type APIDoc struct {
	OpenAPI    string              `json:"openapi"`
	Info       *APIInfo            `json:"info"`
	Servers    []*APIServer        `json:"servers"`
	Paths      map[string]*APIPath `json:"paths"`
	Components *APIComponents      `json:"components"`
	Tags       []*APITag           `json:"tags"`
}

// APIInfo API信息
type APIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// APIServer 服务器信息
type APIServer struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// APIPath API路径
type APIPath struct {
	Get    *APIOperation `json:"get,omitempty"`
	Post   *APIOperation `json:"post,omitempty"`
	Put    *APIOperation `json:"put,omitempty"`
	Delete *APIOperation `json:"delete,omitempty"`
}

// APIOperation API操作
type APIOperation struct {
	Tags        []string                `json:"tags"`
	Summary     string                  `json:"summary"`
	Description string                  `json:"description"`
	OperationID string                  `json:"operationId"`
	Responses   map[string]*APIResponse `json:"responses"`
	Security    []*APISecurity          `json:"security,omitempty"`
}

// APIDocResponse API文档响应
type APIDocResponse struct {
	Description string                 `json:"description"`
	Content     map[string]*APIContent `json:"content,omitempty"`
}

// APIContent 内容类型
type APIContent struct {
	Schema *APISchema `json:"schema"`
}

// APISchema API模式
type APISchema struct {
	Type       string                `json:"type,omitempty"`
	Properties map[string]*APISchema `json:"properties,omitempty"`
	Ref        string                `json:"$ref,omitempty"`
}

// APISecurity 安全定义
type APISecurity struct {
	BearerAuth []string `json:"bearerAuth,omitempty"`
}

// APITag API标签
type APITag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// APIComponents API组件
type APIComponents struct {
	Schemas         map[string]*APISchema         `json:"schemas,omitempty"`
	SecuritySchemes map[string]*APISecurityScheme `json:"securitySchemes,omitempty"`
}

// APISecurityScheme 安全方案
type APISecurityScheme struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
}

// GenerateDocsHandler 生成文档处理器
func GenerateDocsHandler() gin.HandlerFunc {
	doc := &APIDoc{
		OpenAPI: "3.0.3",
		Info: &APIInfo{
			Title:       "TODO API",
			Description: "A comprehensive TODO application API",
			Version:     "1.0.0",
		},
		Servers: []*APIServer{
			{
				URL:         "http://localhost:8080",
				Description: "Development server",
			},
		},
		Paths: make(map[string]*APIPath),
		Components: &APIComponents{
			Schemas:         make(map[string]*APISchema),
			SecuritySchemes: make(map[string]*APISecurityScheme),
		},
		Tags: []*APITag{
			{Name: "Authentication", Description: "User authentication"},
			{Name: "Users", Description: "User management"},
			{Name: "Todos", Description: "TODO management"},
			{Name: "Health", Description: "Health check"},
		},
	}

	// 添加安全方案
	doc.Components.SecuritySchemes["bearerAuth"] = &APISecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}

	return func(c *gin.Context) {
		jsonData, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate documentation"})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Data(200, "application/json", jsonData)
	}
}

// GenerateSwaggerUIHandler 生成Swagger UI处理器
func GenerateSwaggerUIHandler() gin.HandlerFunc {
	swaggerHTML := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/api/v1/docs/json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.Data(200, "text/html", []byte(swaggerHTML))
	}
}
