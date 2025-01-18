## Dependency Injection

# Part-1 : Understanding Repository and Service Layer Pattern

## Core Concept

The Repository and Service layer pattern is an architectural approach that separates data access logic from business logic. Think of it as a clear division of responsibilities:

```
[Database] <-> [Repository Layer] <-> [Service Layer] <-> [API/Handler Layer]
```

## Repository Layer: "How to Store and Retrieve"

The Repository layer is responsible for data persistence operations. It answers the question "How do I store and retrieve data?"

### Repository Layer Responsibilities:
- Direct database interactions (CRUD operations)
- Query construction
- Data mapping between database and domain models
- Transaction handling
- Cache interactions (if any)

```go
// Repository example
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    Update(user *User) error
    Delete(id uint) error
}
```

### Why Have a Repository Layer?
1. **Data Access Abstraction**
   - Hides database implementation details
   - Makes it easier to switch databases (e.g., PostgreSQL to MongoDB)
   - Centralizes data access logic

2. **Query Optimization**
   - Single place to optimize database queries
   - Easier to implement caching
   - Better control over database connections

## Service Layer: "What Should Happen"

The Service layer contains business logic. It answers the question "What should happen when...?"

### Service Layer Responsibilities:
- Business rules enforcement
- Data validation
- Orchestrating multiple operations
- Handling business exceptions
- Implementing business workflows

```go
// Service example
type UserService interface {
    RegisterUser(name, email string) error
    DeactivateAccount(userID uint) error
    UpdateProfile(userID uint, profile ProfileUpdate) error
}
```

### Why Have a Service Layer?
1. **Business Logic Encapsulation**
   - Keeps business rules in one place
   - Makes the codebase easier to understand
   - Prevents business logic from leaking into other layers

2. **Complex Operations Handling**
   - Can coordinate multiple repository operations
   - Handles transactional boundaries
   - Implements business workflows

## Real-World Example

Let's look at a user registration scenario:

```go
// Repository Layer: Handles data operations
type UserRepository interface {
    Create(user *User) error
    GetByEmail(email string) (*User, error)
}

// Service Layer: Handles business logic
type UserService interface {
    RegisterUser(name, email string, password string) error
}

// Service Implementation
func (s *userService) RegisterUser(name, email, password string) error {
    // Business Logic:
    // 1. Validate input
    if !isValidEmail(email) {
        return ErrInvalidEmail
    }

    // 2. Check if user exists
    existing, err := s.repo.GetByEmail(email)
    if err == nil && existing != nil {
        return ErrUserAlreadyExists
    }

    // 3. Hash password (business requirement)
    hashedPassword, err := hashPassword(password)
    if err != nil {
        return err
    }

    // 4. Create user through repository
    user := &User{
        Name: name,
        Email: email,
        Password: hashedPassword,
        Status: "active",
        CreatedAt: time.Now(),
    }
    
    return s.repo.Create(user)
}
```

## Benefits of This Separation

1. **Cleaner Code Organization**
   ```
   Without Separation:
   func CreateUser(name, email string) error {
       // Validation mixed with DB operations
       // Business logic mixed with SQL queries
       // Everything in one place
   }

   With Separation:
   // Repository: DB operations only
   func (r *repo) Create(user *User) error {
       return r.db.Create(user).Error
   }

   // Service: Business logic
   func (s *service) RegisterUser(name, email string) error {
       // Validation
       // Business rules
       // Calls repository when needed
   }
   ```

2. **Easier Testing**
   ```go
   // Mock repository for testing service
   type MockUserRepo struct {
       users map[string]*User
   }

   func TestUserRegistration(t *testing.T) {
       mockRepo := NewMockUserRepo()
       userService := NewUserService(mockRepo)
       
       // Test business logic without DB
       err := userService.RegisterUser("test", "test@email.com")
       assert.NoError(t, err)
   }
   ```

3. **Better Maintenance**
   - Changes to database schema only affect repository layer
   - Business rule changes only affect service layer
   - Less risk of breaking changes

## Common Use Cases

### 1. Data Validation and Business Rules
```go
func (s *service) CreateOrder(order Order) error {
    // Business rule: Check inventory
    if !s.inventoryService.HasStock(order.ProductID) {
        return ErrOutOfStock
    }
    
    // Business rule: Validate user credit
    if !s.paymentService.HasCredit(order.UserID) {
        return ErrInsufficientCredit
    }
    
    // Only after business rules pass, use repository
    return s.orderRepo.Create(order)
}
```

### 2. Complex Operations
```go
func (s *service) ProcessOrder(orderID uint) error {
    // Multiple repository operations in one business transaction
    order, err := s.orderRepo.GetByID(orderID)
    if err != nil {
        return err
    }
    
    // Update inventory
    err = s.inventoryRepo.DecrementStock(order.ProductID)
    if err != nil {
        return err
    }
    
    // Create payment record
    err = s.paymentRepo.Create(orderID)
    if err != nil {
        // Business logic: Rollback inventory if payment fails
        s.inventoryRepo.IncrementStock(order.ProductID)
        return err
    }
    
    return s.orderRepo.UpdateStatus(orderID, "processed")
}
```

### 3. Cross-Cutting Concerns
```go
func (s *service) UpdateUserProfile(userID uint, profile Profile) error {
    // Logging
    s.logger.Info("Updating profile", "userID", userID)
    
    // Authentication check
    if !s.authService.HasPermission(userID) {
        return ErrUnauthorized
    }
    
    // After all checks, use repository
    return s.userRepo.Update(userID, profile)
}
```

## When to Use This Pattern

