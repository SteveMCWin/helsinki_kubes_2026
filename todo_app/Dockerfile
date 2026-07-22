FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod .
COPY main.go .
RUN CGO_ENABLED=0 go build -o /todo_app .

FROM scratch
COPY --from=build /todo_app /todo_app
ENTRYPOINT ["/todo_app"]
