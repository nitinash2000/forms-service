package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Form struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FormId    string             `bson:"form_id"`
	FormName  string             `bson:"form_name"`
	FormPages []FormPage         `bson:"form_pages"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type FormPage struct {
	Page   int     `bson:"page"`
	Fields []Field `bson:"fields"`
}

type Field struct {
	FieldId    string `bson:"field_id"`
	FieldName  string `bson:"field_name"`
	FieldType  string `bson:"field_type"`
	FieldValue string `bson:"field_value"`
}

type FormSubmission struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	FormId          string             `bson:"form_id"`
	UserId          string             `bson:"user_id"`
	Status          string             `bson:"status"`
	CurrentPage     int                `bson:"current_page"`
	SubmittedAt     time.Time          `bson:"submitted_at,omitempty"`
	CompletedFields []Field            `bson:"completed_fields"`
}
