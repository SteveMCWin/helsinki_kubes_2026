FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod .
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o /add_todo_job .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /add_todo_job /add_todo_job
ENTRYPOINT ["/add_todo_job"]
