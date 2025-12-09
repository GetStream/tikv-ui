# Stage 1: Build Frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY app/package.json app/yarn.lock* app/package-lock.json* ./
RUN npm install
COPY app .
# Allow setting API URL at build time, default to relative for same-origin serving
ENV NEXT_PUBLIC_API_URL=""
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o server cmd/main.go

# Stage 3: Final Image
FROM alpine:3:22
WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=backend-builder /app/server .
# Copy static frontend files to ./public which the Go server expects
COPY --from=frontend-builder /app/out ./public

EXPOSE 8081
CMD ["./server"]
