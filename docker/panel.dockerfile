# multi-stage docker build

# stage 1: build node/tailwindcss
FROM node:latest AS node
WORKDIR /app/panel/views
COPY panel/views/ ./
RUN npm install
RUN npx tailwindcss -i input.css -o static/style.css

# stage 2: generate templ files
FROM ghcr.io/a-h/templ:latest AS templ
WORKDIR /app
COPY --chown=65532:65532 --from=node /app/panel/views/ ./
RUN ["templ", "generate"]

# stage 3: build go binary
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd/ ./cmd/
COPY ./types/ ./types/
COPY ./panel/ ./panel/
COPY ./kapyclient/ ./kapyclient/
COPY ./log/ ./log/
COPY ./controller/ ./controller/
COPY ./ ./
COPY --from=templ /app panel/views/

RUN ls -la panel/views/auth
RUN mkdir out
RUN go build -o out/panel ./cmd/panel/main.go

# stage 4: final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/out/panel ./
EXPOSE 8080
ENTRYPOINT ["/root/panel"]
