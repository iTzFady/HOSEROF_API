package config

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var ctx context.Context
var App *firebase.App
var DB *firestore.Client
var opt option.ClientOption

func InitFirebase() {
	ctx = context.Background()
	if firebaseCredsJSON := os.Getenv("FIREBASE_CREDENTIALS_JSON"); firebaseCredsJSON != "" {
		log.Println("Initializing Firebase from environment variable...")
		opt = option.WithCredentialsJSON([]byte(firebaseCredsJSON))
	} else {
		log.Println("Initializing Firebase from credentials file...")
		opt = option.WithCredentialsFile("./hoserof_fb.json")
	}

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Firebase: Error Occured: %v", err)
	}
	App = app

	client, err := App.Firestore(ctx)

	if err != nil {
		log.Fatalf("Firestore: Error Occured: %v", err)

	}
	DB = client
}
