package db

import (
	"context"
	"log"

	"github.com/dhruv42/hraasaka/config"
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

func InitDb(url string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatalf("Cou")
		panic(err)
	}

	return client
}

func (c *Client) Connect(ctx context.Context) error {
	err := c.Conn.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = c.Conn.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	return nil
}

func (c *Client) Disconnect(ctx context.Context) error {
	c.Conn.Disconnect(ctx)
	return nil
}
