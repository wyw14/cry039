FROM golang:1.24-alpine AS build-api
WORKDIR /office
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /office-feedback ./cmd/server
FROM alpine:3.21
RUN addgroup -S office && adduser -S office -G office
USER office
COPY --from=build-api /office-feedback /office-feedback
EXPOSE 8080
ENTRYPOINT ["/office-feedback"]
