FROM fedora

RUN dnf update -y && dnf install -y fuse3-devel && dnf install -y podman

COPY ./image-processor ./image-processor

CMD ["./image-processor"]



