FROM quay.io/prometheus/golang-builder as builder

COPY ./promu /bin/promu
ADD . /go/src/github.com/pytimer/gpu_exporter
WORKDIR /go/src/github.com/pytimer/gpu_exporter
RUN /bin/promu build -v

FROM        quay.io/prometheus/busybox:glibc
MAINTAINER  Pytimer <lixin20101023@gmail.com>

ENV NVIDIA_VISIBLE_DEVICES=all
ENV NVIDIA_DRIVER_CAPABILITIES=utility

COPY --from=builder /lib/x86_64-linux-gnu/libdl-2.24.so /lib/libdl-2.24.so
COPY --from=builder /lib/x86_64-linux-gnu/libdl.so.2 /lib/libdl.so.2
COPY --from=builder /go/src/github.com/pytimer/gpu_exporter/gpu_exporter /bin/gpu_exporter

EXPOSE      9470
USER        nobody
ENTRYPOINT  [ "/bin/gpu_exporter" ]