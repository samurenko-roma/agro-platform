package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samurenkoroma/agro-platform/internal/domain/account/aggregate/user"
	configs "github.com/samurenkoroma/agro-platform/internal/shared/config"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Claims структура JWT токена (с OrganizationID)
type Claims struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
	OrgRole        string `json:"org_role"`
	TokenType      string `json:"token_type"` // "access" | "refresh" — различает, где можно использовать токен
	jwt.RegisteredClaims
}

// TokenPair пара токенов
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

// Service JWT сервис
type Service struct {
	config configs.AuthConfig
}

func NewService(config configs.AuthConfig) *Service {
	return &Service{config: config}
}

// GenerateTokenPair генерирует пару токенов с organizationID
func (s *Service) GenerateTokenPair(userID, username, email, globalRole, organizationID, orgRole string) (*TokenPair, error) {
	accessToken, err := s.generateAccessToken(userID, username, email, globalRole, organizationID, orgRole)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(userID, username, email, globalRole, organizationID, orgRole)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessExpiry.Seconds()),
	}, nil
}

func (s *Service) generateAccessToken(userID, username, email, globalRole, organizationID, orgRole string) (string, error) {
	return s.generateToken(userID, username, email, globalRole, organizationID, orgRole, TokenTypeAccess, s.config.AccessExpiry)
}

func (s *Service) generateRefreshToken(userID, username, email, globalRole, organizationID, orgRole string) (string, error) {
	return s.generateToken(userID, username, email, globalRole, organizationID, orgRole, TokenTypeRefresh, s.config.RefreshExpiry)
}

func (s *Service) generateToken(userID, username, email, globalRole, organizationID, orgRole, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:         userID,
		Username:       username,
		Email:          email,
		Role:           globalRole,
		OrganizationID: organizationID,
		OrgRole:        orgRole,
		TokenType:      tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    s.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

// ValidateToken валидирует токен и возвращает claims — БЕЗ проверки типа.
// Используй ValidateAccessToken/ValidateRefreshToken там, где важно точно
// знать, какой именно токен пришёл (везде, кроме внутренних служебных
// случаев).
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, user.ErrTokenExpired
		}
		return nil, user.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, user.ErrInvalidToken
}

// ValidateAccessToken — валидация для защищённых API-эндпоинтов
// (AuthMiddleware). Отклоняет refresh-токен, даже если подпись валидна.
func (s *Service) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != TokenTypeAccess {
		return nil, user.ErrInvalidToken
	}
	return claims, nil
}

// ValidateRefreshToken — валидация для /auth/refresh. Отклоняет access-
// токен, даже если подпись валидна.
func (s *Service) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != TokenTypeRefresh {
		return nil, user.ErrInvalidToken
	}
	return claims, nil
}

// RefreshToken обновляет пару токенов по refresh-токену
func (s *Service) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return s.GenerateTokenPair(
		claims.UserID,
		claims.Username,
		claims.Email,
		claims.Role,
		claims.OrganizationID,
		claims.OrgRole,
	)
}
