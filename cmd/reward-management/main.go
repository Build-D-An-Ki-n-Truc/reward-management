package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Build-D-An-Ki-n-Truc/reward-management/internal/db/mongodb"
	"github.com/Build-D-An-Ki-n-Truc/reward-management/internal/messaging/api"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func main() {
	url, exists := os.LookupEnv("NATS_URL")
	if !exists {
		url = nats.DefaultURL
	} else {
		url = strings.TrimSpace(url)
	}

	if strings.TrimSpace(url) == "" {
		url = nats.DefaultURL
	}

	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		logrus.Fatal(err)
		return
	}

	err = mongodb.InitializeMongoDBClient()

	if err != nil {
		logrus.Fatal(err)
	}

	// Subcribe to each service
	api.CreateExchangeSubcriber(nc)
	api.CreateGiftHistorySubcriber(nc)
	api.CreateUserItemSubcriber(nc)
	api.GetAllExchangeSubcriber(nc)
	api.GetAllGiftHistorySubcriber(nc)
	api.GetAllUserItemSubcriber(nc)
	api.GetExchangeSubcriber(nc)
	api.GetSenderGiftHistorySubcriber(nc)
	api.GetReceiverGiftHistorySubcriber(nc)
	api.GetOneUserItemSubcriber(nc)

	// Initialize MongoDB

	fmt.Println("Auth service running at port 3010")
	select {}
}
