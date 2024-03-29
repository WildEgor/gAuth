# Base Stage
FROM golang:1.21 as base
ARG APP_PATH=/app
COPY . $APP_PATH
WORKDIR $APP_PATH
RUN go mod download && mkdir -p dist

# Development Stage
FROM base as dev
WORKDIR /app
RUN go install -mod=mod github.com/cosmtrek/air
ENTRYPOINT ["air"]

# # Test Stage
# FROM base as test
# ENTRYPOINT make test

# Build Production Stage
FROM base as builder
ARG APP_PATH=/app
WORKDIR $APP_PATH
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dist/app cmd/main.go

# Production Stage
FROM alpine:latest as production
RUN apk --no-cache add ca-certificates
ARG APP_PATH=/root/
WORKDIR $APP_PATH
COPY --from=builder /app/dist/app .
EXPOSE 8888
CMD ["./app"]