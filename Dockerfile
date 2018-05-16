FROM golang:1.10

WORKDIR /go/src/manager
COPY . .

RUN go get -d -v   # "go get -d -v ./..."
RUN go install -v    # "go install -v ./..."

EXPOSE 8080

CMD ["manager"]
