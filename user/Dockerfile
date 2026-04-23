FROM golang:1.25-alpine

RUN apk add --no-cache git bash curl

WORKDIR /app

# instala air
RUN go install github.com/air-verse/air@latest

CMD ["air"]