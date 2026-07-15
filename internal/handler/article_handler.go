package handler

import (
	"net/http"
	"strconv"

	"github.com/you/sharing-vision-backend-v2/internal/config"
	"github.com/you/sharing-vision-backend-v2/internal/model"
	"github.com/you/sharing-vision-backend-v2/internal/service"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

func getUserID(c *gin.Context) int {
	if u, ok := c.Get("user"); ok {
		if usr, ok := u.(*model.User); ok {
			return usr.ID
		}
	}
	return 0
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := getUserID(c)
	if userID == 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	post, err := h.articleService.Create(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// audit log
	log := &model.AuditLog{
		ActorUserID:  ptrInt(userID),
		Action:       "create",
		ResourceType: "post",
		ResourceID:   ptrInt(post.ID),
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		NewValues: map[string]any{
			"title":    post.Title,
			"category": post.Category,
			"status":   post.Status,
		},
	}
	_ = h.articleService.CreateAuditLog(log)

	c.JSON(http.StatusCreated, post)
}

func (h *ArticleHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	post, err := h.articleService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *ArticleHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldPost, _ := h.articleService.GetByID(id)
	if oldPost == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	m := map[string]interface{}{}
	if req.Title != nil {
		m["title"] = *req.Title
	}
	if req.Content != nil {
		m["content"] = *req.Content
	}
	if req.Category != nil {
		m["category"] = *req.Category
	}
	if req.Status != nil {
		m["status"] = *req.Status
	}
	if len(m) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := h.articleService.Update(id, m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newPost, _ := h.articleService.GetByID(id)
	oldValues := map[string]any{
		"title": oldPost.Title,
		"category": oldPost.Category,
		"status": oldPost.Status,
	}
	newValues := map[string]any{
		"title": newPost.Title,
		"category": newPost.Category,
		"status": newPost.Status,
	}
	log := &model.AuditLog{
		ActorUserID:  ptrInt(getUserID(c)),
		Action:       "update",
		ResourceType: "post",
		ResourceID:   ptrInt(id),
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		OldValues:    oldValues,
		NewValues:    newValues,
	}
	_ = h.articleService.CreateAuditLog(log)

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	oldPost, _ := h.articleService.GetByID(id)
	_ = h.articleService.Update(id, map[string]interface{}{"status": "thrash"})
	_ = h.articleService.CreateAuditLog(&model.AuditLog{
		ActorUserID:  ptrInt(getUserID(c)),
		Action:       "thrash",
		ResourceType: "post",
		ResourceID:   ptrInt(id),
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		OldValues: map[string]any{
			"title": oldPost.Title,
			"category": oldPost.Category,
			"status": oldPost.Status,
		},
	})
	c.JSON(http.StatusOK, gin.H{"message": "moved to thrash"})
}

func (h *ArticleHandler) List(c *gin.Context) {
	var q model.PostListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, total, err := h.articleService.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items":  items,
		"total":  total,
		"limit":  q.Limit,
		"offset": q.Offset,
	})
}

func (h *ArticleHandler) PublicList(c *gin.Context) {
	var q model.PostListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Support path-based pagination: /articles/:limit/:offset
	if limitStr := c.Param("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			q.Limit = l
		}
	}
	if offsetStr := c.Param("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			q.Offset = o
		}
	}
	q.Status = "publish"
	items, total, err := h.articleService.List(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cats, _ := h.articleService.ListCategories()
	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"total":     total,
		"categories": cats,
	})
}

func (h *ArticleHandler) Dashboard(c *gin.Context) {
	published, _ := h.articleService.CountByStatus("publish")
	draft, _ := h.articleService.CountByStatus("draft")
	thrash, _ := h.articleService.CountByStatus("thrash")

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"published": published,
			"draft":     draft,
			"thrash":    thrash,
		},
		"rate_limit": gin.H{
			"window": "1m",
			"limit":  config.Conf.RateLimit.AdminRPS,
		},
	})
}

func (h *ArticleHandler) AuditLogs(c *gin.Context) {
	limit := 50
	offset := 0
	if v := c.Query("limit"); v != "" {
		limit, _ = strconv.Atoi(v)
	}
	if v := c.Query("offset"); v != "" {
		offset, _ = strconv.Atoi(v)
	}
	resourceType := c.Query("resource_type")
	logs, total, err := h.articleService.ListAuditLogs(limit, offset, resourceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items":  logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func ptrInt(i int) *int {
	return &i
}
