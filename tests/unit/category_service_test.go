package service_test

import (
	"testing"

	"github.com/you/sharing-vision-backend-v2/internal/model"
	"github.com/you/sharing-vision-backend-v2/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepo struct {
	mock.Mock
}

func (m *MockCategoryRepo) Create(cat *model.Category) error {
	args := m.Called(cat)
	return args.Error(0)
}

func (m *MockCategoryRepo) FindByID(id int) (*model.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepo) FindByName(name string) (*model.Category, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepo) Update(id int, name string) error {
	args := m.Called(id, name)
	return args.Error(0)
}

func (m *MockCategoryRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCategoryRepo) List() ([]model.Category, error) {
	args := m.Called()
	return args.Get(0).([]model.Category), args.Error(1)
}

func TestCategoryService_Create_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	svc := service.NewCategoryService(mockRepo)

	mockRepo.On("FindByName", "Teknologi").Return(nil, nil).Once()
	mockRepo.On("Create", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		cat := args.Get(0).(*model.Category)
		cat.ID = 1
	}).Once()

	cat, err := svc.Create(model.CreateCategoryRequest{Name: "Teknologi"})

	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, "Teknologi", cat.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Create_Duplicate(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	svc := service.NewCategoryService(mockRepo)

	existing := &model.Category{ID: 1, Name: "Teknologi"}
	mockRepo.On("FindByName", "Teknologi").Return(existing, nil).Once()

	cat, err := svc.Create(model.CreateCategoryRequest{Name: "Teknologi"})

	assert.Error(t, err)
	assert.Nil(t, cat)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	svc := service.NewCategoryService(mockRepo)

	expected := &model.Category{ID: 1, Name: "Teknologi Informasi"}
	mockRepo.On("FindByName", "Teknologi Informasi").Return(nil, nil).Once()
	mockRepo.On("Update", 1, "Teknologi Informasi").Return(nil).Once()
	mockRepo.On("FindByID", 1).Return(expected, nil).Once()

	cat, err := svc.Update(1, model.UpdateCategoryRequest{Name: "Teknologi Informasi"})

	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, "Teknologi Informasi", cat.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update_DuplicateName(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	svc := service.NewCategoryService(mockRepo)

	existing := &model.Category{ID: 2, Name: "Teknologi"}
	mockRepo.On("FindByName", "Teknologi").Return(existing, nil).Once()

	cat, err := svc.Update(1, model.UpdateCategoryRequest{Name: "Teknologi"})

	assert.Error(t, err)
	assert.Nil(t, cat)
	assert.Contains(t, err.Error(), "already taken")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Delete_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	svc := service.NewCategoryService(mockRepo)

	mockRepo.On("FindByID", 1).Return(&model.Category{ID: 1, Name: "Teknologi"}, nil).Once()
	mockRepo.On("Delete", 1).Return(nil).Once()

	err := svc.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	svc := service.NewCategoryService(mockRepo)

	mockRepo.On("FindByID", 999).Return(nil, nil).Once()

	err := svc.Delete(999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_List(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	svc := service.NewCategoryService(mockRepo)

	expected := []model.Category{
		{ID: 1, Name: "Berita"},
		{ID: 2, Name: "Teknologi"},
	}
	mockRepo.On("List").Return(expected, nil).Once()

	cats, err := svc.List()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(cats))
	assert.Equal(t, "Berita", cats[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetByID(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	svc := service.NewCategoryService(mockRepo)

	mockRepo.On("FindByID", 1).Return(&model.Category{ID: 1, Name: "Teknologi"}, nil).Once()

	cat, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, "Teknologi", cat.Name)
	mockRepo.AssertExpectations(t)
}
