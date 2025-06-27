package handlers

import (
	"forms-service/dtos"
	"forms-service/services/forms"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type formHandler struct {
	formService forms.FormService
}

func NewFormHandler(formService forms.FormService) *formHandler {
	return &formHandler{
		formService: formService,
	}
}

func (o *formHandler) GetForm(ctx *gin.Context) {
	id := ctx.Param("id")
	page, _ := strconv.Atoi(ctx.Query("page"))

	form, err := o.formService.GetForm(id, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, form)
}

func (o *formHandler) CreateForm(ctx *gin.Context) {
	var req *dtos.Form

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = o.formService.CreateForm(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "form created successfully"})
}

func (o *formHandler) DeleteForm(ctx *gin.Context) {
	id := ctx.Param("id")

	err := o.formService.DeleteForm(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "form deleted successfully"})
}

func (o *formHandler) UpdateForm(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dtos.Form
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = o.formService.UpdateForm(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Updated form successfully"})
}

func (o *formHandler) GetFormDetails(ctx *gin.Context) {
	formId := ctx.Param("formId")
	userId := ctx.Param("userId")
	page, _ := strconv.Atoi(ctx.Query("page"))

	form, err := o.formService.GetFormDetails(formId, userId, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, form)
}

func (o *formHandler) UpsertFormDetails(ctx *gin.Context) {
	formId := ctx.Param("formId")
	userId := ctx.Param("userId")

	var req *dtos.FormDetails
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	req.FormId = formId
	req.UserId = userId

	err = o.formService.UpsertFormDetails(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "form submitted successfully"})
}
