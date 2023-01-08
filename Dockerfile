FROM golang:1.19-bullseye AS builder
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD . .
RUN echo "package main\n\nconst AppVersion = \"`cat ./VERSION | awk NF`\"" > version.go
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o go-hole .
RUN apt-get update && apt-get install --yes libcap2-bin
RUN which setcap

FROM gcr.io/distroless/base-debian11
COPY --from=builder /go/src/app/go-hole /app/
COPY --from=builder /sbin/getcap /sbin/
COPY --from=builder /sbin/setcap /sbin/
COPY --from=builder /lib/*-linux-gnu/libcap.so.2 /lib/
RUN ["/sbin/setcap", "cap_net_bind_service=+ep", "/app/go-hole"]
ADD config.yaml /app/
WORKDIR /app
EXPOSE 53/udp
USER 65532:65532
CMD ["./go-hole"]