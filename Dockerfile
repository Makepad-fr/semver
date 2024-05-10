FROM golang:1.22.2-bookworm as builder

WORKDIR /app

COPY go.work .
COPY ./cli/ ./cli/
COPY ./semver/ ./semver/
COPY Makefile Makefile

RUN make build

FROM scratch

COPY --from=builder /app/out/semver /

ENTRYPOINT [ "/semver" ]