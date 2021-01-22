FROM golang:1.15.6-alpine3.12 as build
WORKDIR /go/src/github.com/ProlificLabs/shakesearch

COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY go.mod .
COPY go.sum .

WORKDIR /go/src/github.com/ProlificLabs/shakesearch/cmd/search
RUN go build -o app

FROM alpine:3.12
RUN apk add ca-certificates tzdata
COPY --from=build /go/src/github.com/ProlificLabs/shakesearch/cmd/search .

ENTRYPOINT [ "./app" ]