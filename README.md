# Reminder API

A serverless RESTful API for a simple reminder service.  

I am making this as an API for practicing making clients and welcome others to do the same.

## README Index
- [How is serverless done?](#how-is-serverless-done)
- [How do I use this API?](#how-do-i-use-this-api)
	- [Models](#models)
		- [Reminder](#reminder)
		- [User](#user)
		- [Credentials](#credentials)
	- [Authentication](#authentication)
		- [Login](#login)
		- [Register](#register)
	- [Reminder](#reminder-1)
		- [Get](#get)
		- [New](#new)
		- [Update](#update)
		- [Delete](#delete)

## How is serverless done?

For handling the HTTPS requests and executing all of the code I decided to use [Google Cloud Functions](https://cloud.google.com/functions/).  

And for the NoSQL database I used [Google Datastore](https://cloud.google.com/datastore/).  

I decided on using Google's cloud ecosystem as it seems to be reliable, fast, and cheap. Alongside being cheap, it comes with a free tier meaning I pay no money for hosting this API.

## How do I use this API?

As this is a Cloud Functions RESTful API the base path will be:  
`https://europe-west1-froogo-reminder-api.cloudfunctions.net/`  
  
All functions will take input and output in JSON and can respond with a variety of HTTP status codes.

### [Models](model/model.go)

#### Reminder
``` Go
type Reminder struct {
	Name, Description, Key             string `json:",omitempty"`
	Creation, Modification, Activation int64  `json:",omitempty"`
}
```

#### User
```Go
type User struct {
	Username, Email, Password string `json:",omitempty"`
}
```

#### Credentials
```Go
type Credentials struct {
	Auth, Refresh string `json:",omitempty"`
}
```

### Authentication

All secured end-points will required to be sent with a [Credentials](#credentials) object.  

#### Re-auth

At any point, if the auth token is invalid, the server will respond with the status code 401 (Unauthorized).  

If the refresh token is still valid, alongside this response there will be a new [Credentials](#credentials) object with a fresh auth and refresh token.  

Example:  

Status code: 401 (Unauthorized)  
Body:  
```JSON
{
	"Credentials": {
		"Auth": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjY5ODgxMDEsImlhdCI6MTU2NjkwMTcwMSwic3ViIjoiRWc4S0JGVnpaWElRZ0lDQXVJaVJnd28ifQ.chuoJkqKVCgouEQGsNMs00PdVnZTaMhV7BvaV0WfDlI",
		"Refresh": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjgxMTEzMDEsImp0aSI6IkVnOEtCRlZ6WlhJUWdJQ0F1SWlSZ3dvU0Rnb0RTbFJKRUlDQWdMalJ5WUlLIiwiaWF0IjoxNTY2OTAxNzAxLCJzdWIiOiJFZzhLQkZWelpYSVFnSUNBdUlpUmd3byJ9.MySQtyTRieM6PYIDIsNSTyheRu-1bmqkmlP5IyfWXQU"
	}
}
```

#### Login
Path: /login  
Input body: Email (string), Password (string)  
Example input:  
```JSON
{
    "Email": "harry@froogo.co.uk",
    "Password": "superSecretPassword123"
}
```
Response codes:  
```
200 (OK)                    - Success
400 (Bad Request)           - Invalid credentials
500 (Internal Server Error) - Internal server error
```
Response body: [Credentials](#credentials)  
Example output:  
```JSON
{
	"Credentials": {
		"Auth": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjY5ODgxMDEsImlhdCI6MTU2NjkwMTcwMSwic3ViIjoiRWc4S0JGVnpaWElRZ0lDQXVJaVJnd28ifQ.chuoJkqKVCgouEQGsNMs00PdVnZTaMhV7BvaV0WfDlI",
		"Refresh": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjgxMTEzMDEsImp0aSI6IkVnOEtCRlZ6WlhJUWdJQ0F1SWlSZ3dvU0Rnb0RTbFJKRUlDQWdMalJ5WUlLIiwiaWF0IjoxNTY2OTAxNzAxLCJzdWIiOiJFZzhLQkZWelpYSVFnSUNBdUlpUmd3byJ9.MySQtyTRieM6PYIDIsNSTyheRu-1bmqkmlP5IyfWXQU"
	}
}
```

#### Register
Path: /register  
Input body: [User](#user)  
Example input:  
```JSON
{
    "Email": "harry@froogo.co.uk",
    "Password": "superSecretPassword123"
}
```
Response codes:  
```
200 (OK)                    - Success
400 (Bad Request)           - Email taken, email invalid, or username and password blank
500 (Internal Server Error) - Internal server error
```
Response body: [Credentials](#credentials)  
Example output:  
```JSON
{
	"Credentials": {
		"Auth": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjY5ODgxMDEsImlhdCI6MTU2NjkwMTcwMSwic3ViIjoiRWc4S0JGVnpaWElRZ0lDQXVJaVJnd28ifQ.chuoJkqKVCgouEQGsNMs00PdVnZTaMhV7BvaV0WfDlI",
		"Refresh": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjgxMTEzMDEsImp0aSI6IkVnOEtCRlZ6WlhJUWdJQ0F1SWlSZ3dvU0Rnb0RTbFJKRUlDQWdMalJ5WUlLIiwiaWF0IjoxNTY2OTAxNzAxLCJzdWIiOiJFZzhLQkZWelpYSVFnSUNBdUlpUmd3byJ9.MySQtyTRieM6PYIDIsNSTyheRu-1bmqkmlP5IyfWXQU"
	}
}
```

### Reminder

All end-points in this section are secure and will require authentication.

#### Get
Path: /get  
Input body: [Credentials](#credentials)  
Example input:  
```JSON
{
	"Credentials": {
		"Auth": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjY5ODgxMDEsImlhdCI6MTU2NjkwMTcwMSwic3ViIjoiRWc4S0JGVnpaWElRZ0lDQXVJaVJnd28ifQ.chuoJkqKVCgouEQGsNMs00PdVnZTaMhV7BvaV0WfDlI",
		"Refresh": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjgxMTEzMDEsImp0aSI6IkVnOEtCRlZ6WlhJUWdJQ0F1SWlSZ3dvU0Rnb0RTbFJKRUlDQWdMalJ5WUlLIiwiaWF0IjoxNTY2OTAxNzAxLCJzdWIiOiJFZzhLQkZWelpYSVFnSUNBdUlpUmd3byJ9.MySQtyTRieM6PYIDIsNSTyheRu-1bmqkmlP5IyfWXQU"
	}
}
```
Response codes:  
```
200 (OK)                    - Success
401 (Unauthorized)          - Authentication expired, read authentication
500 (Internal Server Error) - Internal server error
```
Response body: Reminders (array of [Reminder](#reminder))  
Example output:  
```JSON
{
    "Reminders": [
        {
            "Name": "Get Milk",
            "Description": "Buy some ultra heat treated milk to end my existence.",
            "Key": "Eg8KBFVzZXIQgICAuIiRgwoSEwoIUmVtaW5kZXIQgICA2NrSiAo",
            "Creation": 1566920662,
            "Modification": 1566920662,
            "Activation": 1567516900
        }
    ]
}
```

#### New
Path: /new  
Input body: [Reminder](#reminder), [Credentials](#credentials)  
Example input:  
```JSON
{
	"Reminder": {
		"Name": "Get Milk",
		"Description": "Buy some ultra heat treated milk to end my existence.",
		"Activation": 1567516900
	},
	"Credentials": {
		"Auth": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjY5ODgxMDEsImlhdCI6MTU2NjkwMTcwMSwic3ViIjoiRWc4S0JGVnpaWElRZ0lDQXVJaVJnd28ifQ.chuoJkqKVCgouEQGsNMs00PdVnZTaMhV7BvaV0WfDlI",
		"Refresh": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjgxMTEzMDEsImp0aSI6IkVnOEtCRlZ6WlhJUWdJQ0F1SWlSZ3dvU0Rnb0RTbFJKRUlDQWdMalJ5WUlLIiwiaWF0IjoxNTY2OTAxNzAxLCJzdWIiOiJFZzhLQkZWelpYSVFnSUNBdUlpUmd3byJ9.MySQtyTRieM6PYIDIsNSTyheRu-1bmqkmlP5IyfWXQU"
	}
}
```
Response codes:  
```
200 (OK)                    - Success
401 (Unauthorized)          - Authentication expired, read authentication
500 (Internal Server Error) - Internal server error
```
Response body: Key (string)  
Example output:  
```JSON
{
	"Key": "Eg8KBFVzZXIQgICAuIiRgwoSEwoIUmVtaW5kZXIQgICA2NrSiAo"
}
```

#### Update
Path: /update  
Input body: [Reminder](#reminder), [Credentials](#credentials)  
Example input:  
```JSON
{
	"Reminder": {
		"Name": "Get Milk",
		"Description": "Buy some ultra heat treated milk to end my existence.",
		"Key": "Eg8KBFVzZXIQgICAuIiRgwoSEwoIUmVtaW5kZXIQgICA2NrSiAo",
		"Activation": 1567516900
	},
	"Credentials": {
		"Auth": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjY5ODgxMDEsImlhdCI6MTU2NjkwMTcwMSwic3ViIjoiRWc4S0JGVnpaWElRZ0lDQXVJaVJnd28ifQ.chuoJkqKVCgouEQGsNMs00PdVnZTaMhV7BvaV0WfDlI",
		"Refresh": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjgxMTEzMDEsImp0aSI6IkVnOEtCRlZ6WlhJUWdJQ0F1SWlSZ3dvU0Rnb0RTbFJKRUlDQWdMalJ5WUlLIiwiaWF0IjoxNTY2OTAxNzAxLCJzdWIiOiJFZzhLQkZWelpYSVFnSUNBdUlpUmd3byJ9.MySQtyTRieM6PYIDIsNSTyheRu-1bmqkmlP5IyfWXQU"
	}
}
```
Response codes:  
```
200 (OK)                    - Success
401 (Unauthorized)          - Authentication expired, read authentication
403 (Forbidden)             - Not owner of reminder
500 (Internal Server Error) - Internal server error
```
Response body: none provided, only response code.  

#### Delete
Path: /delete  
Input body: [Reminder](#reminder), [Credentials](#credentials)  
Example input:  
```JSON
{
	"Key": "Eg8KBFVzZXIQgICAuIiRgwoSEwoIUmVtaW5kZXIQgICA2NrSiAo",
	"Credentials": {
		"Auth": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjY5ODgxMDEsImlhdCI6MTU2NjkwMTcwMSwic3ViIjoiRWc4S0JGVnpaWElRZ0lDQXVJaVJnd28ifQ.chuoJkqKVCgouEQGsNMs00PdVnZTaMhV7BvaV0WfDlI",
		"Refresh": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjgxMTEzMDEsImp0aSI6IkVnOEtCRlZ6WlhJUWdJQ0F1SWlSZ3dvU0Rnb0RTbFJKRUlDQWdMalJ5WUlLIiwiaWF0IjoxNTY2OTAxNzAxLCJzdWIiOiJFZzhLQkZWelpYSVFnSUNBdUlpUmd3byJ9.MySQtyTRieM6PYIDIsNSTyheRu-1bmqkmlP5IyfWXQU"
	}
}
```
Response codes:  
```
200 (OK)                    - Success
401 (Unauthorized)          - Authentication expired, read authentication
403 (Forbidden)             - Not owner of reminder
500 (Internal Server Error) - Internal server error
```
Response body: none provided, only response code.  
