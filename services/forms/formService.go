package forms

import (
	"fmt"
	"forms-service/dtos"
	"forms-service/models"
	"forms-service/repository"
	"time"
)

type FormService interface {
	CreateForm(req *dtos.Form) error
	UpdateForm(id string, req *dtos.Form) error
	GetForm(formId string, page int) (*dtos.Form, error)
	DeleteForm(formId string) error

	GetFormDetails(formId, userId string, page int) (*dtos.FormDetails, error)
	UpsertFormDetails(req *dtos.FormDetails) error
}

type formService struct {
	formRepo repository.FormRepo
}

func NewFormService(formRepo repository.FormRepo) FormService {
	return &formService{
		formRepo: formRepo,
	}
}

func (f *formService) CreateForm(req *dtos.Form) error {
	var formPagesModel []models.FormPage

	for _, formPage := range req.FormPages {
		var fields []models.Field

		for _, v := range formPage.Fields {
			fields = append(fields, models.Field{
				FieldId:    v.FieldId,
				FieldName:  v.FieldName,
				FieldType:  v.FieldType,
				FieldValue: v.FieldValue,
			})
		}

		formPagesModel = append(formPagesModel, models.FormPage{
			Page:   formPage.Page,
			Fields: fields,
		})
	}

	formModel := &models.Form{
		FormId:    req.FormId,
		FormName:  req.FormName,
		FormPages: formPagesModel,
		CreatedAt: time.Now(),
	}

	err := f.formRepo.Create(formModel)
	if err != nil {
		return err
	}

	return nil
}

func (f *formService) UpdateForm(id string, req *dtos.Form) error {
	var formPagesModel []models.FormPage

	for _, formPage := range req.FormPages {
		var fields []models.Field

		for _, v := range formPage.Fields {
			fields = append(fields, models.Field{
				FieldId:    v.FieldId,
				FieldName:  v.FieldName,
				FieldType:  v.FieldType,
				FieldValue: v.FieldValue,
			})
		}

		formPagesModel = append(formPagesModel, models.FormPage{
			Page:   formPage.Page,
			Fields: fields,
		})
	}

	formModel := &models.Form{
		FormId:    req.FormId,
		FormName:  req.FormName,
		FormPages: formPagesModel,
		CreatedAt: req.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err := f.formRepo.Update(id, formModel)
	if err != nil {
		return err
	}

	return nil
}

func (f *formService) GetForm(formId string, page int) (*dtos.Form, error) {
	formModel, err := f.formRepo.Get(formId)
	if err != nil {
		return nil, err
	}

	var formPages []dtos.FormPage
	for _, v := range formModel.FormPages {
		var fields []dtos.Field

		if page == 0 || v.Page == page {
			for _, field := range v.Fields {
				fields = append(fields, dtos.Field{
					FieldId:   field.FieldId,
					FieldName: field.FieldName,
					FieldType: field.FieldType,
				})
			}

			formPages = append(formPages, dtos.FormPage{
				Page:   v.Page,
				Fields: fields,
			})
		}
	}

	result := &dtos.Form{
		FormId:    formModel.FormId,
		FormName:  formModel.FormName,
		FormPages: formPages,
		NoOfPages: len(formModel.FormPages),
		CreatedAt: formModel.CreatedAt,
		UpdatedAt: formModel.UpdatedAt,
	}

	return result, nil
}

func (f *formService) DeleteForm(formId string) error {
	err := f.formRepo.Delete(formId)
	if err != nil {
		return err
	}

	return nil
}

func (f *formService) getCurrentPageFields(formPages []models.FormPage, completedFieldsMap map[string]dtos.Field, page int) (int, string, []dtos.Field) {
	var completedFields []dtos.Field

	isCurrentPage := false
	nextField := ""

	for _, formPage := range formPages {
		completedFields = nil

		if page == formPage.Page {
			isCurrentPage = true
		}

		for _, field := range formPage.Fields {
			completedField, exists := completedFieldsMap[field.FieldId]
			if exists {
				completedFields = append(completedFields, completedField)
			} else if nextField == "" {
				nextField = field.FieldName
				isCurrentPage = true
			}
		}

		if isCurrentPage {
			return formPage.Page, nextField, completedFields
		}
	}

	return formPages[len(formPages)-1].Page, "", completedFields
}

func (f *formService) GetFormDetails(formId, userId string, page int) (*dtos.FormDetails, error) {
	form, err := f.formRepo.Get(formId)
	if err != nil {
		return nil, err
	}

	formSubmission, err := f.formRepo.GetFormSubmission(formId, userId)
	if err != nil {
		return nil, err
	}

	completedFieldsMap := make(map[string]dtos.Field)
	for _, field := range formSubmission.CompletedFields {
		if field.FieldValue != "" {
			completedFieldsMap[field.FieldId] = dtos.Field{
				FieldId:    field.FieldId,
				FieldName:  field.FieldName,
				FieldType:  field.FieldType,
				FieldValue: field.FieldValue,
			}
		}
	}

	result := &dtos.FormDetails{
		FormId:      formId,
		UserId:      userId,
		Status:      formSubmission.Status,
		CurrentPage: 1,
		SubmittedAt: formSubmission.SubmittedAt,
	}

	result.CurrentPage, result.NextField, result.CompletedFields = f.getCurrentPageFields(form.FormPages, completedFieldsMap, page)

	return result, nil
}

func (f *formService) getLastPageFields(form *models.Form, prevFormSubmission *models.FormSubmission) map[string]struct{} {
	prevFieldsMap := make(map[string]struct{})
	for _, v := range prevFormSubmission.CompletedFields {
		prevFieldsMap[v.FieldId] = struct{}{}
	}

	isCurrentPage := false

	for _, v := range form.FormPages {
		allowedFieldsMap := make(map[string]struct{})

		for _, field := range v.Fields {
			if _, exists := prevFieldsMap[field.FieldId]; !exists {
				isCurrentPage = true
			}

			allowedFieldsMap[field.FieldId] = struct{}{}
		}

		if isCurrentPage {
			return allowedFieldsMap
		}
	}

	return nil
}

func (f *formService) UpsertFormDetails(req *dtos.FormDetails) error {
	var completedFields []models.Field

	previousFormSubmission, err := f.formRepo.GetFormSubmission(req.FormId, req.UserId)
	if err != nil {
		return err
	}

	form, err := f.formRepo.Get(req.FormId)
	if err != nil {
		return err
	}

	allowedFieldsMap := f.getLastPageFields(form, previousFormSubmission)

	completedFieldsMap := make(map[string]struct{})
	for _, field := range req.CompletedFields {
		if _, exists := allowedFieldsMap[field.FieldId]; !exists {
			return fmt.Errorf("field %s not accepted", field.FieldId)
		}

		completedFieldsMap[field.FieldId] = struct{}{}

		completedFields = append(completedFields, models.Field{
			FieldId:    field.FieldId,
			FieldName:  field.FieldName,
			FieldType:  field.FieldType,
			FieldValue: field.FieldValue,
		})
	}

	for _, v := range previousFormSubmission.CompletedFields {
		if _, exists := completedFieldsMap[v.FieldId]; !exists {
			completedFields = append(completedFields, v)
		}
	}

	formSubmissionModel := &models.FormSubmission{
		FormId:          req.FormId,
		UserId:          req.UserId,
		Status:          req.Status,
		CurrentPage:     req.CurrentPage,
		SubmittedAt:     time.Now(),
		CompletedFields: completedFields,
	}

	err = f.formRepo.UpsertFormSubmission(formSubmissionModel)
	if err != nil {
		return err
	}

	return nil
}