Use this pattern when:
1. Your application has complex business rules
2. You need to maintain multiple data sources
3. You want to write testable code
4. You need to separate concerns clearly
5. Your application might grow in complexity

Don't use this pattern when:
1. Your application is very simple (CRUD only)
2. You have no business logic
3. You're building a quick prototype

## Common Pitfalls to Avoid

1. **Don't Put Business Logic in Repositories**
   ```go
   // BAD: Business logic in repository
   func (r *repo) CreateUser(user *User) error {
       if !isValidEmail(user.Email) { // Business validation shouldn't be here
           return ErrInvalidEmail
       }
       return r.db.Create(user)
   }

   // GOOD: Keep repositories focused on data operations
   func (r *repo) Create(user *User) error {
       return r.db.Create(user)
   }
   ```

2. **Don't Access Database Directly from Service**
   ```go
   // BAD: Service accessing database directly
   func (s *service) GetUser(id uint) (*User, error) {
       var user User
       return s.db.First(&user, id) // Direct DB access
   }

   // GOOD: Use repository
   func (s *service) GetUser(id uint) (*User, error) {
       return s.userRepo.GetByID(id)
   }
   ```

3. **Don't Skip Service Layer for Simple Operations**
   ```go
   // BAD: Skipping service layer
   func HandleGetUser(w http.ResponseWriter, r *http.Request) {
       user, err := userRepo.GetByID(id) // Directly using repo
   }

   // GOOD: Always go through service
   func HandleGetUser(w http.ResponseWriter, r *http.Request) {
       user, err := userService.GetUser(id)
   }
   ```

# Part-2 : Go Dependency Injection with GORM

This project demonstrates a clean, well-structured implementation of dependency injection in Go using GORM as the ORM for PostgreSQL. It follows best practices for package organization and separation of concerns.

## Project Structure

```
project/
├── models/
│   └── user.go           # Domain models
├── repository/
│   ├── interfaces.go     # Repository interfaces
│   └── gorm_repository.go # GORM implementation
├── service/
│   ├── interfaces.go     # Service interfaces
│   └── user_service.go   # Service implementation
├── config/
│   └── database.go       # Database configuration
├── tests/
│   └── user_service_test.go # Tests
└── main.go              # Application entry point
```

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- GORM

## Dependencies

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

## Setup

1. Clone the repository:
```bash
git clone https://github.com/giridharmb/depedency_injection.git depedency_injection
cd depedency_injection
go mod init depedency_injection
go mod tidy
```

2. Set up PostgreSQL database:
```sql
CREATE DATABASE testdb;
```

3. Update database configuration in `config/database.go`:
```go
dsn := "host=localhost user=postgres password=*** dbname=testdb port=5432 sslmode=disable"
```

4. Run the application:
```bash
go run main.go
```

## Running Tests

```bash
go test ./tests -v
```

## Package Details

### Models
Contains domain models that represent the business entities.
```go
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"not null"`
    Email string `gorm:"uniqueIndex;not null"`
}
```

### Repository
Implements data access layer using the repository pattern.

- `interfaces.go`: Defines repository interfaces
- `gorm_repository.go`: Implements interfaces using GORM

```go
type UserRepository interface {
    Create(user *models.User) error
    GetByID(id uint) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
}
```

### Service
Contains business logic and implements service interfaces.

- `interfaces.go`: Defines service interfaces
- `user_service.go`: Implements business logic

```go
type UserService interface {
    CreateUser(name, email string) error
    GetUser(id uint) (*models.User, error)
    UpdateUser(id uint, name, email string) error
    DeleteUser(id uint) error
}
```

### Config
Handles application configuration including database setup.

### Tests
Contains test implementations including mocks.

## Dependency Injection Benefits

1. **Loose Coupling**
   - Components are independent of their concrete implementations
   - Easy to switch implementations (e.g., different databases)
   - Better separation of concerns

2. **Testability**
   - Easy to mock dependencies for unit testing
   - Tests run without actual database connections
   - Isolated component testing

3. **Maintainability**
   - Clear dependency graph
   - Easy to modify individual components
   - Reduces code duplication

## Example Usage

```go
// Initialize database
db, err := config.InitDB()
if err != nil {
    log.Fatal(err)
}

// Initialize repository and service with dependency injection
userRepo := repository.NewGormUserRepository(db)
userService := service.NewUserService(userRepo)

// Create a new user
err = userService.CreateUser("John Doe", "john@example.com")
```

## Adding New Features

To add new features:

1. Add new models in the `models` package
2. Create new repository interfaces and implementations
3. Create new service interfaces and implementations
4. Inject dependencies through constructors

Example of adding a new feature:

```go
// models/product.go
type Product struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"not null"`
    Price float64
}

// repository/interfaces.go
type ProductRepository interface {
    Create(product *models.Product) error
    // ... other methods
}

// service/interfaces.go
type ProductService interface {
    CreateProduct(name string, price float64) error
    // ... other methods
}
```

## Testing

The project includes examples of mock implementations for testing:

```go
type MockUserRepository struct {
    users map[uint]*models.User
}

func TestUserService(t *testing.T) {
    mockRepo := NewMockUserRepository()
    userService := service.NewUserService(mockRepo)
    // ... test cases
}
```

## Best Practices

1. **Interface Segregation**
   - Keep interfaces focused and minimal
   - Split large interfaces into smaller ones

2. **Constructor Injection**
   - Use constructor injection over field injection
   - Make dependencies explicit

3. **Package Organization**
   - Clear separation of concerns
   - Logical grouping of related code
   - Clean dependency graph