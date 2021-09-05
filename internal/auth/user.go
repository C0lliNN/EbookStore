package auth

type UserRole string

const (
	Admin    UserRole = "ADMIN"
	Customer UserRole = "CUSTOMER"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Role      UserRole
	Password  string
	CreatedAt int64
}
