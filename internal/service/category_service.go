package service

import (
	"errors"

	"github.com/you/sharing-vision-backend-v2/internal/model"
)

type CategoryRepo interface {
	Create(cat *model.Category) error
	FindByID(id int) (*model.Category, error)
	FindByName(name string) (*model.Category, error)
	Update(id int, name string) error
	Delete(id int) error
	List() ([]model.Category, error)
}

type CategoryService struct {
	repo CategoryRepo
}

func NewCategoryService(repo CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(req model.CreateCategoryRequest) (*model.Category, error) {
	existing, _ := s.repo.FindByName(req.Name)
	if existing != nil {
		return nil, errors.New("category already exists")
	}
	cat := &model.Category{Name: req.Name}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) GetByID(id int) (*model.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) Update(id int, req model.UpdateCategoryRequest) (*model.Category, error) {
	existing, _ := s.repo.FindByName(req.Name)
	if existing != nil && existing.ID != id {
		return nil, errors.New("category name already taken")
	}
	if err := s.repo.Update(id, req.Name); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *CategoryService) Delete(id int) error {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if cat == nil {
		return errors.New("category not found")
	}
	return s.repo.Delete(id)
}

func (s *CategoryService) List() ([]model.Category, error) {
	return s.repo.List()
}
