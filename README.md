# README

## Introduction
This README provides an overview of the user management API and its endpoints.

## Endpoints

### 1. POST /reward/createUserItem
- Description: Create a user item to store user's item.
- Data:{
		"username": username,
		"amount": amouunt (amount of items, should be 0),
	},
- Response: create userItem status.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/createUserItem.png)

### 2. GET /reward/getUserItem?username=username
- Description: Retrieves a specific userItem by their username.
- parameter ?username=username
- Response: JSON object containing user details.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/getUserItem.png)
### 3. GET /reward/getAllUserItem
- Description: Retrieves all userItem .
- parameter ?username=username
- Response: JSON object containing user details.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/getAllUserItem.png)

### 4. PUT /reward/updateUserItem
- Description: Updates an existing userItem.
- Data:{
		"username": username,
		"amount": amouunt,
		"voucher": voucherID (string)
	},
- Response: Update failed or succeed
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/updateUserItem.png)
### 5. POST /reward/createGiftHistory
- Description: create a gift history between 2 users.
- Data:{
//			"sender": senderUsername,
//			"receiver": receiverUsername,
//			"amount": amount,
//		},
- Response: create status
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/createUserItem.png)


### 6. GET /reward/getSenderGiftHistory?username=username
- Description: Retrieves a specific sender giftHistory by their username.
- parameter ?username=username
- Response: JSON object containing details.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/getSenderGiftHistory.png)

### 7. GET /reward/getReceiverGiftHistory?username=username
- Description: Retrieves a specific receiver giftHistory by their username.
- parameter ?username=username
- Response: JSON object containing details.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/getReceiverGiftHistory.png)

### 8. GET /reward/getAllGiftHistory
- Description: Retrieves a all GiftHistory
- Response: JSON object containing details.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/getAllGiftHistory.png)


### 9. POST /reward/createExchange
- Description: create an exchange (user change item into voucher)
- Data:{
//			"username": username,
//			"voucher": voucher, string (ObjectID of the voucher)
//		}, 
- Response: create status
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/createExchange.png)


### 10. GET /reward/getExchange?username=username
- Description: get an exchange from a user by their username
- Response: JSON object containing details.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/getExchange.png)

### 11. GET /reward/getAllExchange
- Description: get All exchange 
- Response: JSON object containing details.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/reward-manegment/blob/main/image/getAllExchange.png)

