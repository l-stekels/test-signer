FROM golang:1.21 as base

# Dev container
FROM base as dev
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
WORKDIR /app
CMD ["air"]

# Build
FROM dev as build
ENV CGO_ENABLED=0
WORKDIR /app
RUN go install ./cmd
RUN go build -o /tmp/app ./cmd
