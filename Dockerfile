FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.19.2-alpine3.16 AS builder
ARG RELEASE_VERSION=devel
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
WORKDIR /go/src/github.com/mazay/mikromanager
RUN apk add git curl
COPY ./ ./
RUN go mod download
RUN go build

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.16.2
ARG TARGETPLATFORM
LABEL maintainer="Yevgeniy Valeyev <z.mazay@gmail.com>"
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY templates ./templates
COPY static ./static
COPY config.yml .
COPY --from=builder /go/src/github.com/mazay/mikromanager/mikromanager .
CMD ["./mikromanager"]
