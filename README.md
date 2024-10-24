# Kiln Exercise

## Overview

This project is a Go application designed to poll delegations from the Tezos blockchain using the TzKT API and store them in a PostgreSQL database.
By default, it will poll the delegations from the start and then poll the new delegations every 10 seconds.
It includes the delegation listing functionality.

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

1. Run the following command to start the application:
    ```sh
   make
    ```

2. The application will start polling delegations and storing them in the PostgreSQL database.

## Environment Variables

- `TZKT_URL`: The URL of the TzKT API. Default: https://api.tzkt.io.
- `POLLING_INTERVAL_SECONDS`: The interval in seconds at which the application polls for new delegations. Default: 10.
- `DEFAULT_POLLING_FROM`: The default start date for polling delegations. Format: YYYY-MM-DD. Default: 2018-01-01.
- `POLLING_BATCH_SIZE`: The number of delegations to fetch in each polling batch. Default: 10000.
- `POSTGRES_HOST`: The hostname of the PostgreSQL database.
- `POSTGRES_PORT`: The port number of the PostgreSQL database.
- `POSTGRES_USER`: The username for the PostgreSQL database.
- `POSTGRES_PASSWORD`: The password for the PostgreSQL database.
- `POSTGRES_DB`: The name of the PostgreSQL database.

## Running the Tests

1. Run the unit tests:
    ```sh
    go test ./...