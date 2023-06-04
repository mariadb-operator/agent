FROM golang:1.20.4-alpine3.18 AS builder

ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH}

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app
RUN go build -o agent main.go

FROM gcr.io/distroless/static AS app

WORKDIR /
COPY --from=builder /app/agent /bin/agent 
USER 65532:65532

ENTRYPOINT ["/bin/agent"]
