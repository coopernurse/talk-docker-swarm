FROM golang:1.9
WORKDIR /go/src/github.com/coopernurse/swarm-demo-clock/
RUN go get -d -v github.com/dinever/golf
COPY clock.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.6 
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=0 /go/src/github.com/coopernurse/swarm-demo-clock/app .
CMD ["./app"]
