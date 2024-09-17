FROM europe-docker.pkg.dev/kyma-project/prod/external/library/golang:1.23.1-alpine3.20 as build


WORKDIR /src/
COPY main.go go.* /src/
COPY internal/ /src/internal/
RUN CGO_ENABLED=0 go build -o /bin/exporter

FROM scratch
COPY --from=build /bin/exporter /bin/exporter
ENTRYPOINT ["/bin/exporter"]
