package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is a struct that represents an admin user
type ExchangeStruct struct {
	Username string             `json:"username" bson:"username"`
	Time     time.Time          `json:"time" bson:"time"`
	Voucher  primitive.ObjectID `json:"voucher" bson:"voucher"`
}

type GiftHistoryStruct struct {
	Sender   string    `json:"sender" bson:"sender"`
	Receiver string    `json:"receiver" bson:"receiver"`
	Time     time.Time `json:"time" bson:"time"`
	Amount   int       `json:"amount" bson:"amount"`
}

type UserItemStruct struct {
	Username string               `json:"username" bson:"username"`
	Amount   int                  `json:"amount" bson:"amount"`
	Voucher  []primitive.ObjectID `json:"voucher" bson:"voucher"`
}
