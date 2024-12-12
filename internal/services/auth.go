package services

import (
	"cvwo-backend/internal/models"

	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{userService}
}

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

// Given a username and password, check if the password matches and if so, generate a jwt
func (service *AuthService) Authenticate(authInput *models.AuthInput) (string, error) {
	user, err := service.userService.repo.GetByUsername(authInput.Username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authInput.Password)); err != nil {
		return "", NewUnauthorizedError("incorrect password")
	}

	// Set jwt claims including userID and expiration time
	claims := jwt.RegisteredClaims{
		ID: string(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, err
}

func (service *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check that the token's signing method/algorithm matches
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, NewUnauthorizedError("invalid or expired token")
	}

	// Check if the token's claims match the RegisteredClaims type
	claims, ok := token.Claims.(jwt.RegisteredClaims); 
	if !ok {
		return nil, NewUnauthorizedError("invalid token")
	}

	// Check if token has expired
	if claims.ExpiresAt.Compare(time.Now()) == -1 {
		return nil, NewUnauthorizedError("expired token")
	}

	// Convert token ID to int
	userId, err := strconv.Atoi(claims.ID)
	if err != nil {
		return nil, NewUnauthorizedError("invalid user id")
	}

	// Retrieve the user whose ID corresponds to the token ID
	user, err := service.userService.GetByID(uint(userId))
	if err != nil {
		// If no user corresponds to given ID, return unauthorized error
		var notFoundErr *NotFoundError
		if errors.As(err, &notFoundErr) {
			return nil, NewUnauthorizedError("invalid user id")
		}
		return nil, err
	}

	return user, nil
}
