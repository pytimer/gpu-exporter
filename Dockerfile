FROM quay.io/prometheus/golang-builder as builder

RUN go get -v -u github.com/prometheus/promu
ADD . /go/src/github.com/pytimer/gpu-exporter
WORKDIR /go/src/github.com/pytimer/gpu-exporter
RUN /go/bin/promu build -v

FROM quay.io/prometheus/busybox:glibc

ENV NVIDIA_VISIBLE_DEVICES=all
ENV NVIDIA_DRIVER_CAPABILITIES=utility

COPY --from=builder /lib/x86_64-linux-gnu/libdl-2.24.so /lib/libdl-2.24.so
COPY --from=builder /lib/x86_64-linux-gnu/libdl.so.2 /lib/libdl.so.2
COPY --from=builder /go/src/github.com/pytimer/gpu-exporter/gpu-exporter /bin/gpu-exporter

EXPOSE      9470
USER        nobody
ENTRYPOINT  [ "/bin/gpu-exporter" ]
