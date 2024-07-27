FROM gcr.io/distroless/base:latest
LABEL maintainer="Jonathan Wright <jon@than.io>" \
  org.opencontainers.image.title="dashboard" \
  org.opencontainers.image.description="A container for running dashboard, an event status web dashboard service." \
  org.opencontainers.image.authors="Jonathan Wright <jon@than.io>" \
  org.opencontainers.image.url="https://github.com/n3tuk/dashboard" \
  org.opencontainers.image.source="https://github.com/n3tuk/dashboard/blob/Dockerfile" \
  org.opencontainers.image.vendor="n3t.uk"

COPY dashboard /go/bin/dashboard

ENTRYPOINT ["/go/bin/dashboard"]
CMD ["serve"]
