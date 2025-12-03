FROM golang:1.25

WORKDIR /app
COPY ./web .
RUN go build -o main ./cmd/subsaggregator

EXPOSE 8080
ENTRYPOINT ["/app/main"]
