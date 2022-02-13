package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
	"github.com/teal-seagull/lyre-be-v4/pkg/config"
)

func main() {
	var (
		result struct {
			Version string `json:"version"`
		}
		status bson.D = bson.D{{Key: "serverStatus", Value: 1}}
	)

	// Init config instance and override options
	config.TheConfig().Database.URI = "mongodb://localhost:27017"
	config.TheConfig().Database.Name = "testing"
	config.TheConfig().Database.Timeout = 1

	_ = models.DB().RunCommand(context.TODO(), status).Decode(&result)

	if result.Version == "" {
		log.Fatalf("Error getting MongoDB version for %s:%s", config.TheConfig().Database.URI, config.TheConfig().Database.Name)
	}
}
