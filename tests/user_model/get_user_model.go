package main

// WARNING
// v4 database should be initialized with sample user/permission/application

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
	"github.com/teal-seagull/lyre-be-v4/pkg/config"
)

func main() {
	var (
		us   *models.UserScheme = &models.UserScheme{Name: "user1"}
		user *models.User
		b    []byte
		err  error
	)

	// Init config instance and override options
	config.TheConfig().Database.URI = "mongodb://localhost:27017"
	config.TheConfig().Database.Name = "v4"
	config.TheConfig().Database.Timeout = 1

	if user, err = models.NewUser(); err != nil {
		log.Fatalf("Error creating user model - %s\n", err)
	}

	if us, err = user.FindByName(us); err != nil {
		log.Fatalf("Error getting user - %s\n", err)
	}

	if b, err = json.MarshalIndent(us, "", "    "); err != nil {
		log.Fatalf("Error marshalling result - %s\n", b)
	}

	fmt.Printf("%s\n", string(b))
}
