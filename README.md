# Kiln Exercise

## Overview

This project is a Go application designed to poll delegations from the Tezos blockchain using the TzKT API and store them in a PostgreSQL database. It includes the delegation listing functionality.

## Project Structure

- `cmd/api/main.go`: Entry point for the API application.
- `cmd/polling/main.go`: Entry point for the polling application.
- `internal/model`: Contains the domain models.
- `internal/usecase/delegation/list`: Contains the use case and tests for listing delegations.
- `internal/usecase/delegation/poll`: Contains the use case for polling delegations.
- `internal/handler`: Contains the HTTP handlers for the API.
- `internal/pg`: Contains PostgreSQL repository implementations.
- `pkg/*`: Contains shared packages and utils (mostly taken from other personal projects).

## Running the Application

1. Build and run the application using Docker Compose:
    ```sh
    docker-compose up
    ```

2. The application will start polling delegations and storing them in the PostgreSQL database.

## Running Tests

1. Run the unit tests:
    ```sh
    go test ./...