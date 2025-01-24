FROM golang:1.23.4
WORKDIR /app
COPY . .
RUN touch .env
RUN go mod download
RUN go build -o main
CMD ["./main"]