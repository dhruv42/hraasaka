package db

import (
	"context"
	"log"

	"github.com/dhruv42/hraasaka/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
}

type Client struct {
	Conn *mongo.Client
}

func GetDBConnection(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.DbUrl))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CreateIndexOnHash(ctx context.Context, cfg *config.Config, client *mongo.Client) error {
	col := client.Database(cfg.DbName).Collection("links")

	options := options.Index().SetUnique(true)
	mod := mongo.IndexModel{
		Keys: bson.M{
			"hash": 1,
		},
		Options: options,
	}

	idx, err := col.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return err
	}

	log.Printf("%s index created", idx)

	return nil
}
