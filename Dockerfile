# Base Stage
FROM golang:1.22-alpine AS base

LABEL maintainer="Kartashov Egor <kartashov_egor96@mail.ru>"

ARG GITHUB_TOKEN
RUN apk update && apk add ca-certificates git openssh
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download && mkdir -p dist

# Development Stage
FROM base as dev
WORKDIR /app
COPY . .
RUN go install -mod=mod github.com/cosmtrek/air
ENTRYPOINT ["air"]

# # Test Stage
# FROM base as test
# ENTRYPOINT make test

# Build Production Stage
FROM base as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o dist/app cmd/main.go

# Production Stage
FROM cgr.dev/chainguard/busybox:latest-glibc as production
WORKDIR /app/
COPY --from=builder /app/dist/app .
COPY --from=builder /app/.env.local .
CMD ["./app"]