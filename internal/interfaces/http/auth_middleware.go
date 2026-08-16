package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/samurenkoroma/agro-platform/internal/domain/account/aggregate/organization"
	"github.com/samurenkoroma/agro-platform/internal/domain/account/aggregate/user"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/jwt"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

type AuthMiddleware struct {
	jwtService *jwt.Service
}

func NewAuthMiddleware(jwtService *jwt.Service) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

// Authenticate проверяет токен и добавляет все данные в контекст
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WriteUnauthorized(w, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.WriteUnauthorized(w, "invalid authorization format")
			return
		}

		token := parts[1]

		claims, err := m.jwtService.ValidateAccessToken(token) // ← было ValidateToken
		if err != nil {
			response.WriteUnauthorized(w, err.Error())
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "organization_id", claims.OrganizationID)
		ctx = context.WithValue(ctx, "org_role", organization.OrganizationRole(claims.OrgRole))
		ctx = context.WithValue(ctx, "role", user.Role(claims.Role))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetOrganizationID извлекает ID текущей организации из контекста
func GetOrganizationID(ctx context.Context) string {
	if orgID, ok := ctx.Value("organization_id").(string); ok {
		return orgID
	}
	return ""
}

// GetOrgRole извлекает роль в текущей организации из контекста
func GetOrgRole(ctx context.Context) organization.OrganizationRole {
	if role, ok := ctx.Value("org_role").(organization.OrganizationRole); ok {
		return role
	}
	return ""
}
