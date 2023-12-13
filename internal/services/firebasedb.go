package services

import (
	"context"
	"encoding/json"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

type FireDB struct {
	*db.Client
}

var fireDB FireDB

type FirebasePrivateKey map[string] string

func (db *FireDB) Connect(fpk FirebasePrivateKey) error {
	byteValue, err := json.Marshal(fpk)
	if err != nil {
		return err
	}

	ctx := context.Background()
	opt := option.WithCredentialsJSON(byteValue)
	config := &firebase.Config{DatabaseURL: "https://lan-coworking-bot-default-rtdb.firebaseio.com/"}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
	return fmt.Errorf("error initializing database: %v", err)
	}

	db.Client = client
	return nil
}

func FirebaseDB() *FireDB {
	return &fireDB
}