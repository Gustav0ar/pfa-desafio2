FROM golang:1.16-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .

RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .

FROM scratch
COPY --from=builder /dist/main /

EXPOSE 3000

ENTRYPOINT ["/main"]