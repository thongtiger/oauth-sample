
#build stage
FROM golang:1.13-alpine AS builder
WORKDIR /app
# enable Go modules support
ENV GO111MODULE=on
# manage dependencies
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app


#final stage
FROM alpine:latest
WORKDIR /root/
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache tzdata
COPY --from=builder /app .
ENTRYPOINT ./app
LABEL Name=autopay Version=0.0.1
EXPOSE $PORT
