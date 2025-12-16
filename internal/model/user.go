package model

import "time"

type User struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Username  string     `db:"username" json:"username"`
	Password  string     `db:"password" json:"password"`
	Role      Role       `db:"role" json:"role"`
	Age       *int       `db:"age" json:"age"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)
