package helpers

import (
	"errors"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	SaveUser(*models.User) (int64, error)
	GetUserById(*models.User, uint) error
	GetUserByUsername(*models.User, string) error
	GetUserBySubId(*models.User, string) error
	GetAllUsers(*[]models.User) error
}

func (m *MainRepository) SaveUser(user *models.User) (int64, error) {
	return m.DB.Save(user).RowsAffected, m.DB.Error
}

func (m *MainRepository) GetUserById(user *models.User, u uint) error {
	if u <= 0 {
		return errors.New("invalid userID")
	}
	m.DB.Preload(clause.Associations).Find(user, u)
	return m.DB.Error
}

func (m *MainRepository) GetUserByUsername(user *models.User, s string) error {
	if len(s) == 0 {
		return errors.New("invalid username")
	}
	m.DB.Preload(clause.Associations).Where("username = ?", s).First(user)
	return m.DB.Error
}

func (m *MainRepository) GetUserBySubId(user *models.User, id string) error {
	if len(id) == 0 {
		return errors.New("invalid subscriber id")
	} // SubId
	m.DB.Preload(clause.Associations).Where("sub_id = ?", id).First(user)
	return m.DB.Error
}

func (m *MainRepository) GetAllUsers(users *[]models.User) error {
	m.DB.Preload(clause.Associations).Find(users)
	return m.DB.Error
}
