FROM golang:1.17 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s -extldflags "-static"' -a -o /hacker_news .


FROM scratch

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /hacker_news /hacker_news

ENTRYPOINT ["/hacker_news"]