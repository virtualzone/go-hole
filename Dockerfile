FROM golang:1.19-bullseye AS server-builder
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD . .
RUN echo "package main\n\nconst AppVersion = \"`cat ./VERSION | awk NF`\"" > version.go
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o go-hole .

FROM gcr.io/distroless/static-debian11
COPY --from=server-builder /go/src/app/go-hole /app/
ADD config.yaml /app/
WORKDIR /app
EXPOSE 53
USER 65532:65532
CMD ["./go-hole"]