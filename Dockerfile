FROM golang:1.24.6-bookworm AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . . 
ARG COMMIT_SHA
RUN CGO=0 go build -o clean -ldflags "-X main.CommitString=${COMMIT_SHA}" -tags "osusergo netgo static_build" .

FROM busybox:1.37.0

COPY --from=builder /app/clean /usr/local/bin/clean
RUN chmod +x /usr/local/bin/clean