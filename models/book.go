package models

type BookDetail struct {
	Title           string `json:"title"`
	AvailableCopies int    `json:"available_copies"`
}
