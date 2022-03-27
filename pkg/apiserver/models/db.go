package models

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/teal-seagull/lyre-be-v4/pkg/config"
)

var (
	timeout  = 10 * time.Second
	instance *mongo.Client
	once     sync.Once
)

// TheClient is a singleton of DataBase connection pool
// Please be aware of this one is panics on connection failure
func TheClient() *mongo.Client {
	once.Do(func() {
		var err error

		if instance == nil {
			if instance, err = connectDatabase(); err != nil {
				panic(err)
			}
		}
	})

	return instance
}

// DB is just a shourtcut
func DB() *mongo.Database {
	return TheClient().Database(config.TheConfig().Database.Name)
}

func connectDatabase() (*mongo.Client, error) {
	var (
		client *mongo.Client
		ctx    context.Context
		cancel context.CancelFunc
		err    error
	)

	opt := options.Client().ApplyURI(config.TheConfig().Database.URI)

	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if client, err = mongo.Connect(ctx, opt); err != nil {
		return nil, fmt.Errorf("error connection to database - %s", err)
	}

	return client, nil
}
