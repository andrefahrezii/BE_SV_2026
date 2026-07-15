package model

import "time"

type Post struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	Category    string    `json:"category" db:"category"`
	CreatedDate time.Time `json:"created_date" db:"created_date"`
	UpdatedDate time.Time `json:"updated_date" db:"updated_date"`
	Status      string    `json:"status" db:"status"`
	AuthorID    int       `json:"author_id" db:"author_id"`
}

type CreatePostRequest struct {
	Title    string `json:"title" binding:"required,min=20"`
	Content  string `json:"content" binding:"required,min=200"`
	Category string `json:"category" binding:"required,min=3"`
	Status   string `json:"status" binding:"required,oneof=publish draft thrash"`
}

type UpdatePostRequest struct {
	Title    *string `json:"title" binding:"omitempty,min=20"`
	Content  *string `json:"content" binding:"omitempty,min=200"`
	Category *string `json:"category" binding:"omitempty,min=3"`
	Status   *string `json:"status" binding:"omitempty,oneof=publish draft thrash"`
}

type PostListQuery struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Q        string `form:"q"`
	Category string `form:"category"`
	Status   string `form:"status"`
	SortBy   string `form:"sort_by"`
	Order    string `form:"order"`
}

type PostListItem struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`
	CreatedDate time.Time `json:"created_date"`
}

type AuditLog struct {
	ID          int             `json:"id" db:"id"`
	ActorUserID *int            `json:"actor_user_id" db:"actor_user_id"`
	Action      string          `json:"action" db:"action"`
	ResourceType string         `json:"resource_type" db:"resource_type"`
	ResourceID  *int            `json:"resource_id" db:"resource_id"`
	IPAddress   string          `json:"ip_address" db:"ip_address"`
	UserAgent   string          `json:"user_agent" db:"user_agent"`
	OldValues   map[string]any  `json:"old_values" db:"old_values"`
	NewValues   map[string]any  `json:"new_values" db:"new_values"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}
