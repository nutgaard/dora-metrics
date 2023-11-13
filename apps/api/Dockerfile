FROM golang:1.21.3-alpine as builder

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

COPY go.mod go.sum ./
RUN go mod download

COPY ../.. .

RUN go build -o dora-metrics-app
RUN ls -la
RUN chmod +x dora-metrics-app

#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/static-debian11

WORKDIR /

COPY --from=builder /app/dora-metrics-app /dora-metrics-app

EXPOSE 8080

ENTRYPOINT ["/dora-metrics-app"]