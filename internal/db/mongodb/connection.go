package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/Build-D-An-Ki-n-Truc/reward-management/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is the MongoDB client object
var Client *mongo.Client

// ExchangeColl is the collection object for the Exchange collection
var ExchangeColl *mongo.Collection

// GiftHistoryColl is the collection object for the GiftHistory collection
var GiftHistoryColl *mongo.Collection

// UserItemColl is the collection object for the UserItem collection
var UserItemColl *mongo.Collection

// Initialize a connection to MongoDB
func InitializeMongoDBClient() error {
	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is canceled after the function returns

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DbUrl))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Test the connection
	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// Set the global client object
	Client = client
	ExchangeColl = Client.Database("RewardDB").Collection("exchange")
	GiftHistoryColl = Client.Database("RewardDB").Collection("gift_history")
	UserItemColl = Client.Database("RewardDB").Collection("user_item")

	return nil

}

// Disconnect from MongoDB
func DisconnectMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is canceled after the function returns

	// Disconnect the MongoDB client
	err := Client.Disconnect(ctx)

	if err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
	}

	return nil
}
