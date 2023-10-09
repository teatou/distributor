FROM golang:1.21.0
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN cd cmd; go build -o main .
CMD ["/app/cmd/main"]