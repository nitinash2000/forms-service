package forms

import (
	"forms-service/dtos"
	"forms-service/models"
	"forms-service/repository"
	"time"
)

type FormService interface {
	CreateForm(req *dtos.Form) error
	UpdateForm(id string, req *dtos.Form) error
	GetForm(formId string) (*dtos.Form, error)
	DeleteForm(formId string) error

	GetFormDetails(formId, userId string) (*dtos.FormDetails, error)
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

func (f *formService) GetForm(formId string) (*dtos.Form, error) {
	formModel, err := f.formRepo.Get(formId)
	if err != nil {
		return nil, err
	}

	var formPages []dtos.FormPage
	for _, v := range formModel.FormPages {
		var fields []dtos.Field
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

	result := &dtos.Form{
		FormId:    formModel.FormId,
		FormName:  formModel.FormName,
		FormPages: formPages,
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

func (f *formService) getCurrentPage(formPages []models.FormPage, completedFieldsMap map[string]struct{}) int {
	for _, formPage := range formPages {
		for _, field := range formPage.Fields {
			if _, exists := completedFieldsMap[field.FieldId]; !exists {
				return formPage.Page
			}
		}
	}

	return formPages[len(formPages)-1].Page
}

func (f *formService) GetFormDetails(formId, userId string) (*dtos.FormDetails, error) {
	form, err := f.formRepo.Get(formId)
	if err != nil {
		return nil, err
	}

	formSubmission, err := f.formRepo.GetFormSubmission(formId, userId)
	if err != nil {
		return nil, err
	}

	var fields []dtos.Field
	for _, v := range formSubmission.CompletedFields {
		fields = append(fields, dtos.Field{
			FieldId:    v.FieldId,
			FieldName:  v.FieldName,
			FieldType:  v.FieldType,
			FieldValue: v.FieldValue,
		})
	}

	result := &dtos.FormDetails{
		FormId:      formId,
		UserId:      userId,
		Status:      formSubmission.Status,
		CurrentPage: 1,
		SubmittedAt: formSubmission.SubmittedAt,
		Fields:      fields,
	}

	completedFieldsMap := make(map[string]struct{})
	for _, field := range formSubmission.CompletedFields {
		if field.FieldValue != "" {
			completedFieldsMap[field.FieldId] = struct{}{}
		}
	}

	result.CurrentPage = f.getCurrentPage(form.FormPages, completedFieldsMap)

	return result, nil
}

func (f *formService) UpsertFormDetails(req *dtos.FormDetails) error {
	var completedFields []models.Field

	for _, field := range req.Fields {
		completedFields = append(completedFields, models.Field{
			FieldId:    field.FieldId,
			FieldName:  field.FieldName,
			FieldType:  field.FieldType,
			FieldValue: field.FieldValue,
		})
	}

	formSubmissionModel := &models.FormSubmission{
		FormId:          req.FormId,
		UserId:          req.UserId,
		Status:          req.Status,
		CurrentPage:     req.CurrentPage,
		SubmittedAt:     time.Now(),
		CompletedFields: completedFields,
	}

	err := f.formRepo.UpsertFormSubmission(formSubmissionModel)
	if err != nil {
		return err
	}

	return nil
}
