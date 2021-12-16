package models

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

type ReviewJson struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Header    string
	Body      string
	Images    []string
}

func (j *ReviewJson) GetString() string {
	return strings.Join(j.Images, ";")
}

func (j *ReviewJson) ToRegularVersion(r *Review) {
	r.ID = j.ID
	r.CreatedAt = j.CreatedAt
	r.UpdatedAt = j.UpdatedAt
	r.DeletedAt = j.DeletedAt
	r.Header = j.Header
	r.Body = j.Body
	r.Images = j.GetString()
}

type Review struct {
	gorm.Model
	Restaurant *Restaurant `gorm:"foreignKey:ID"`
	Header     string
	Body       string
	Images     string
	//https://stackoverflow.com/questions/64035165/unsupported-data-type-error-on-gorm-field-where-custom-valuer-returns-nil
}

func (r *Review) ToJsonVersion(j *ReviewJson) {
	j.ID = r.ID
	j.CreatedAt = r.CreatedAt
	j.UpdatedAt = r.UpdatedAt
	j.DeletedAt = r.DeletedAt
	j.Header = r.Header
	j.Body = r.Body
	j.Images = r.GetStringArr()
}
func (r *Review) AppendString(s string) {
	if len(r.Images) == 0 {
		r.Images = s
		return
	}
	r.Images = r.Images + ";" + s
}

func (r *Review) GetStringArr() []string {
	return strings.Split(r.Images, ";")
}
