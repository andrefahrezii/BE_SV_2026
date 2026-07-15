package repository

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/you/sharing-vision-backend-v2/internal/model"
)

type ArticleRepo struct {
	db *sql.DB
}

func NewArticleRepo(db *sql.DB) *ArticleRepo {
	return &ArticleRepo{db: db}
}

func (r *ArticleRepo) Create(post *model.Post) (int, error) {
	query := `
		INSERT INTO sv_portal.posts (title, content, category, status, author_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_date, updated_date`
	err := r.db.QueryRow(query, post.Title, post.Content, post.Category, post.Status, post.AuthorID).
		Scan(&post.ID, &post.CreatedDate, &post.UpdatedDate)
	return post.ID, err
}

func (r *ArticleRepo) FindByID(id int) (*model.Post, error) {
	post := &model.Post{}
	query := `SELECT id, title, content, category, created_date, updated_date, status, author_id
		FROM sv_portal.posts WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.Category,
		&post.CreatedDate, &post.UpdatedDate, &post.Status, &post.AuthorID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return post, err
}

func (r *ArticleRepo) Update(id int, m map[string]interface{}) error {
	parts := []string{}
	args := []interface{}{}
	idx := 1
	if v, ok := m["title"]; ok {
		parts = append(parts, "title=$"+itoa(idx))
		args = append(args, v)
		idx++
	}
	if v, ok := m["content"]; ok {
		parts = append(parts, "content=$"+itoa(idx))
		args = append(args, v)
		idx++
	}
	if v, ok := m["category"]; ok {
		parts = append(parts, "category=$"+itoa(idx))
		args = append(args, v)
		idx++
	}
	if v, ok := m["status"]; ok {
		parts = append(parts, "status=$"+itoa(idx))
		args = append(args, v)
		idx++
	}
	parts = append(parts, "updated_date=NOW()")
	args = append(args, id)
	query := "UPDATE sv_portal.posts SET " + strings.Join(parts, ", ") + " WHERE id=$" + itoa(idx)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ArticleRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM sv_portal.posts WHERE id=$1", id)
	return err
}

func (r *ArticleRepo) List(q model.PostListQuery) ([]model.PostListItem, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	idx := 1

	if q.Q != "" {
		// Full-text search instead of LIKE
		where = append(where, "to_tsvector('indonesian', coalesce(title,'') || ' ' || coalesce(category,'')) @@ plainto_tsquery('indonesian', $"+itoa(idx)+")")
		args = append(args, q.Q)
		idx++
	}
	if q.Category != "" {
		where = append(where, "category=$"+itoa(idx))
		args = append(args, q.Category)
		idx++
	}
	if q.Status != "" {
		where = append(where, "status=$"+itoa(idx))
		args = append(args, q.Status)
		idx++
	}

	base := "SELECT id, title, category, status, created_date FROM sv_portal.posts WHERE " + strings.Join(where, " AND ")
	countQ := "SELECT COUNT(*) FROM sv_portal.posts WHERE " + strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	sortBy := "created_date"
	order := "DESC"
	if q.SortBy != "" {
		sortBy = q.SortBy
	}
	if q.Order != "" {
		order = strings.ToUpper(q.Order)
	}

	limit := q.Limit
	offset := q.Offset
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	args = append(args, limit, offset)
	query := base + " ORDER BY " + sortBy + " " + order + " LIMIT $" + itoa(idx) + " OFFSET $" + itoa(idx+1)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := []model.PostListItem{}
	for rows.Next() {
		var it model.PostListItem
		if err := rows.Scan(&it.ID, &it.Title, &it.Category, &it.Status, &it.CreatedDate); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, rows.Err()
}

func (r *ArticleRepo) CountByStatus(status string) (int, error) {
	var c int
	err := r.db.QueryRow("SELECT COUNT(*) FROM sv_portal.posts WHERE status=$1", status).Scan(&c)
	return c, err
}

func (r *ArticleRepo) ListCategories() ([]string, error) {
	// Union with master categories table so admin-created categories also appear
	query := `SELECT DISTINCT category FROM sv_portal.posts WHERE category IS NOT NULL
		UNION
		SELECT name FROM sv_portal.categories
		ORDER BY category`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, rows.Err()
}

func (r *ArticleRepo) ListAuditLogs(limit, offset int, resourceType string) ([]model.AuditLog, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	idx := 1
	if resourceType != "" {
		where = append(where, "resource_type=$"+itoa(idx))
		args = append(args, resourceType)
		idx++
	}
	base := "SELECT id, actor_user_id, action, resource_type, resource_id, ip_address, user_agent, old_values, new_values, created_at FROM sv_portal.audit_logs WHERE " + strings.Join(where, " AND ")
	countQ := "SELECT COUNT(*) FROM sv_portal.audit_logs WHERE " + strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	args = append(args, limit, offset)
	query := base + " ORDER BY created_at DESC LIMIT $" + itoa(idx) + " OFFSET $" + itoa(idx+1)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	logs := []model.AuditLog{}
	for rows.Next() {
		var log model.AuditLog
		var oldJSON, newJSON []byte
		if err := rows.Scan(&log.ID, &log.ActorUserID, &log.Action, &log.ResourceType, &log.ResourceID, &log.IPAddress, &log.UserAgent, &oldJSON, &newJSON, &log.CreatedAt); err != nil {
			return nil, 0, err
		}
		if oldJSON != nil {
			log.OldValues = map[string]any{}
			json.Unmarshal(oldJSON, &log.OldValues)
		}
		if newJSON != nil {
			log.NewValues = map[string]any{}
			json.Unmarshal(newJSON, &log.NewValues)
		}
		logs = append(logs, log)
	}
	return logs, total, rows.Err()
}

func (r *ArticleRepo) CreateAuditLog(log *model.AuditLog) error {
	query := `INSERT INTO sv_portal.audit_logs (actor_user_id, action, resource_type, resource_id, ip_address, user_agent, old_values, new_values)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	var oldJSON, newJSON []byte
	if log.OldValues != nil {
		oldJSON, _ = json.Marshal(log.OldValues)
	}
	if log.NewValues != nil {
		newJSON, _ = json.Marshal(log.NewValues)
	}
	return r.db.QueryRow(query, log.ActorUserID, log.Action, log.ResourceType, log.ResourceID, log.IPAddress, log.UserAgent, oldJSON, newJSON).Scan(&log.ID)
}

func itoa(i int) string {
	return strconv.Itoa(i)
}
