FROM golang:1.21-alpine3.18 as builder

WORKDIR /build

COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o app cmd/main.go

FROM gcr.io/distroless/base-debian11
COPY --from=builder /build/app /build/app
USER nonroot:nonroot
ENTRYPOINT ["/build/app"]
