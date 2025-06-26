package dtos

import (
	"time"
)

type Form struct {
	FormId    string     `json:"form_id"`
	FormName  string     `json:"form_name"`
	FormPages []FormPage `json:"form_pages"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type FormPage struct {
	Page   int     `bson:"page"`
	Fields []Field `bson:"fields"`
}

type FormDetails struct {
	FormId      string    `json:"form_id"`
	UserId      string    `json:"user_id"`
	Status      string    `json:"status"`
	CurrentPage int       `json:"current_page"`
	SubmittedAt time.Time `json:"submitted_at"`
	Fields      []Field   `json:"fields"`
}

type Field struct {
	FieldId    string `json:"field_id"`
	FieldName  string `json:"field_name"`
	FieldType  string `json:"field_type"`
	FieldValue string `json:"field_value"`
}
