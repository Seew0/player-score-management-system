FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache make git

RUN go mod download && go mod verify

RUN make

EXPOSE 4000

CMD ["./bin/api"]  