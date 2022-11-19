FROM golang:1.18-alpine AS builder

WORKDIR /chat
COPY . .
RUN go get
RUN go build -o app

FROM alpine:3.17
RUN touch .env
COPY --from=builder /chat/app .

ENTRYPOINT [ "./app" ]