package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"forms-service/dtos"
	"forms-service/services/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type formHandlerTestSuite struct {
	suite.Suite
	mockCtrl        *gomock.Controller
	mockFormService *mocks.MockFormService
	formHandler     *formHandler
}

func TestFormHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(formHandlerTestSuite))
}

func (suite *formHandlerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())

	suite.mockFormService = mocks.NewMockFormService(suite.mockCtrl)

	suite.formHandler = NewFormHandler(suite.mockFormService)
}

func (suite *formHandlerTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *formHandlerTestSuite) TestGetForm() {
	expected := &dtos.Form{
		FormId:    "",
		FormName:  "",
		FormPages: []dtos.FormPage{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "id", Value: "123"},
	}

	suite.mockFormService.EXPECT().GetForm("123", 0).Return(expected, nil).Times(1)
	suite.formHandler.GetForm(c)

	var result *dtos.Form

	err := json.Unmarshal(w.Body.Bytes(), &result)
	suite.NoError(err)

	suite.Equal(expected, result)
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *formHandlerTestSuite) TestGetFormError() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "id", Value: "123"},
	}

	suite.mockFormService.EXPECT().GetForm("123", 0).Return(nil, errors.New("error getting form")).Times(1)
	suite.formHandler.GetForm(c)

	var result *dtos.Form

	err := json.Unmarshal(w.Body.Bytes(), &result)
	suite.Error(err)
	suite.Empty(result)
	suite.Equal(http.StatusInternalServerError, w.Code)
}

func (suite *formHandlerTestSuite) TestCreateForm() {
	req := &dtos.Form{
		FormId:    "",
		FormName:  "",
		FormPages: []dtos.FormPage{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/form-masters", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.mockFormService.EXPECT().CreateForm(req).Return(nil).Times(1)

	suite.formHandler.CreateForm(c)
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *formHandlerTestSuite) TestCreateFormError() {
	req := &dtos.Form{
		FormId:    "",
		FormName:  "",
		FormPages: []dtos.FormPage{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/form-masters", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.mockFormService.EXPECT().CreateForm(req).Return(errors.New("error creating form")).Times(1)

	suite.formHandler.CreateForm(c)
	suite.Equal(http.StatusInternalServerError, w.Code)
}

func (suite *formHandlerTestSuite) TestCreateForm_BadRequest() {
	invalidJSON := `{"form_id": 123, "formName": "Test Form", "price": "not_a_number", "stock": "50"}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/form-masters", bytes.NewReader([]byte(invalidJSON)))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.formHandler.CreateForm(c)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *formHandlerTestSuite) TestDeleteForm() {
	suite.mockFormService.EXPECT().DeleteForm("123").Return(nil).Times(1)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "id", Value: "123"},
	}
	c.Request = httptest.NewRequest(http.MethodDelete, "/forms/123", nil)

	suite.formHandler.DeleteForm(c)
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *formHandlerTestSuite) TestDeleteFormError() {
	suite.mockFormService.EXPECT().DeleteForm("123").Return(errors.New("error deleting form")).Times(1)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "id", Value: "123"},
	}
	c.Request = httptest.NewRequest(http.MethodDelete, "/form-masters/123", nil)

	suite.formHandler.DeleteForm(c)
	suite.Equal(http.StatusInternalServerError, w.Code)
}

func (suite *formHandlerTestSuite) TestUpdateForm() {
	req := &dtos.Form{
		FormId:    "",
		FormName:  "",
		FormPages: []dtos.FormPage{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "id", Value: "123"},
	}
	c.Request = httptest.NewRequest(http.MethodPut, "/form-masters/123", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.mockFormService.EXPECT().UpdateForm("123", req).Return(nil).Times(1)

	suite.formHandler.UpdateForm(c)
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *formHandlerTestSuite) TestUpdateFormError() {
	req := &dtos.Form{
		FormId:    "",
		FormName:  "",
		FormPages: []dtos.FormPage{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "id", Value: "123"},
	}
	c.Request = httptest.NewRequest(http.MethodPut, "/form-masters/123", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.mockFormService.EXPECT().UpdateForm("123", req).Return(errors.New("error updating form")).Times(1)

	suite.formHandler.UpdateForm(c)
	suite.Equal(http.StatusInternalServerError, w.Code)
}

func (suite *formHandlerTestSuite) TestUpdateFormBadRequest() {
	invalidJSON := `{"form_id": 123, "formName": "Test Form", "price": "not_a_number", "stock": "50"}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "id", Value: "123"},
	}
	c.Request = httptest.NewRequest(http.MethodPut, "/forms/123", bytes.NewReader([]byte(invalidJSON)))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.formHandler.UpdateForm(c)
	suite.Equal(http.StatusBadRequest, w.Code)
}

/*

func (o *formHandler) UpsertFormDetails(ctx *gin.Context) {
	var req *dtos.FormDetails

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = o.formService.UpsertFormDetails(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
*/

func (suite *formHandlerTestSuite) TestGetFormDetails() {
	expected := &dtos.FormDetails{
		FormId:          "123",
		UserId:          "200",
		Status:          "",
		CurrentPage:     0,
		SubmittedAt:     time.Time{},
		CompletedFields: []dtos.Field{},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "formId", Value: "123"},
		{Key: "userId", Value: "200"},
	}

	suite.mockFormService.EXPECT().GetFormDetails("123", "200", 0).Return(expected, nil).Times(1)
	suite.formHandler.GetFormDetails(c)

	var result *dtos.FormDetails

	err := json.Unmarshal(w.Body.Bytes(), &result)
	suite.NoError(err)

	suite.Equal(expected, result)
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *formHandlerTestSuite) TestGetFormDetailsError() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "formId", Value: "123"},
		{Key: "userId", Value: "200"},
	}

	suite.mockFormService.EXPECT().GetFormDetails("123", "200", 0).Return(nil, errors.New("error getting form")).Times(1)
	suite.formHandler.GetFormDetails(c)

	var result *dtos.FormDetails

	err := json.Unmarshal(w.Body.Bytes(), &result)
	suite.Error(err)
	suite.Empty(result)
	suite.Equal(http.StatusInternalServerError, w.Code)
}

