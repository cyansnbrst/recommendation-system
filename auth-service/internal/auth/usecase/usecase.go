package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"cyansnbrst/auth-service/config"
	"cyansnbrst/auth-service/internal/auth"
	"cyansnbrst/auth-service/internal/models"
)

// Auth usecase struct
type authUC struct {
	cfg      *config.Config
	authRepo auth.Repository
	logger   *zap.Logger
}

// Auth usecase constructor
func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, logger *zap.Logger) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, logger: logger}
}

// Custom JWT claims
type CustomClaims struct {
	UserUID string `json:"user_uid"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// Create a user
func (u *authUC) Create(email, password string) (string, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("failed to hash password", zap.Error(err))
		return "", "", err
	}

	newUser := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	newUID, err := u.authRepo.Insert(newUser)
	if err != nil {
		u.logger.Error("failed to create user", zap.Error(err))
		return "", "", err
	}

	token, err := u.GenerateJWT(*newUser)
	if err != nil {
		u.logger.Error("failed to generate token",
			zap.String("error", err.Error()),
		)
		return "", "", err
	}

	return token, newUID, nil
}

// Generate JWT token
func (u *authUC) GenerateJWT(user models.User) (string, error) {
	claims := CustomClaims{
		UserUID: user.ID,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.cfg.Timeout.Token)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.cfg.SecretKey))
}

// Validate JWT token
func (u *authUC) ValidateToken(tokenString string) (string, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(u.cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return "", false, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return "", false, errors.New("invalid claims structure")
	}

	return claims.UserUID, claims.IsAdmin, nil
}

// Validate  user's credentials
func (u *authUC) ValidateCredentials(email, password string) (*models.User, error) {
	user, err := u.authRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// Send request for create user's profile
func (u *authUC) CreateProfile(uid string, name string, cookie *http.Cookie) error {
	profileRequestBody := models.CreateProfileDTO{Name: name}

	data, err := json.Marshal(profileRequestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal profile request: %w", err)
	}

	url := fmt.Sprintf("%s/%s", u.cfg.ProfileLink, uid)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send profile creation request: %w", err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			u.logger.Error("error closing response body", zap.Error(closeErr))
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create profile, status: %s", resp.Status)
	}

	return nil
}
