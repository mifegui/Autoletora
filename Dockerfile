FROM golang

COPY ./app /go/src/github.com/mifegui/Autoletora/app
WORKDIR /go/src/github.com/mifegui/Autoletora/app

RUN go get ./
RUN go build

CMD ["./app"]
EXPOSE 8080
