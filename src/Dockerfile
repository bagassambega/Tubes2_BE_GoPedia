FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]

# docker build -t gopedia-backend
# docker run -p 8080:8080 gopedia-backend