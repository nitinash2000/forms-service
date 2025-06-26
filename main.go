package main

import (
	"context"
	"encoding/json"
	"fmt"
	"forms-service/config"
	"forms-service/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConfig(filePath string) (*config.Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config config.Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	env := os.Getenv("env")
	if env == "" {
		env = "default"
	}

	config, err := GetConfig(fmt.Sprintf("%s.json", env))
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		os.Exit(1)
	}

	clientOptions := options.Client().ApplyURI(config.DbUrl)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}

	r := gin.Default()

	routes.Router(r, client)

	r.Run(":" + config.ServerPort)
}
