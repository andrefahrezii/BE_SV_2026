package service_test

import (
	"testing"

	"github.com/you/sharing-vision-backend-v2/internal/model"
	"github.com/you/sharing-vision-backend-v2/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockArticleRepo struct {
	mock.Mock
}

func (m *MockArticleRepo) Create(post *model.Post) (int, error) {
	args := m.Called(post)
	return args.Int(0), args.Error(1)
}

func (m *MockArticleRepo) FindByID(id int) (*model.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockArticleRepo) Update(id int, fields map[string]interface{}) error {
	args := m.Called(id, fields)
	return args.Error(0)
}

func (m *MockArticleRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockArticleRepo) List(q model.PostListQuery) ([]model.PostListItem, int, error) {
	args := m.Called(q)
	return args.Get(0).([]model.PostListItem), args.Int(1), args.Error(2)
}

func (m *MockArticleRepo) CountByStatus(status string) (int, error) {
	args := m.Called(status)
	return args.Int(0), args.Error(1)
}

func (m *MockArticleRepo) ListCategories() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockArticleRepo) ListAuditLogs(limit, offset int, resourceType string) ([]model.AuditLog, int, error) {
	args := m.Called(limit, offset, resourceType)
	return args.Get(0).([]model.AuditLog), args.Int(1), args.Error(2)
}

func (m *MockArticleRepo) CreateAuditLog(log *model.AuditLog) error {
	args := m.Called(log)
	return args.Error(0)
}

type MockArticleCache struct {
	mock.Mock
}

func (m *MockArticleCache) GetList(q model.PostListQuery) ([]model.PostListItem, int, bool, error) {
	args := m.Called(q)
	return args.Get(0).([]model.PostListItem), args.Int(1), args.Bool(2), args.Error(3)
}

func (m *MockArticleCache) SetList(q model.PostListQuery, items []model.PostListItem, total int) {
	m.Called(q, items, total)
}

func (m *MockArticleCache) InvalidateList() error {
	args := m.Called()
	return args.Error(0)
}

func newTestArticleService(repo *MockArticleRepo) (*MockArticleCache, *service.ArticleService) {
	cache := new(MockArticleCache)
	cache.On("InvalidateList").Return(nil).Maybe()
	return cache, service.NewArticleService(repo, cache)
}

func TestArticleService_Create(t *testing.T) {
	mockRepo := new(MockArticleRepo)
	mockCache, svc := newTestArticleService(mockRepo)

	mockRepo.On("Create", mock.Anything).Return(1, nil).Run(func(args mock.Arguments) {
		post := args.Get(0).(*model.Post)
		post.ID = 1
	}).Once()

	post, err := svc.Create(model.CreatePostRequest{Title: "Test Title", Content: "This is a test content that is long enough.", Category: "Tech", Status: "publish"}, 1)

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, 1, post.ID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestArticleService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockArticleRepo)
	_, svc := newTestArticleService(mockRepo)

	mockRepo.On("FindByID", 999).Return(nil, nil).Once()

	post, err := svc.GetByID(999)

	assert.NoError(t, err)
	assert.Nil(t, post)
	mockRepo.AssertExpectations(t)
}

func TestArticleService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockArticleRepo)
	_, svc := newTestArticleService(mockRepo)

	expected := &model.Post{ID: 1, Title: "Found"}
	mockRepo.On("FindByID", 1).Return(expected, nil).Once()

	post, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, post)
	mockRepo.AssertExpectations(t)
}

func TestArticleService_Update(t *testing.T) {
	mockRepo := new(MockArticleRepo)
	mockCache, svc := newTestArticleService(mockRepo)

	mockRepo.On("Update", 1, mock.Anything).Return(nil).Once()

	err := svc.Update(1, map[string]interface{}{"title": "New Title"})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestArticleService_Delete(t *testing.T) {
	mockRepo := new(MockArticleRepo)
	mockCache, svc := newTestArticleService(mockRepo)

	mockRepo.On("Delete", 1).Return(nil).Once()

	err := svc.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestArticleService_CountByStatus(t *testing.T) {
	mockRepo := new(MockArticleRepo)
	_, svc := newTestArticleService(mockRepo)

	mockRepo.On("CountByStatus", "publish").Return(5, nil).Once()

	count, err := svc.CountByStatus("publish")

	assert.NoError(t, err)
	assert.Equal(t, 5, count)
	mockRepo.AssertExpectations(t)
}

func TestArticleService_ListCategories(t *testing.T) {
	mockRepo := new(MockArticleRepo)
	_, svc := newTestArticleService(mockRepo)

	expected := []string{"Tech", "News"}
	mockRepo.On("ListCategories").Return(expected, nil).Once()

	cats, err := svc.ListCategories()

	assert.NoError(t, err)
	assert.Equal(t, expected, cats)
	mockRepo.AssertExpectations(t)
}
