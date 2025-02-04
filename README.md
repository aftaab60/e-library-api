# e-Library API

This is a simple e-Library API built using Go and Gin, with in-memory storage for managing books and loans.
In memory can be switched to persistent DB (PostgreSQL) by changing repository layer which is an interface on purpose. 

## Features
- Retrieve book details and available copies
- Borrow a book (loan period: 4 weeks)
- Extend a loan (extend by 3 weeks from return date)
- Return a book

## Installation
Clone the repository and navigate into the project directory:

```sh
git clone https://github.com/aftaab60/e-library-api.git
cd e-library-api
```

Install dependencies:

```sh
go mod tidy
```

## Running the Project
Start the API server:

```sh
go run main.go
```

By default, the server runs on `localhost:3000`.

If using PostgreSQL DB, run docker command
```
docker-compose -f ./internal/docker/docker-compose.yml up -d 
```

## API Endpoints

### 1. Get Book Details
**GET /book/:title**

#### Example Request:
```sh
curl -X GET "http://localhost:3000/book/book1"
```

#### Response:
```json
{
  "title": "book1",
  "available_copies": 5
}
```

### 2. Borrow a Book
**POST /borrow**

#### Example Request:
```sh
curl --location 'localhost:3000/borrow' \
--header 'Content-Type: application/json' \
--data '{
    "title": "book1",
    "borrower_name": "user1"
}'
```

#### Response:
```json
{
  "name_of_borrower": "user1",
  "loan_date": "2025-02-03T16:17:53.439944+08:00",
  "return_date": "2025-03-03T16:17:53.439944+08:00"
}
```

### 3. Extend Loan
**POST /extend**

#### Example Request:
```sh
curl --location 'localhost:3000/extend' \
--header 'Content-Type: application/json' \
--data '{
    "title": "book1",
    "borrower_name": "user1"
}'
```

#### Response:
```json
{
  "name_of_borrower": "user1",
  "loan_date": "2025-02-03T16:17:53.439944+08:00",
  "return_date": "2025-03-24T16:17:53.439944+08:00"
}
```

### 4. Return a Book
**POST /return**

#### Example Request:
```sh
curl --location 'localhost:3000/return' \
--header 'Content-Type: application/json' \
--data '{
    "title": "book1",
    "borrower_name": "user1"
}'
```

#### Response:
```json
{
  "message": "book returned"
}
```

## Running Tests
To run unit tests:

```sh
go test ./...
```


