package service

import (
	"dependency_injection/models"
	"dependency_injection/repository"
)

type DefaultUserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &DefaultUserService{repo: repo}
}

func (s *DefaultUserService) CreateUser(name, email string) error {
	user := &models.User{
		Name:  name,
		Email: email,
	}
	return s.repo.Create(user)
}

func (s *DefaultUserService) GetUser(id uint) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *DefaultUserService) UpdateUser(id uint, name, email string) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	user.Name = name
	user.Email = email
	return s.repo.Update(user)
}

func (s *DefaultUserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
