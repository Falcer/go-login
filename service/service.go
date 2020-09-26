package service

import (
	"errors"
	"log"
	"time"

	"github.com/Falcer/go-login/model"
	"github.com/Falcer/go-login/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/crypto/bcrypt"
)

type (
	// UserService interface
	UserService interface {
		Users() ([]model.User, error)
		Login(user model.User) (*string, error)
		Register(user model.User) error
		Hash(password string) (string, error)
		CheckPasswordHash(password, hash string) bool
		GenerateJWT(username string) (*string, error)
		CheckJWT(tokenString string) bool
	}
	// UserServiceImpl struct
	UserServiceImpl struct {
		repo repository.UserRepository
	}
	userConfig struct {
		JWTKEY string `envconfig:"JWT_KEY"`
	}
	// Claims struct
	Claims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}
)

// NewUserService function
func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo}
}

// Users methods
func (s *UserServiceImpl) Users() ([]model.User, error) {
	return s.repo.GetAllUser()
}

// Login Methods
func (s *UserServiceImpl) Login(user model.User) (*string, error) {
	userRes, err := s.repo.GetUserByUsername(user.Username)
	if err != nil {
		return nil, err
	}

	if !s.CheckPasswordHash(user.Password, userRes.Password) {
		return nil, errors.New("Wrong Password")
	}

	token, err := s.GenerateJWT(user.Username)

	if !s.CheckPasswordHash(user.Password, userRes.Password) {
		return nil, errors.New("Generated Token Failed")
	}

	return token, nil

}

// Register methods
func (s *UserServiceImpl) Register(user model.User) error {
	hash, err := s.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	return s.repo.AddUser(user)
}

// Hash methods
func (s *UserServiceImpl) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash methods
func (s *UserServiceImpl) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT methods
func (s *UserServiceImpl) GenerateJWT(username string) (*string, error) {

	var env userConfig

	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err)
	}

	expTime := time.Now().Add(30 * 24 * time.Hour)
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{Username: username, StandardClaims: jwt.StandardClaims{ExpiresAt: expTime.Unix()}}).SignedString([]byte(env.JWTKEY))

	if err != nil {
		return nil, err
	}

	return &tokenString, nil

}

// CheckJWT methods
func (s *UserServiceImpl) CheckJWT(tokenString string) bool {

	var env userConfig

	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err)
	}

	claims := &Claims{}
	tokenParse, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.JWTKEY), nil
	})

	if err != nil || !tokenParse.Valid {
		return false
	}

	return true

}
