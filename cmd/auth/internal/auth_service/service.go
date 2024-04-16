// service.go
package authservice

import (
	"context"
	"errors"
	"log"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	RegisterUser(ctx context.Context, name, email, password string) (*User, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
	LogoutUser(ctx context.Context, token string) error
	GetUser(ctx context.Context, token string) (*User, error)
}

type AuthServiceImpl struct {
	userRepo   UserRepository
	jwtManager *JWTManager
}

// NewAuthServiceImpl creates a new AuthServiceImpl
func NewAuthServiceImpl(userRepo UserRepository, jwtManager *JWTManager) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo: userRepo, jwtManager: jwtManager}
}

// RegisterUser handles user registration
func (s *AuthServiceImpl) RegisterUser(ctx context.Context, name, email, password string) (*User, error) {
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
		Password: password,
		Name:     name, // Hash the password before storing
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

func (s *AuthServiceImpl) GetUser(ctx context.Context, token string) (*User, error) {
	// Verify the JWT token
	log.Println("this is token from get user", token)
	claims, err := s.jwtManager.VerifyJWT(token)
	if err != nil {
		return nil, err
	}

	// Get the user ID from the claims
	// userID, ok := claims["user_id"].(string)
	userEmail, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid user ID")
	}

	// Get the user from the repository
	user, err := s.userRepo.GetUserByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Add other service method implementations as needed
