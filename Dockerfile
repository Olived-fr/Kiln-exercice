FROM golang:1.23.2-alpine AS build

ARG BUILD_TARGET=api

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o out ./cmd/${BUILD_TARGET}

FROM alpine

COPY --from=build /build/out /app

RUN chmod +x /app

EXPOSE 8080

CMD ["./app"]