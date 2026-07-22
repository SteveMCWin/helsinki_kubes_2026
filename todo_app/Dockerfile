FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum .
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o /todo_app .

FROM scratch
COPY --from=build /todo_app /todo_app
COPY front /front
ENTRYPOINT ["/todo_app"]
