FROM golang:1.24.3-alpine3.21 AS build

WORKDIR /src/

COPY main.go go.* /src/
COPY internal/ /src/internal/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go mod tidy && go build -a -o /bin/exporter

FROM scratch

LABEL org.opencontainers.image.source="https://github.com/kyma-project/directory-size-exporter"

COPY --from=build /bin/exporter /bin/exporter

USER 65532:65532

ENTRYPOINT ["/bin/exporter"]
