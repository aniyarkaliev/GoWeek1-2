FROM golang:latest

RUN mkdir /app
ADD . /app/
WORKDIR /app

EXPOSE 5432

EXPOSE 4000

# Build final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app .

CMD go run ./cmd/web