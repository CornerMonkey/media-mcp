FROM golang:1.26.1 AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN go test ./...

ARG MAIN_PACKAGE=./cmd/media-mcp
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/media-mcp ${MAIN_PACKAGE}

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/media-mcp /media-mcp

USER nonroot:nonroot
ENTRYPOINT ["/media-mcp"]
