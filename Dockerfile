FROM --platform=${BUILDPLATFORM:-linux/amd64} public.ecr.aws/docker/library/golang:1.24.5-alpine3.21 AS builder
ARG RELEASE_VERSION=devel
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
WORKDIR /go/src/github.com/mazay/mikromanager
# hadolint ignore=DL3018
RUN apk --no-cache add git curl
COPY ./ ./
RUN go mod download
# hadolint ignore=DL3059
RUN go build

FROM public.ecr.aws/docker/library/alpine:3.21.3
ARG TARGETPLATFORM
LABEL maintainer="Yevgeniy Valeyev <z.mazay@gmail.com>"
# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates
# hadolint ignore=DL3059
RUN adduser \
    --disabled-password \
    --no-create-home \
    -u 8888 \
    mikromanager
USER mikromanager
WORKDIR /app/
COPY templates ./templates
COPY static ./static
COPY config.yml .
COPY --from=builder /go/src/github.com/mazay/mikromanager/mikromanager .
CMD ["./mikromanager"]
