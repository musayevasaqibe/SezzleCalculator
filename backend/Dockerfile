# Build stage
FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /sezzle-backend

# Runtime
FROM scratch
COPY --from=build /sezzle-backend /sezzle-backend
EXPOSE 8080
ENTRYPOINT ["/sezzle-backend"]
