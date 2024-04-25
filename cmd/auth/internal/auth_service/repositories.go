// repositories.go
package authservice

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	// "strings"
	"time"

	"lokesh-katari/code-realm/cmd/auth/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	// import any other required packages
)

type User struct {
	ID                   int
	Email                string
	Password             string
	Name                 string
	Easy_Problem_count   int
	Medium_Problem_count int
	Hard_Problem_count   int
	Submission           []string
}

// UserRepository defines the interface for user data access operations
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	UpdateUserSubmissions(user *User, queId string, difficulty string) error

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
	log.Println("NewPostgresUserRepository")
	db, err := db.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	return &PostgresUserRepository{db: db}, nil
}

func (r *PostgresUserRepository) CreateUser(user *User) error {

	user.GenerateHashPassword()
	query := "INSERT INTO users (email, password, name) VALUES ($1, $2, $3)"
	re, err := r.db.Exec(query, user.Email, user.Password, user.Name)
	fmt.Println(re)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, email, password, name, easy_problem_count, medium_problem_count, hard_problem_count, submissions FROM users WHERE email = $1"
	row := r.db.QueryRow(query, email)
	var user User
	var submissionBytes []byte

	// Scan the query result into the User object
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Easy_Problem_count, &user.Medium_Problem_count, &user.Hard_Problem_count, &submissionBytes)
	if err != nil {
		// Check if the error is due to no rows being returned
		if err == sql.ErrNoRows {
			return nil, nil // User not found, return nil without error
		}
		// Return any other error encountered during scanning
		return nil, err
	}
	submissionStr := strings.Trim(string(submissionBytes), "{}")

	// Split the string by comma to get individual elements
	user.Submission = strings.Split(submissionStr, ",")

	fmt.Println("user", user.Submission)
	// Return the User object
	return &user, nil
}

func (r *PostgresUserRepository) UpdateUser(user *User) error {
	query := "UPDATE users SET email = $1, password = $2, name = $3 WHERE id = $4"
	// submission := strings.Join(user.Submission, ",")
	_, err := r.db.Exec(query, user.Email, user.Password, user.Name, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) UpdateUserSubmissions(user *User, queId string, difficulty string) error {
	query := `
	    UPDATE users
	    SET submissions = CASE
	        WHEN $1 = ANY(submissions) THEN submissions
	        ELSE array_append(submissions, $1)
	    END,
		hard_problem_count = CASE
		WHEN $3 = 'hard' AND $1 <> ALL(submissions) THEN hard_problem_count + 1
		ELSE hard_problem_count
	END,
	medium_problem_count = CASE
		WHEN $3 = 'medium' AND $1 <> ALL(submissions) THEN medium_problem_count + 1
		ELSE medium_problem_count
	END,
	easy_problem_count = CASE
		WHEN $3 = 'easy' AND $1 <> ALL(submissions) THEN easy_problem_count + 1
		ELSE easy_problem_count
	END
	    WHERE id = $2;
	`
	fmt.Println(difficulty, queId, user.ID)
	_, err := r.db.Exec(query, queId, user.ID, difficulty)
	if err != nil {
		return err
	}
	fmt.Println("success fully updated the user submission")

	return nil
}
