# README

## Introduction
This README provides an overview of the user management API and its endpoints.

## Endpoints

### 1. GET /users/getall
- Description: Retrieves a list of all users.
- Response: JSON array containing user objects.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/user-service-go/blob/main/image/getall.png)

### 2. GET /users/get
- Description: Retrieves a specific user by their ID.
- parameter ?username=username
- Response: JSON object containing user details.
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/user-service-go/blob/main/image/getone.png)
### 3. POST /users/create
- Description: Creates a new user.
- Data:{
		"username": username,
		"password": password,
		"name": name,
		"email": email,
		"role": role,
		"phone": phone,
		"isLocked": isLocked,
	},
- Response: Create failed or succeed
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/user-service-go/blob/main/image/create.png)
### 4. PUT /users/update
- Description: Updates an existing user.
- Data:{
		"username": username (same as existing user),
		"password": password,
		"name": name,
		"email": email,
		"role": role,
		"phone": phone,
		"isLocked": isLocked,
},
- Response: Update failed or succeed
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/user-service-go/blob/main/image/update.png)
### 5. DELETE /users/delete
- Description: Deletes a user.
- Data:{
		"username": username,
},
- Response: Delete failed or succeed.

## Placeholder Image
![Placeholder Image](https://github.com/Build-D-An-Ki-n-Truc/user-service-go/blob/main/image/delete.png)
