package models

import "gorm.io/gorm"

type Testimonial struct {
	gorm.Model
	//ID      int32  `json:"id"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Rating  Rating `json:"rating"`
}

type Rating int8

func OneStarRating() Rating {
	return 1
}
func TwoStarRating() Rating {
	return 2
}
func ThreeStarRating() Rating {
	return 3
}
func FourStarRating() Rating {
	return 4
}
func FiveStarRating() Rating {
	return 5
}
