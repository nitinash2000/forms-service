package repository

import (
	"context"
	"fmt"
	"forms-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FormRepo interface {
	Create(form *models.Form) error
	Update(formId string, form *models.Form) error
	Get(formId string) (*models.Form, error)
	Delete(formId string) error

	UpsertFormSubmission(form *models.FormSubmission) error
	GetFormSubmission(formId, userId string) (*models.FormSubmission, error)
}

type formRepo struct {
	formCollection           *mongo.Collection
	formSubmissionCollection *mongo.Collection
}

func NewFormRepo(client *mongo.Client, database string) FormRepo {
	return &formRepo{
		formCollection:           client.Database(database).Collection("forms"),
		formSubmissionCollection: client.Database(database).Collection("form_submissions"),
	}
}

func (r *formRepo) Create(form *models.Form) error {
	_, err := r.formCollection.InsertOne(context.Background(), form)
	if err != nil {
		return fmt.Errorf("could not create form: %w", err)
	}
	return nil
}

func (r *formRepo) Update(formId string, form *models.Form) error {
	filter := bson.M{"form_id": formId}
	update := bson.M{"$set": form}

	_, err := r.formCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("could not update form: %w", err)
	}
	return nil
}

func (r *formRepo) Get(formId string) (*models.Form, error) {
	var form models.Form
	filter := bson.M{"form_id": formId}
	err := r.formCollection.FindOne(context.Background(), filter).Decode(&form)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("form with id %s not found", formId)
		}
		return nil, fmt.Errorf("could not retrieve form: %w", err)
	}
	return &form, nil
}

func (r *formRepo) Delete(formId string) error {
	filter := bson.M{"form_id": formId}
	_, err := r.formCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("could not delete form: %w", err)
	}
	return nil
}

func (r *formRepo) UpsertFormSubmission(form *models.FormSubmission) error {
	filter := bson.M{"form_id": form.FormId, "user_id": form.UserId}

	update := bson.M{
		"$set": form,
	}

	opts := options.Update().SetUpsert(true)

	_, err := r.formSubmissionCollection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return fmt.Errorf("could not upsert form submission: %w", err)
	}

	return nil
}

func (r *formRepo) GetFormSubmission(formId, userId string) (*models.FormSubmission, error) {
	var formSubmission models.FormSubmission
	filter := bson.M{"form_id": formId, "user_id": userId}

	err := r.formSubmissionCollection.FindOne(context.Background(), filter).Decode(&formSubmission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("form submission with id %s not found", formId)
		}
		return nil, fmt.Errorf("could not retrieve form: %w", err)
	}
	return &formSubmission, nil
}
