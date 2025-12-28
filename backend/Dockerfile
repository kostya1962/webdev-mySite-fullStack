FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY backend/ . 
RUN go build -o server ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY backend/images ./images
EXPOSE 3000
CMD ["./server"]