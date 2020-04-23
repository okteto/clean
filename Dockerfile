FROM golang:1.14 as builder
ARG COMMIT 
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . . 
RUN CGO=0 GOOS=linux go build -o remote -ldflags "-X main.CommitString=${COMMIT}" -tags "osusergo netgo static_build" .

FROM busybox

COPY --from=builder /app/clean /usr/local/bin/clean
RUN chmod +x /usr/local/bin/clean