FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /bin/server ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /bin/server /bin/server
EXPOSE 8080
CMD ["/bin/server"]
