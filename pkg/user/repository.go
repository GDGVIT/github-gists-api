package user

type Repository interface {
	FindByID(id float64) (*User, error)

	FindByEmail(email string) (*User, error)

	Register(user *User) (*User, error)

	DoesEmailExist(email string) (bool, error)
}