func (suite *formHandlerTestSuite) TestUpsertFormDetails() {
	req := &dtos.FormDetails{
		FormId:          "123",
		UserId:          "200",
		Status:          "",
		CurrentPage:     0,
		SubmittedAt:     time.Time{},
		CompletedFields: []dtos.Field{},
	}

	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "formId", Value: "123"},
		{Key: "userId", Value: "200"},
	}
	c.Request = httptest.NewRequest(http.MethodPut, "/forms/123/users/200", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.mockFormService.EXPECT().UpsertFormDetails(req).Return(nil).Times(1)

	suite.formHandler.UpsertFormDetails(c)
	suite.Equal(http.StatusOK, w.Code)
}

// func (suite *formHandlerTestSuite) TestUpdateFormError() {
// 	req := &dtos.Form{
// 		FormId:    "",
// 		FormName:  "",
// 		FormPages: []dtos.FormPage{},
// 		CreatedAt: time.Time{},
// 		UpdatedAt: time.Time{},
// 	}

// 	body, _ := json.Marshal(req)

// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)
// 	c.Params = []gin.Param{
// 		{Key: "id", Value: "123"},
// 	}
// 	c.Request = httptest.NewRequest(http.MethodPut, "/form-masters/123", bytes.NewReader(body))
// 	c.Request.Header.Set("Content-Type", "application/json")

// 	suite.mockFormService.EXPECT().UpdateForm("123", req).Return(errors.New("error updating form")).Times(1)

// 	suite.formHandler.UpdateForm(c)
// 	suite.Equal(http.StatusInternalServerError, w.Code)
// }

// func (suite *formHandlerTestSuite) TestUpdateFormBadRequest() {
// 	invalidJSON := `{"form_id": 123, "formName": "Test Form", "price": "not_a_number", "stock": "50"}`

// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)
// 	c.Params = []gin.Param{
// 		{Key: "id", Value: "123"},
// 	}
// 	c.Request = httptest.NewRequest(http.MethodPut, "/forms/123", bytes.NewReader([]byte(invalidJSON)))
// 	c.Request.Header.Set("Content-Type", "application/json")

// 	suite.formHandler.UpdateForm(c)
// 	suite.Equal(http.StatusBadRequest, w.Code)
// }
