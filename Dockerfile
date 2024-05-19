# Run superlist-backend in docker:
#
# Build as docker build -t sample-go-app .
# Run as docker run -it --rm -p 4001:4001 sample-go-app
#

FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go mod tidy
RUN go build -o sample-go-app.out

WORKDIR /app

ENTRYPOINT [ "sh", "./docker-entrypoint.sh" ]
# Ignore below, its only for quick debugging
# CMD [ "sleep", "infinity" ]
