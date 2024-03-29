# GophKeeper

GophKeeper is a secure storage application that allows users to store their private data, such as passwords, bank card information, and various binary and text data. This repository contains both the server and client applications, which communicate using the gRPC protocol. The server stores information in a PostgreSQL database, while the client maintains a local copy of data encrypted with the user's password. 

The client application provides functionality for user registration, authentication, retrieval of metadata about stored data, retrieval of data by ID, and creation and sending of new data to the server.

## Features

- User registration and authentication
- Retrieval of metadata about stored data
- Retrieval of data by ID
- Creation and submission of new data to the server

## Configuration

Both the server and client applications can be configured using environment variables, which are specified in the `.env` file.

## Getting Started

To get started with GophKeeper, follow the steps below:

1. Start the PostgreSQL container using Docker Compose:
   ```
   docker-compose up -d
   ```

2. Apply the database migrations:
   ```
    migrate -path ./migration/ -database "postgres://postgres:YourPassword@localhost:5432/postgres?sslmode=disable" -verbose up
   ```

3. Set the required environment variables in the `.env` file.

4. Run the server:
   ```
   go run ./cmd/server/main.go
   ```

5. Run the client:
   ```
   go run ./cmd/client/main.go
   ```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for details.

Feel free to contribute and provide feedback!