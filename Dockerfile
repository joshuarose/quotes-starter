## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY main.go quotes.sql ./

RUN go mod download

COPY *.go ./

RUN go build -o /quotes

## Deploy
FROM gcr.io/distroless/base

WORKDIR /

COPY --from=build /quotes /quotes

EXPOSE 8080

USER nonroot:nonroot

CMD ["/quotes"]