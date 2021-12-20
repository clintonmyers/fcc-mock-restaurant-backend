package helpers

import (
	"errors"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
)

type UserRepository interface {
	SaveUser(*models.User) (int64, error)
	GetUserById(*models.User, uint) error
	GetUserByUsername(*models.User, string) error
	GetAllUsers(*[]models.User) error
}

func (m *MainRepository) SaveUser(user *models.User) (int64, error) {
	return m.DB.Save(user).RowsAffected, m.DB.Error
}

func (m *MainRepository) GetUserById(user *models.User, u uint) error {
	if u <= 0 {
		return errors.New("invalid userID")
	}
	m.DB.Find(user, u)
	return m.DB.Error
}

func (m *MainRepository) GetUserByUsername(user *models.User, s string) error {
	if len(s) == 0 {
		return errors.New("invalid username")
	}
	m.DB.Where("username = ?", s).First(user)
	return m.DB.Error
}

func (m *MainRepository) GetAllUsers(users *[]models.User) error {
	m.DB.Find(users)
	return m.DB.Error
}
