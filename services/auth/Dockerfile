FROM golang:1.25-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o app ./cmd

EXPOSE 8080

CMD ["./app"]
