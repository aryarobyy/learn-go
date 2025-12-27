package model

import "github.com/golang-jwt/jwt/v5"

type Logincredential struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type RegisterCredential struct {
	Name     string `db:"name" json:"name"`
	Username string `db:"username" json:"username" validate:"required"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password" validate:"required,min=8"`
}

type VerifyEmail struct {
	Email string `db:"email" json:"email"`
	ID    int    `db:"id" json:"id"`
}

type VerifyUsername struct {
	Username string `db:"username" json:"username"`
	ID       int    `db:"id" json:"id"`
}

type ClaimsModel struct {
	UserID   int    `json:"id"`
	Role     Role   `json:"role"`
	Name     string `json:"name"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
