package service

import (
	"crypto/sha1"
	"errors"
	"fmt"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/repository"
	"golang.org/x/exp/rand"
)

const (
	key  = "jfaopajsfojadsf"
	salt = "fkdsajl3214u98ujkj"
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	Guid string `json:"guid"`
	Ip   string
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repo: repos}
}

func (a *AuthService) CreateUser(user models.User) (string, error) {
	user.Password = a.generatePasswordHash(user.Password)
	return a.repo.CreateUser(user)
}

func (a *AuthService) GenerateTokens(guid, ip string) (string, string, error) {
	user, err := a.repo.GetUser(guid)
	if err != nil {
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Guid,
		ip,
	})

	accessToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := a.NewRefreshToken()

	return accessToken, refreshToken, err
}

func (a *AuthService) GetUserByGuid(guid string) (models.User, error) {
	return a.repo.GetUser(guid)
}

func (a *AuthService) ParseToken(acessToken string) (string, string, error) {
	token, err := jwt.ParseWithClaims(acessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(key), nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", "", errors.New("bad claims format")
	}
	return claims.Guid, claims.Ip, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	randVal := rand.NewSource(uint64(time.Now().Unix()))
	r := rand.New(randVal)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
