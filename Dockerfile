FROM golang:1.9.2-alpine
ADD . /go/src/github.com/loogo/ws
RUN ls -al
RUN apk add --no-cache git
RUN go get -d -v github.com/gin-gonic/gin
RUN go get -d -v github.com/gorilla/websocket
RUN go install github.com/loogo/ws

FROM alpine:latest
COPY --from=0 /go/bin/ws .
RUN ls -al
RUN ls -al /go/src/github.com/loogo/ws
COPY  /go/src/github.com/loogo/ws/home.html .
CMD ["./ws"]