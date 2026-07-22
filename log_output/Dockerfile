FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum .
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o /log_output .

FROM scratch
COPY --from=build /log_output /log_output
ENTRYPOINT ["/log_output"]
