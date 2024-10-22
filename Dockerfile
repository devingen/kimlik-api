FROM golang:latest
WORKDIR /app

ARG GIT_TOKEN
RUN go env -w GOPRIVATE=github.com/devingen
RUN git config --global url."https://golang:${GIT_TOKEN}@github.com".insteadOf "https://github.com"

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN env GOARCH=amd64 GOOS=linux go build -o api cmd/api/api.go

EXPOSE 1002

CMD ["./api"]
