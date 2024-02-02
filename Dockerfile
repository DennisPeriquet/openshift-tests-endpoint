FROM golang:latest

WORKDIR /app
COPY . .
RUN make build
ENTRYPOINT ["./endpoint_server"]
CMD ["-mode", "server"]
