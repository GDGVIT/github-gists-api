package user

type Service interface {
	Register(user *User) (*User, error)

	GetUserByID(id float64) (*User, error)

	GetRepo() Repository
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) Register(user *User) (*User, error) {
	exists, err := s.repo.DoesEmailExist(user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		//noinspection GoErrorStringFormat
		u, err := s.repo.FindByEmail(user.Email)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	return s.repo.Register(user)
}

func (s *service) GetUserByID(id float64) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *service) GetRepo() Repository {
	return s.repo
}
