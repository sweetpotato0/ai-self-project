#!/bin/bash

echo "登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "更新文章内容..."
curl -X PUT http://localhost:8080/api/v1/articles/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "content": "# 大语言模型技术深度解析\n\n## 引言\n\n大语言模型是AI领域的重要突破，从GPT到Claude展现了技术的演进。\n\n## 核心技术\n\n- Transformer架构\n- 自注意力机制\n- 预训练与微调\n\n## 应用场景\n\n- 自然语言处理\n- 代码生成\n- 创意内容创作\n\n## 未来趋势\n\n- 多模态融合\n- 推理能力增强\n- 个性化定制",
    "summary": "深入探讨大语言模型的技术原理、发展历程和未来趋势",
    "cover_image": "https://images.unsplash.com/photo-1677442136019-21780ecad995?w=800&h=400&fit=crop"
  }'

echo "更新完成！"
