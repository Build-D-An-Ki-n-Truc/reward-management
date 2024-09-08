package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CRUD Operation for each collection

// CreateExchange creates a new ExchangeTransaction
func CreateExchange(exchange ExchangeStruct) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the new Exchange
	_, err := ExchangeColl.InsertOne(ctx, exchange)
	if err != nil {
		return err
	}
	return nil
}

// ReadExchange reads all exchanges from a user
func ReadExchange(username string) ([]ExchangeStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find all exchanges from the user with username
	cursor, err := ExchangeColl.Find(ctx, bson.M{"username": username})
	if err != nil {
		return nil, err
	}
	var exchanges []ExchangeStruct
	if err = cursor.All(ctx, &exchanges); err != nil {
		return nil, err
	}
	return exchanges, nil
}

// ReadAllExchange reads all exchanges from all users
func ReadAllExchange() ([]ExchangeStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ExchangeColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var exchanges []ExchangeStruct
	if err = cursor.All(ctx, &exchanges); err != nil {
		return nil, err
	}
	return exchanges, nil
}

// CreateGiftHistory creates a new GiftHistory
func CreateGiftHistory(giftHistory GiftHistoryStruct) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the new History
	_, err := GiftHistoryColl.InsertOne(ctx, giftHistory)
	if err != nil {
		return err
	}
	return nil
}

// ReadSenderGiftHistory reads all gift history from a user
func ReadSenderGiftHistory(username string) ([]GiftHistoryStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find all gift history from the user with username
	cursor, err := GiftHistoryColl.Find(ctx, bson.M{"sender": username})
	if err != nil {
		return nil, err
	}
	var giftHistory []GiftHistoryStruct
	if err = cursor.All(ctx, &giftHistory); err != nil {
		return nil, err
	}
	return giftHistory, nil
}

// ReadReceiverGiftHistory reads all gift history to a user
func ReadReceiverGiftHistory(username string) ([]GiftHistoryStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find all gift history to the user with username
	cursor, err := GiftHistoryColl.Find(ctx, bson.M{"receiver": username})
	if err != nil {
		return nil, err
	}
	var giftHistory []GiftHistoryStruct
	if err = cursor.All(ctx, &giftHistory); err != nil {
		return nil, err
	}
	return giftHistory, nil
}

// ReadAllGiftHistory reads all gift history
func ReadAllGiftHistory() ([]GiftHistoryStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := GiftHistoryColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var giftHistory []GiftHistoryStruct
	if err = cursor.All(ctx, &giftHistory); err != nil {
		return nil, err
	}
	return giftHistory, nil
}

// CreateUserItem creates a new UserItem
func CreateUserItem(userItem UserItemStruct) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// check if username already exists
	var userItemCheck UserItemStruct
	err := UserItemColl.FindOne(ctx, bson.M{"username": userItem.Username}).Decode(&userItemCheck)
	if err == nil {
		return errors.New("username already exists")
	}

	// Insert the new UserItem
	_, err = UserItemColl.InsertOne(ctx, userItem)
	if err != nil {
		return err
	}
	return nil
}

// ReadUserItem reads all user items from a user
func ReadUserItem(username string) ([]UserItemStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find all user items from the user with username
	cursor, err := UserItemColl.Find(ctx, bson.M{"username": username})
	if err != nil {
		return nil, err
	}
	var userItems []UserItemStruct
	if err = cursor.All(ctx, &userItems); err != nil {
		return nil, err
	}
	return userItems, nil
}

// ReadAllUserItem reads all user items
func ReadAllUserItem() ([]UserItemStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := UserItemColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var userItems []UserItemStruct
	if err = cursor.All(ctx, &userItems); err != nil {
		return nil, err
	}
	return userItems, nil
}

// UpdateUserItem updates a user item
func UpdateUserItem(username string, voucher primitive.ObjectID, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// check if the user item exists by find username
	var userItem UserItemStruct
	err := UserItemColl.FindOne(ctx, bson.M{"username": username}).Decode(&userItem)
	if err != nil {
		return err
	}

	// update the quantity of the user item and push the voucher
	_, err = UserItemColl.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"quantity": quantity}, "$push": bson.M{"voucher": voucher}})

	if err != nil {
		return err
	}
	return nil
}
