# ---------- Build Stage ----------
FROM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS:-linux} \
    GOARCH=${TARGETARCH:-amd64}

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o amialive ./cmd

# ---------- Runtime Stage ----------
FROM gcr.io/distroless/base-debian12

COPY --from=ghcr.io/tarampampam/microcheck:1.3.0 /bin/httpcheck /bin/httpcheck
COPY --from=builder /app/amialive /app/amialive

WORKDIR /app

ENV JOKE_CHANCE=30
ENV PORT=8080

EXPOSE 8080

USER nonroot:nonroot

CMD ["/app/amialive"]
