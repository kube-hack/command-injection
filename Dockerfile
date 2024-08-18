FROM golang

COPY go.mod main.go ./

RUN go mod download && \
    go build -o main main.go && \
    apt-get update && \
    apt-get install -y iputils-ping

ENTRYPOINT [ "./main" ]