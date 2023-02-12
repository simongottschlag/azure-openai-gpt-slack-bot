FROM golang:1.19.5-bullseye as builder
WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY internal/ internal/

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -installsuffix "static" -o azure-openai-gpt-slack-bot .

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static /tini
RUN chmod +x /tini

# hadolint ignore=DL3008
RUN apt-get update && \
    apt-get -y --no-install-recommends install ca-certificates && \
    update-ca-certificates

FROM gcr.io/distroless/static-debian11:nonroot
LABEL org.opencontainers.image.source="https://github.com/simongottschlag/azure-openai-gpt-slack-bot"

WORKDIR /
COPY --from=builder /workspace/azure-openai-gpt-slack-bot /azure-openai-gpt-slack-bot
COPY --from=builder /tini /tini
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT [ "/tini", "--", "/azure-openai-gpt-slack-bot"]