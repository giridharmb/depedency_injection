package service

import "dependency_injection/models"

type UserService interface {
	CreateUser(name, email string) error
	GetUser(id uint) (*models.User, error)
	UpdateUser(id uint, name, email string) error
	DeleteUser(id uint) error
}
