package service

import (
	"crypto/sha1"
	"errors"
	"fmt"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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
	Guid string
	Ip   string
	Key  string
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repo: repos}
}

func (a *AuthService) CreateUser(user models.User) (string, error) {
	user.Guid, _ = a.generateRandSeq()
	logrus.Info(user.Guid)
	user.Password = a.generatePasswordHash(user.Password)
	return a.repo.CreateUser(user)
}

func (a *AuthService) GenerateTokens(guid, ip string) (string, string, error) {
	user, err := a.repo.GetUser(guid)
	if err != nil {
		return "", "", err
	}

	pair_key, _ := a.generateRandSeq()

	accessToken, err := a.newJWT(guid, ip, pair_key, 12*time.Hour)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := a.newJWT(guid, ip, pair_key, 1000*time.Hour)
	if err != nil {
		return "", "", err
	}

	hash, err := a.hashRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	logrus.Info("check")

	if _, err := a.repo.SaveRefreshToken(user.Guid, hash); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

func (a *AuthService) newJWT(guid, ip, pair_key string, expTime time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		guid,
		ip,
		pair_key,
	})

	return token.SignedString([]byte(key))
}

func (a *AuthService) GetUserByGuid(guid string) (models.User, error) {
	return a.repo.GetUser(guid)
}

func (a *AuthService) parseToken(acessToken string) (string, string, string, error) {
	token, err := jwt.ParseWithClaims(acessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(key), nil
	})

	if err != nil {
		return "", "", "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", "", "", errors.New("bad claims format")
	}
	return claims.Guid, claims.Ip, claims.Key, nil
}

func (s *AuthService) Refresh(accessToken, refreshToken, ip string) (string, string, error) {
	guid, lastIp, pair_key1, err := s.parseToken(accessToken)

	if err != nil {
		return "", "", err
	}

	tokens, err := s.repo.GetUserTokens(guid)
	if err != nil {
		return "", "", err
	}

	exists := false
	var tokenId int

	for _, token := range tokens {
		if s.checkEqHash(token.TokenHash, refreshToken) {
			exists = true
			tokenId = token.Id
			break
		}
	}

	if !exists {
		return "", "", errors.New("unknown refresh token")
	}

	_, _, pair_key2, err := s.parseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	if pair_key1 != pair_key2 {
		return "", "", errors.New("incorrect token pair")
	}

	if ip != lastIp {
		logrus.Printf("send warning email to user")
	}

	s.repo.RemoveToken(tokenId)

	return s.GenerateTokens(guid, ip)
}

func (s *AuthService) hashRefreshToken(refreshToken string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken[:72]), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *AuthService) checkEqHash(hash, refreshToken string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(refreshToken[:72]))
	return err == nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) generateRandSeq() (string, error) {
	b := make([]byte, 32)

	randVal := rand.NewSource(uint64(time.Now().Unix()))
	r := rand.New(randVal)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
