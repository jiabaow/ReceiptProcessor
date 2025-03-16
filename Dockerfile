FROM golang:1.20-alpine
LABEL authors="wenjiabao"

ENV GO111MODULE=on
ENV PORT=8080

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ReceiptProcessor .

EXPOSE 8080

CMD ["./ReceiptProcessor"]