FROM golang:1.9.2-alpine
ADD . /go/src/websocket
RUN apk add --no-cache git
RUN go get -d -v github.com/gin-gonic/gin
RUN go get -d -v github.com/gorilla/websocket
RUN go install websocket

FROM alpine:latest
COPY --from=0 /go/bin/websocket .

CMD ["./websocket"]