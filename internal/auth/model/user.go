package model

type UserRole string

const (
	Admin    UserRole = "ADMIN"
	Customer UserRole = "CUSTOMER"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email    string
	Role     UserRole
	Password string
	CreatedAt int64
}

func (u User) IsAdmin() bool {
	return u.Role == Admin
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}
