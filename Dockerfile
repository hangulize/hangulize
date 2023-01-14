FROM golang:1.19.5-alpine3.17 AS builder
WORKDIR /hangulize
COPY . .

RUN apk add --update make git
RUN make -C /hangulize/cmd/hangulize

FROM alpine:3.17
COPY --from=builder /hangulize/cmd/hangulize/hangulize /bin/hangulize

ENTRYPOINT ["/bin/hangulize"]
