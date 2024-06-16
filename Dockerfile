FROM golang:1.22.1-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-buildvcs=false go build -trimpath -ldflags "-w -s" -o app

FROM gcr.io/distroless/static-debian12 as deploy
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]

FROM golang:1.22.1-bullseye as dev
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app","80"]