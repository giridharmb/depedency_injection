package main

import (
	"fmt"
	"log"

	"github.com/giridharmb/depedency_injection/config"
	"github.com/giridharmb/depedency_injection/repository"
	"github.com/giridharmb/depedency_injection/service"
)

func main() {
	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repository and service with dependency injection
	userRepo := repository.NewGormUserRepository(db)
	userService := service.NewUserService(userRepo)

	// Example usage
	err = userService.CreateUser("John Doe", "john@example.com")
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return
	}

	user, err := userService.GetUser(1)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return
	}

	fmt.Printf("Found user: %+v\n", user)
}
