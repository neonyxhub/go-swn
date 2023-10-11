FROM golang:1.21

ARG GO_FLAGS=-v

WORKDIR /app

# copy dependencies only to cache image layer
COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./

RUN go build -o /app "$GO_FLAGS" ./...

EXPOSE 8080
ENTRYPOINT ["/app/swn"]