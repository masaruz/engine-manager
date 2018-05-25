FROM golang:1.10

WORKDIR /go/src/engine-manager
COPY . .

RUN go install -v # "go install -v ./..."

EXPOSE 3000

CMD ["engine-manager"]
