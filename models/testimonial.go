package models

import (
	"gorm.io/gorm"
)

//type ReviewJson struct {
//	ID        uint `gorm:"primarykey"`
//	CreatedAt time.Time
//	UpdatedAt time.Time
//	DeletedAt gorm.DeletedAt `gorm:"index"`
//	Header    string
//	Body      string
//	Images    []string
//}
//
//func (j *ReviewJson) GetString() string {
//	return strings.Join(j.Images, ";")
//}
//
//func (j *ReviewJson) ToRegularVersion(r *Testimonial) {
//	r.ID = j.ID
//	r.CreatedAt = j.CreatedAt
//	r.UpdatedAt = j.UpdatedAt
//	r.DeletedAt = j.DeletedAt
//	r.Header = j.Header
//	r.Body = j.Body
//	r.Images = j.GetString()
//}

type Testimonial struct {
	gorm.Model
	Restaurant   *Restaurant `gorm:"foreignKey:ID"`
	Header       string      `json:"header"`
	Body         string      `json:"body"`
	ImageUrls    []string    `json:"imageUrls"`
	RestaurantId uint        `json:"restaurantId"`
	//https://stackoverflow.com/questions/64035165/unsupported-data-type-error-on-gorm-field-where-custom-valuer-returns-nil
}

type TestimonialImage struct {
	gorm.Model
	URL           string `json:"URL"`
	TestimonialId uint   `json:"TestimonialId"`
}

//func (r *Testimonial) ToJsonVersion(j *ReviewJson) {
//	j.ID = r.ID
//	j.CreatedAt = r.CreatedAt
//	j.UpdatedAt = r.UpdatedAt
//	j.DeletedAt = r.DeletedAt
//	j.Header = r.Header
//	j.Body = r.Body
//	j.Images = r.GetStringArr()
//}
//
//func (r *Testimonial) AppendString(s string) {
//	if len(r.Images) == 0 {
//		r.Images = s
//		return
//	}
//	r.Images = r.Images + ";" + s
//}
//
//func (r *Testimonial) GetStringArr() []string {
//	return strings.Split(r.Images, ";")
//}
