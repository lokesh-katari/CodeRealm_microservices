// repositories.go
package authservice

import (
	"database/sql"
	"fmt"
	"time"

	"lokesh-katari/code-realm/cmd/auth/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	// import any other required packages
)

type User struct {
	ID       int
	Email    string
	Password string
	Username string
}

// UserRepository defines the interface for user data access operations
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)

	// Add other repository methods as needed
}

//generate hash password for the user before saving to the database

func (u *User) GenerateHashPassword() error {
	// hash the password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	fmt.Println("hashedPassword", string(hashedPassword))
	u.Password = string(hashedPassword)
	return nil
}

//Generate the JWT Token

func (u *User) GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["email"] = u.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// check the password is valid or not
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// PostgresUserRepository implements the UserRepository using PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository() (*PostgresUserRepository, error) {
	db, err := db.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	return &PostgresUserRepository{db: db}, nil
}

func (r *PostgresUserRepository) CreateUser(user *User) error {

	user.GenerateHashPassword()
	query := "INSERT INTO users (email, password, name) VALUES ($1, $2, $3)"
	re, err := r.db.Exec(query, user.Email, user.Password, user.Username)
	fmt.Println(re)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, email, password, name FROM users WHERE email = $1"
	row := r.db.QueryRow(query, email)
	var user User

	// Scan the query result into the User object
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Username)
	if err != nil {
		// Check if the error is due to no rows being returned
		if err == sql.ErrNoRows {
			return nil, nil // User not found, return nil without error
		}
		// Return any other error encountered during scanning
		return nil, err
	}

	// Return the User object
	return &user, nil
}
