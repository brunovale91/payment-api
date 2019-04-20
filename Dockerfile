FROM golang:1.12.4-alpine3.9
RUN mkdir -p /go/src/github.com/brunovale91/payment-api 
ADD . /go/src/github.com/brunovale91/payment-api/
WORKDIR /go/src/github.com/brunovale91/payment-api 
RUN go build -o main .
RUN adduser -S -D -H -h /app appuser
USER appuser
CMD ["./main"]