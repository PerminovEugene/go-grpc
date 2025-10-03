package models

type RatingCategory struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Weight int    `json:"weight" db:"weight"`
}

// type CreateRatingCategoryRequest struct {
// 	Name   string `json:"name"`
// 	Weight int    `json:"weight"`
// }

// type UpdateRatingCategoryRequest struct {
// 	ID     int    `json:"id"`
// 	Name   string `json:"name"`
// 	Weight int    `json:"weight"`
// }
