package models

import (
	"github.com/markbates/goth"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string        `json:"username"`
	SubId     string        `json:"-"` // This is the ID from the OIDC and we don't want that sent back and forth
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	UserRole  []UserRole    `json:"userRole"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Addresses []UserAddress `json:"addresses"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *User) GetTruncatedUserRoles() []TruncatedUserRole {
	truncated := make([]TruncatedUserRole, 0, len(u.UserRole))
	for _, role := range u.UserRole {
		if !role.DeletedAt.Valid {
			truncated = append(truncated, role.TruncateRole())
		}
	}
	return truncated
}

func UserFromGormUser(u *User, g *goth.User) {
	u.Username = g.Email
	u.SubId = g.UserID
	u.FirstName = g.FirstName
	u.LastName = g.LastName
}
