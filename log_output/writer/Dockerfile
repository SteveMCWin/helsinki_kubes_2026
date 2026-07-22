FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod .
COPY main.go .
RUN CGO_ENABLED=0 go build -o /writer .

FROM scratch
COPY --from=build /writer /writer
ENTRYPOINT ["/writer"]
