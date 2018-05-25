FROM golang:1.10

WORKDIR /go/src/engine-manager
COPY . .

RUN go install -v # "go install -v ./..."

EXPOSE 8080

CMD ["engine-manager"]
