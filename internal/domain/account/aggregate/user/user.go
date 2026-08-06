package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samurenkoroma/agro-platform/internal/domain/shared/aggregate"
	"golang.org/x/crypto/bcrypt"
)

// UserStatus статус пользователя
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBlocked  UserStatus = "blocked"
)

// User пользователь системы
type User struct {
	aggregate.BaseAggregate
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`

	// Текущая активная организация
	CurrentOrganizationID *string `json:"current_organization_id,omitempty"`

	Role      Role       `json:"role"`
	Status    UserStatus `json:"status"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// NewUser создает нового пользователя
func NewUser(email, username, password, firstName, lastName, phone string) (*User, error) {
	if email == "" {
		return nil, ErrEmailRequired
	}
	if username == "" {
		return nil, ErrUsernameRequired
	}
	if password == "" || len(password) < 6 {
		return nil, ErrPasswordTooShort
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		Role:      RoleClient,
		ID:        uuid.New().String(),
		Email:     email,
		Username:  username,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Status:    UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// SetCurrentOrganization устанавливает текущую организацию
func (u *User) SetCurrentOrganization(orgID string) {
	u.CurrentOrganizationID = &orgID
	u.UpdatedAt = time.Now()
}

// ClearCurrentOrganization очищает текущую организацию
func (u *User) ClearCurrentOrganization() {
	u.CurrentOrganizationID = nil
	u.UpdatedAt = time.Now()
}

// GetCurrentOrganizationID возвращает ID текущей организации
func (u *User) GetCurrentOrganizationID() string {
	if u.CurrentOrganizationID == nil {
		return ""
	}
	return *u.CurrentOrganizationID
}

// hashPassword хеширует пароль
func hashPassword(password string) (string, error) {
	passwordBytes := []byte(password)

	// Дополнительная проверка, что преобразование прошло успешно
	if passwordBytes == nil {
		err := errors.New("failed to convert password to bytes")
		fmt.Print(err)
		return "", err
	}
	bytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword проверяет пароль
func (u *User) CheckPassword(password string) bool {
	passwordBytes := []byte(password)

	// Дополнительная проверка, что преобразование прошло успешно
	if passwordBytes == nil {
		fmt.Print(errors.New("failed to convert password to bytes"))
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), passwordBytes)
	return err == nil
}

// ChangePassword меняет пароль
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if !u.CheckPassword(oldPassword) {
		return ErrInvalidPassword
	}
	if len(newPassword) < 6 {
		return ErrPasswordTooShort
	}

	hashed, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	u.Password = hashed
	u.UpdatedAt = time.Now()
	return nil
}

// Activate активирует пользователя
func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
}

// Deactivate деактивирует пользователя
func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now()
}

// Block блокирует пользователя
func (u *User) Block() {
	u.Status = UserStatusBlocked
	u.UpdatedAt = time.Now()
}

// UpdateLastLogin обновляет время последнего входа
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

// IsActive проверяет активен ли пользователь
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// HasRole проверяет наличие роли
func (u *User) HasRole(role Role) bool {
	return u.Role == role || u.Role == RoleAdmin
}
