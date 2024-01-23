package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func ConnectToMongoDB(ctx context.Context, dsn string) (*DB, error) {
	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return &DB{client: client}, nil
}

func (d *DB) Client() *mongo.Client {
	return d.client
}

func (d *DB) Close(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}