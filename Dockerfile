FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /server /app/server
COPY web /app/web
EXPOSE 8080
CMD ["/app/server"]
