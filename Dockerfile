FROM golang:1.16

WORKDIR /
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
COPY ./public .

RUN go mod download
COPY . .
RUN go build -o main .

CMD ["./main"]
