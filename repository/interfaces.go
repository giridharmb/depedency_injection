package repository

import "dependency_injection/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}
