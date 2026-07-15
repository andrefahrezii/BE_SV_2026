package service

import (
	"github.com/you/sharing-vision-backend-v2/internal/model"
)

type ArticleRepo interface {
	Create(post *model.Post) (int, error)
	FindByID(id int) (*model.Post, error)
	Update(id int, fields map[string]interface{}) error
	Delete(id int) error
	List(q model.PostListQuery) ([]model.PostListItem, int, error)
	CountByStatus(status string) (int, error)
	ListCategories() ([]string, error)
	ListAuditLogs(limit, offset int, resourceType string) ([]model.AuditLog, int, error)
	CreateAuditLog(log *model.AuditLog) error
}

type ArticleCache interface {
	GetList(q model.PostListQuery) ([]model.PostListItem, int, bool, error)
	SetList(q model.PostListQuery, items []model.PostListItem, total int)
	InvalidateList() error
}

type ArticleService struct {
	repo  ArticleRepo
	cache ArticleCache
}

func NewArticleService(repo ArticleRepo, cache ArticleCache) *ArticleService {
	return &ArticleService{repo: repo, cache: cache}
}

func (s *ArticleService) Create(req model.CreatePostRequest, authorID int) (*model.Post, error) {
	post := &model.Post{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
		AuthorID: authorID,
	}
	id, err := s.repo.Create(post)
	if err != nil {
		return nil, err
	}
	post.ID = id
	_ = s.cache.InvalidateList()
	return post, nil
}

func (s *ArticleService) GetByID(id int) (*model.Post, error) {
	return s.repo.FindByID(id)
}

func (s *ArticleService) Update(id int, m map[string]interface{}) error {
	err := s.repo.Update(id, m)
	if err != nil {
		return err
	}
	_ = s.cache.InvalidateList()
	return nil
}

func (s *ArticleService) Delete(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	_ = s.cache.InvalidateList()
	return nil
}

func (s *ArticleService) List(q model.PostListQuery) ([]model.PostListItem, int, error) {
	if items, total, hit, err := s.cache.GetList(q); hit && err == nil {
		return items, total, nil
	}
	items, total, err := s.repo.List(q)
	if err != nil {
		return nil, 0, err
	}
	s.cache.SetList(q, items, total)
	return items, total, nil
}

func (s *ArticleService) ListCategories() ([]string, error) {
	return s.repo.ListCategories()
}

func (s *ArticleService) CountByStatus(status string) (int, error) {
	return s.repo.CountByStatus(status)
}

func (s *ArticleService) ListAuditLogs(limit, offset int, resourceType string) ([]model.AuditLog, int, error) {
	return s.repo.ListAuditLogs(limit, offset, resourceType)
}

func (s *ArticleService) CreateAuditLog(log *model.AuditLog) error {
	return s.repo.CreateAuditLog(log)
}
