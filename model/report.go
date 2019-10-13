package model

// 举报系统
type Report struct {
	ID int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`
	UpdatedAt int64 `json:"updated_at" meddler:"updated_at"`

	// 举报帖子/用户ID
	EntityId int64 `json:"entity_id" meddler:"entity_id"`
	EntityTy int64 `json:"entity_ty" meddler:"entity_ty"`

	// 大概内容
	Content  string `json:"content" meddler:"content"`

	// 举报次数
	Counter  int64 `json:"counter" meddler:"counter"`

	// 状态：待处理/已处理
	Status   int64 `json:"status" meddler:"status"`

	// 举报人
	UserId   int64 `json:"user_id" meddler:"user_id"`

	// 举报类型
	ReportTy int64 `json:"report_ty" meddler:"report_ty"`

	// 附言
	Other    string `json:"other" meddler:"other"`

	// 图片 [...]
	Images   string `json:"images" meddler:"images"`
}

