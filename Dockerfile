FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /out/wc ./cmd/wc

FROM scratch
COPY --from=build /out/wc /wc
ENTRYPOINT ["/wc"]
