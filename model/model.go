package model

// Document struct
type Document struct {
	ID   string      `json:"id"`
	Body interface{} `json:"body"`
	Date string      `json:"date"`
}
