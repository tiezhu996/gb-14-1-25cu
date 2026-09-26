package dto

// IDRequest 路径 ID 参数。
type IDRequest struct {
	ID string `uri:"id" binding:"required"`
}

// PageRequest 分页查询参数。
type PageRequest struct {
	Page     int64 `form:"page" json:"page"`
	PageSize int64 `form:"page_size" json:"page_size"`
}
