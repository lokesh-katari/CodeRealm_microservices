// service.go
package authservice

import (
	"context"
	"errors"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	RegisterUser(ctx context.Context, email, password string) (*User, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
	LogoutUser(ctx context.Context, token string) error
}

type AuthServiceImpl struct {
	userRepo UserRepository
}

// NewAuthServiceImpl creates a new AuthServiceImpl
func NewAuthServiceImpl(userRepo UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo: userRepo}
}

// RegisterUser handles user registration
func (s *AuthServiceImpl) RegisterUser(ctx context.Context, email, password string) (*User, error) {
	// Check if the user already exists
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Create a new user
	newUser := &User{
		Email:    email,
		Password: (password), // Hash the password before storing
		// Add other user fields as needed
	}

	// Save the new user to the repository
	err = s.userRepo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// LoginUser handles user login
func (s *AuthServiceImpl) LoginUser(ctx context.Context, email, password string) (string, error) {
	// Get the user from the repository
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Compare the provided password with the user's password
	if user.ComparePassword(password) {
		// Generate a JWT token for the user
		token, err := user.GenerateJWT()
		if err != nil {
			return "", err
		}
		return token, nil
	}

	return "", errors.New("invalid password")
}

func (s *AuthServiceImpl) LogoutUser(ctx context.Context, token string) error {

	return nil
}

// Add other service method implementations as needed
