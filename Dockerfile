# Stage 1: Build Frontend (architecture-independent)
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-builder
WORKDIR /app
COPY app/package.json app/yarn.lock* app/package-lock.json* ./
RUN npm install
COPY app .
ENV NEXT_PUBLIC_API_URL=""
RUN npm run build

# Stage 2: Build Backend
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS backend-builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o server cmd/main.go

# Stage 3: Final Image
FROM alpine:3.23.0
WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /app/out ./public

EXPOSE 8081
CMD ["./server"]
