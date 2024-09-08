package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Build-D-An-Ki-n-Truc/reward-management/internal/db/mongodb"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message Pattern
/*
	{
	  "pattern": {
	    "service": "example-nestjs",
	    "endpoint": "hello",
	    "method": "GET"
	  },
	  "data": {
	    "headers": {},
	    "authorization": {},
	    "params": {
	      "name": "hai"
	    },
	    "payload": {}
	  },
	  "id": "5cb26e8dfd533783314c4"
	}
*/

type Pattern struct {
	Service  string `json:"service"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

type Header struct {
	Authorization string `json:"Authorization"`
}

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type Authorization struct {
	User User `json:"user"`
}

// In Data should have username and password for login
type Payload struct {
	Type   []string    `json:"type"`
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type Data struct {
	Headers       Header            `json:"headers"`
	Authorization Authorization     `json:"authorization"`
	Params        map[string]string `json:"params"`
	Payload       Payload           `json:"payload"`
}

// Struct for a Request
type Request struct {
	Pattern Pattern `json:"pattern"`
	Data    Data    `json:"data"`
	ID      string  `json:"id"`
}

// Struct for a Response

type Response struct {
	Headers       Header            `json:"headers"`
	Authorization Authorization     `json:"authorization"`
	Params        map[string]string `json:"params"`
	Payload       Payload           `json:"payload"`
}

func createSubscriptionString(endpoint, method, service string) string {
	return fmt.Sprintf(`{"endpoint":"%s","method":"%s","service":"%s"}`, endpoint, method, service)
}

// Subcriber for create a UserItem, Payload should have data:
//
//	Payload: Payload{
//		Data:{
//			"username": username,
//			"amount": amount, (initialy 0 or maybe a different number)
//		},
//	},
//
// reward/createUserItem/ POST	-> create a user
func CreateUserItemSubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("createUserItem", "POST", "reward")
	// Subscribe to users/create
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			// Get data from request
			// type assertion
			Data := request.Data.Payload.Data.(map[string]interface{})
			username := Data["username"].(string)
			amount := Data["amount"].(float64)

			// Create a new user
			NewUser := mongodb.UserItemStruct{
				Username: username,
				Amount:   int(amount),
				Voucher:  []primitive.ObjectID{},
			}

			err := mongodb.CreateUserItem(NewUser)
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to create userItem, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   "UserItem created",
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}
}

// Subcriber for get a UserItem, Request should have data:
//
// reward/getUserItem?username=...	-> get a user
//
//	reward/getUserItem GET	-> get a userItem
func GetOneUserItemSubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("getUserItem", "GET", "reward")
	// Subscribe to users/get
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			// Get data from request
			// type assertion
			Data := request.Data.Params
			username := Data["username"]
			// get that user
			targetUser, err := mongodb.ReadUserItem(username)
			// check if there is an error when getting that user
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to get userItem, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   targetUser,
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}
}

// Subcriber for get all UserItem, Request should have data:
//
// reward/getAllUserItem?
//
//	reward/getAllUserItem GET	-> get all userItem
func GetAllUserItemSubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("getAllUserItem", "GET", "reward")
	// Subscribe to users/get
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			targetUser, err := mongodb.ReadAllUserItem()
			// check if there is an error when getting that user
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to get userItem, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   targetUser,
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}
}

// Subcriber for update a UserItem, Payload should have data:
//
//	Payload: Payload{
//		Data:{
//			"username": username,
//			"amount": amount,
//			"voucher": voucher, string (ObjectID of the voucher)
//		},
//	},
//
// reward/updateUserItem/ PUT	-> update a user
func UpdateUserItemSubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("updateUserItem", "PUT", "reward")
	// Subscribe to users/update
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			// Get data from request
			// type assertion
			Data := request.Data.Payload.Data.(map[string]interface{})
			username := Data["username"].(string)
			amount := Data["amount"].(float64)
			voucher := Data["voucher"].(string)

			// Convert to ObjectID
			convertedVoucherID, err := primitive.ObjectIDFromHex(voucher)
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to convert voucherID, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
				return
			}

			err = mongodb.UpdateUserItem(username, convertedVoucherID, int(amount))
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to update userItem, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   "UserItem updated",
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
				return
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}

}

// Subcriber for create a GiftHistory, Payload should have data:
//
//	Payload: Payload{
//		Data:{
//			"sender": sender,
//			"receiver": receiver,
//			"amount": amount,
//		},
//	},
//
// reward/createGiftHistory/ POST	-> create a gift history
func CreateGiftHistorySubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("createGiftHistory", "POST", "reward")
	// Subscribe to users/create
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			// Get data from request
			// type assertion
			Data := request.Data.Payload.Data.(map[string]interface{})
			sender := Data["sender"].(string)
			receiver := Data["receiver"].(string)
			amount := Data["amount"].(float64)

			// Create a new user
			NewGiftHistory := mongodb.GiftHistoryStruct{
				Sender:   sender,
				Receiver: receiver,
				Time:     time.Now(),
				Amount:   int(amount),
			}

			err := mongodb.CreateGiftHistory(NewGiftHistory)
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to create giftHistory, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   "GiftHistory created",
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}
}

// Subcriber for get sender GiftHistory, Request should have data:
//
// reward/getSenderGiftHistory?username=...	-> get a gift history
//
//	reward/getSenderGiftHistory GET	-> get a gift history
func GetSenderGiftHistorySubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("getSenderGiftHistory", "GET", "reward")
	// Subscribe to users/get
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			// Get data from request
			// type assertion
			Data := request.Data.Params
			username := Data["username"]
			// get that user
			targetUser, err := mongodb.ReadSenderGiftHistory(username)
			// check if there is an error when getting that user
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to get giftHistory, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   targetUser,
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}
}

// Subcriber for get receiver GiftHistory, Request should have data:
//
// reward/getReceiverGiftHistory?username=...	-> get a gift history
//
//	reward/getReceiverGiftHistory GET	-> get a gift history
func GetReceiverGiftHistorySubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("getReceiverGiftHistory", "GET", "reward")
	// Subscribe to users/get
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			// Get data from request
			// type assertion
			Data := request.Data.Params
			username := Data["username"]
			// get that user
			targetUser, err := mongodb.ReadReceiverGiftHistory(username)
			// check if there is an error when getting that user
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to get giftHistory, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   targetUser,
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}
}

// Subcriber for get all GiftHistory, Request should have data:
//
// reward/getAllGiftHistory
//
//	reward/getAllGiftHistory GET	-> get all gift history
func GetAllGiftHistorySubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("getAllGiftHistory", "GET", "reward")
	// Subscribe to users/get
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			targetUser, err := mongodb.ReadAllGiftHistory()
			// check if there is an error when getting that user
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to get giftHistory, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   targetUser,
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}
}

// Subcriber for create a Exchange, Payload should have data:
//
//	Payload: Payload{
//		Data:{
//			"username": username,
//			"voucher": voucher, string (ObjectID of the voucher)
//		},
//	},
//
// reward/createExchange/ POST	-> create a exchange
func CreateExchangeSubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("createExchange", "POST", "reward")
	// Subscribe to users/create
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			// Get data from request
			// type assertion
			Data := request.Data.Payload.Data.(map[string]interface{})
			username := Data["username"].(string)
			voucher := Data["voucher"].(string)

			// Convert to ObjectID
			convertedVoucherID, err := primitive.ObjectIDFromHex(voucher)
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to convert voucherID, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
				return
			}

			// Create a new user
			NewExchange := mongodb.ExchangeStruct{
				Username: username,
				Time:     time.Now(),
				Voucher:  convertedVoucherID,
			}

			err = mongodb.CreateExchange(NewExchange)
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to create exchange, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   "Exchange created",
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}
}

// Subcriber for get a Exchange, Request should have data:
//
// reward/getExchange?username=...	-> get a exchange
//
//	reward/getExchange GET	-> get a exchange
func GetExchangeSubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("getExchange", "GET", "reward")
	// Subscribe to users/get
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			// Get data from request
			// type assertion
			Data := request.Data.Params
			username := Data["username"]
			// get that user
			targetUser, err := mongodb.ReadExchange(username)
			// check if there is an error when getting that user
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to get exchange, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   targetUser,
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}
}

// Subcriber for get all Exchange, Request should have data:
//
// reward/getAllExchange
//
//	reward/getAllExchange GET	-> get all exchange
func GetAllExchangeSubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("getAllExchange", "GET", "reward")
	// Subscribe to users/get
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			targetUser, err := mongodb.ReadAllExchange()
			// check if there is an error when getting that user
			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "Failed to get exchange, err: " + err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   targetUser,
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if errUser != nil {
		log.Println(errUser)
	}
}
