FROM public.ecr.aws/docker/library/golang:1.22.4-alpine3.20 AS builder

RUN apk add git
WORKDIR /src/app-build
ADD . .

RUN apk update && apk add bash ca-certificates git gcc g++

# Build
RUN \
  VERSION=$(date '+%Y%m%d.%H%M%S') && \
  COMMIT=$(git rev-parse HEAD) && \
  BRANCH=$(git rev-parse --abbrev-ref HEAD) && \
  HOST=$(hostname) && \
  GO111MODULE=on \
  GOOS=linux \
  GOARCH=amd64 \
  go build \
    -tags musl \
    -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH} -X main.BUILDHOST=${HOST}" \
    -o /go/bin/artemis ./cmd/serve

# ----  Now build final image  ----
FROM public.ecr.aws/docker/library/golang:1.22.4-alpine3.20

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ENV BIND_ADDRESS=0.0.0.0:9123

COPY --from=builder /go/bin/. /app/
WORKDIR /app

EXPOSE 9123
RUN date > BUILD_DATE
CMD ["./artemis"]
