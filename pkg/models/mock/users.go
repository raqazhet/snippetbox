package mock

import (
	"alex/pkg/models"
	"time"
)

var mockUser = &models.User{
	Id:      1,
	Name:    "Arman",
	Email:   "aramn@gmail.com",
	Created: time.Now(),
}

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "raz@gmail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "arman@gmail.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentails
	}
}
func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
