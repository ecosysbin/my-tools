# Build the manager binary
FROM harbor.zetyun.cn/gcp/golang-alpine3.19-tini:1.22.1 as builder

ENV GOPROXY https://goproxy.cn

# Make sure we use go modules
WORKDIR vcluster

# Copy the Go Modules manifests
COPY go.mod go.mod

# cache deps before building and copying source so that we don't need to re-download as much
RUN go mod download -x

COPY . .

# Build cmd
RUN CGO_ENABLED=0 GO111MODULE=on go build -o /plugin main.go

# we use alpine for easier debugging
FROM alpine

# Set root path as working directory
WORKDIR /

COPY --from=builder /plugin .

ENTRYPOINT ["/plugin"]
