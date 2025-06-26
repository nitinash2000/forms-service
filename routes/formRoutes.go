package routes

import (
	"forms-service/handlers"
	"forms-service/repository"
	"forms-service/services/forms"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func formRoutes(r *gin.RouterGroup, client *mongo.Client) {
	formRepo := repository.NewFormRepo(client, "forms-service")
	formService := forms.NewFormService(formRepo)
	formHandler := handlers.NewFormHandler(formService)

	forms := r.Group("/forms")
	{
		forms.GET("/:formId/users/:userId", formHandler.GetFormDetails)
		forms.POST("/:formId/users/:userId", formHandler.UpsertFormDetails)
	}

	formMasters := r.Group("/form-masters")
	{
		formMasters.GET("/:id", formHandler.GetForm)
		formMasters.POST("", formHandler.CreateForm)
		formMasters.PUT("/:id", formHandler.UpdateForm)
		formMasters.DELETE("/:id", formHandler.DeleteForm)
	}
}
