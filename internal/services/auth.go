package services

import (
	errs "cvwo-backend/internal/errors"
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"

	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repos.UserRepo
}

func NewAuthService(userRepo *repos.UserRepo) *AuthService {
	return &AuthService{userRepo}
}

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

// Given a username and password, check if the password matches and if so, generate a jwt
func (service *AuthService) Authenticate(authInput *models.AuthInput) (*models.User, string, error) {
	// Check if there is any user with that username
	user, err := service.userRepo.GetByUsername(authInput.Username)
	if err != nil {
		return nil, "", errs.New(errs.ErrUnauthorized, "Username not found")
	}

	// Check if password matches
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authInput.Password)); err != nil {
		return nil, "", errs.New(errs.ErrUnauthorized, "Incorrect password")
	}

	// Set jwt claims including userID and expiration time
	claims := &jwt.RegisteredClaims{
		ID:        strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, "", err
	}

	return user, signedToken, err
}

func (service *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	// Parse jwt while checking for correct claims type and signing method	
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Handle invalid token
	if err != nil || !token.Valid {
		return nil, errs.New(errs.ErrUnauthorized, "Invalid token")
	}

	// Check if token has expired
	if claims.ExpiresAt.Compare(time.Now()) == -1 {
		return nil, errs.New(errs.ErrUnauthorized, "Expired token")
	}

	// Convert token ID to int
	userId, err := strconv.Atoi(claims.ID)
	if err != nil {
		return nil, errs.New(errs.ErrUnauthorized, fmt.Sprintf("Invalid user id: %s", claims.ID))
	}

	// Retrieve the user whose ID corresponds to the token ID
	user, err := service.userRepo.GetByID(uint(userId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.ErrUnauthorized, "User not found")
		}
		return nil, err
	}

	return user, nil
}
