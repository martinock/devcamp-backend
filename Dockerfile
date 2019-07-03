FROM golang:1.12

LABEL maintainer="Nino <martino.aksel.11@gmail.com>"

WORKDIR $GOPATH/src/hello

COPY . .

RUN go get -v 

EXPOSE 3000

CMD ["hello"]
