FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum .
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o /ping_pong .

FROM scratch
COPY --from=build /ping_pong /ping_pong
ENTRYPOINT ["/ping_pong"]
