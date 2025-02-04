package models

type Book struct {
	Id              int    `json:"id"`
	Title           string `json:"title"`
	AvailableCopies int    `json:"available_copies"`
}

type BookDetail struct {
	Title           string `json:"title"`
	AvailableCopies int    `json:"available_copies"`
}
