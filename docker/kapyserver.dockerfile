FROM --platform=${BUILDPLATFORM} golang:1.22 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY ./kapyserver/ kapyserver/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o kapyserver kapyserver/cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/kapyserver .
USER 65532:65532

ENTRYPOINT ["/kapyserver"]
