package domain

import "database/sql"

type UserRole string

const (
	RoleAdmin UserRole = "Administrator"
	RoleUser  UserRole = "User"
)

type User struct {
	ID           int64        `db:"id"`
	FirstName    string       `db:"first_name"`
	LastName     string       `db:"last_name"`
	Email        string       `db:"email"`
	PasswordHash string       `db:"password_hash"`
	Role         UserRole     `db:"role"`
	CreatedAt    sql.NullTime `db:"created_at"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

type UserRequest struct {
	FirstName string `form:"first_name" validate:"required"`
	LastName  string `form:"last_name" validate:"required"`
	Email     string `form:"email" validate:"required,email"`
	Password  string `form:"password" validate:"required"`
}

func (ur *UserRequest) ToDomain() *User {
	return &User{
		FirstName: ur.FirstName,
		LastName:  ur.LastName,
		Email:     ur.Email,
		Role:      RoleUser,
	}
}
