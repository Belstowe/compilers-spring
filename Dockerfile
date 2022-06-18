FROM golang:1.17-alpine

RUN apk add --no-cache bash openjdk11-jre

WORKDIR /opt/ruster/
COPY . .
RUN chmod 755 scripts/build.sh && \
    ./scripts/build.sh
RUN go build ./cmd/ruster/ruster.go
ENTRYPOINT ["/bin/bash"]