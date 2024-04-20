package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	CodeQueCollection    *mongo.Collection
	SubmissionCollection *mongo.Collection
	TemplateCollection   *mongo.Collection
)

var MONGO_URI = "mongodb+srv://lokesh:21341A0571@cluster0.yh7v13e.mongodb.net/coderealm_ms?retryWrites=true&w=majority&appName=Cluster0"

var Client = Connect()

func Connect() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// opts := options.Client().ApplyURI("mongodb+srv://lokeshkatari921:aI09VQyUrAm78tQS@cluster0.uhgrpv4.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// opts := options.Client().ApplyURI("mongodb://localhost:27017/coderealm").SetServerAPIOptions(serverAPI)
	opts := options.Client().ApplyURI(MONGO_URI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		fmt.Println(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	CodeQueCollection = client.Database("coderealm_ms").Collection("CodeQues")
	SubmissionCollection = client.Database("coderealm_ms").Collection("submissions")
	TemplateCollection = client.Database("coderealm_ms").Collection("Templates")
	// return client
	return client
}
