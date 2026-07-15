package service_test

import (
	"testing"

	"github.com/you/sharing-vision-backend-v2/internal/auth"
	"github.com/you/sharing-vision-backend-v2/internal/config"
	"github.com/you/sharing-vision-backend-v2/internal/model"
	"github.com/you/sharing-vision-backend-v2/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	config.Load()
	m.Run()
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(id int) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) CountAdmins() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := service.NewAuthService(mockRepo)

	mockRepo.On("FindByEmail", "new@example.com").Return(nil, nil).Once()
	mockRepo.On("Create", mock.Anything).Return(nil).Once()

	resp, err := svc.Register(model.RegisterRequest{Email: "new@example.com", Password: "password123", Role: "user"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "new@example.com", resp.User.Email)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := service.NewAuthService(mockRepo)

	existingUser := &model.User{ID: 1, Email: "existing@example.com"}
	mockRepo.On("FindByEmail", "existing@example.com").Return(existingUser, nil).Once()

	resp, err := svc.Register(model.RegisterRequest{Email: "existing@example.com", Password: "password123", Role: "user"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "email already registered")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := service.NewAuthService(mockRepo)

	hash, _ := auth.HashPassword("password123")
	existingUser := &model.User{ID: 1, Email: "user@example.com", PasswordHash: hash, Role: "user"}
	mockRepo.On("FindByEmail", "user@example.com").Return(existingUser, nil).Once()

	resp, err := svc.Login(model.LoginRequest{Email: "user@example.com", Password: "password123"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "user@example.com", resp.User.Email)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := service.NewAuthService(mockRepo)

	mockRepo.On("FindByEmail", "user@example.com").Return(nil, nil).Once()

	resp, err := svc.Login(model.LoginRequest{Email: "user@example.com", Password: "wrongpassword"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_BootstrapAdmin_AlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := service.NewAuthService(mockRepo)

	existingUser := &model.User{ID: 1, Email: "admin@sharingvision.id", Role: "admin"}
	mockRepo.On("FindByEmail", "admin@sharingvision.id").Return(existingUser, nil).Once()

	err := svc.BootstrapAdmin()

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_BootstrapAdmin_CreatesNew(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := service.NewAuthService(mockRepo)

	mockRepo.On("FindByEmail", "admin@sharingvision.id").Return(nil, nil).Once()
	mockRepo.On("Create", mock.Anything).Return(nil).Once()

	err := svc.BootstrapAdmin()

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
