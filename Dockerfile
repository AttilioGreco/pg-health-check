FROM golang:1.21-alpine3.19 as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go build -v -o pg-health-check .

FROM alpine:3.19

COPY --from=builder /usr/src/app/pg-health-check /usr/local/bin/pg-health-check

# CMD ["sleep", "3600"]
CMD ["/usr/local/bin/pg-health-check", "--config", "/etc/pg-health-check/config.yaml"]