FROM golang:alpine as builder
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
COPY cars-fleet-service /
COPY cars.json /
COPY config.json /
ENTRYPOINT ["/cars-fleet-service"]
