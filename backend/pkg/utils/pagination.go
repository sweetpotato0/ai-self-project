package utils

// CalculateOffset 计算分页偏移量
func CalculateOffset(page, limit int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * limit
}

// CalculateTotalPages 计算总页数
func CalculateTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	return int((total + int64(limit) - 1) / int64(limit))
}

// ValidatePaginationParams 验证分页参数
func ValidatePaginationParams(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20 // 默认每页20条
	}
	return page, limit
}

// PaginationInfo 分页信息结构
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	Offset     int   `json:"-"`
}

// NewPaginationInfo 创建分页信息
func NewPaginationInfo(page, limit int, total int64) *PaginationInfo {
	page, limit = ValidatePaginationParams(page, limit)
	offset := CalculateOffset(page, limit)
	totalPages := CalculateTotalPages(total, limit)
	
	return &PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Offset:     offset,
	}
}