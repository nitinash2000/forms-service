package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(r *gin.Engine, client *mongo.Client) {
	v1 := r.Group("/v1")

	formRoutes(v1, client)
	userRoutes(v1)
}
