FROM golang:1.17-alpine

RUN apk add --no-cache bash openjdk11-jre curl clang llvm13 musl-dev

WORKDIR /opt/ruster/
COPY . .
RUN chmod 755 scripts/build.sh && \
    ./scripts/build.sh
RUN chmod 755 ./scripts/run.sh
RUN go build ./cmd/ruster/ruster.go
ENTRYPOINT ["/bin/bash"]
