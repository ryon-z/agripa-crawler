FROM golang:1.14

WORKDIR /crawler
COPY go.mod . 
COPY go.sum .

RUN apt-get update
RUN apt-get install vim -y
RUN apt-get install mariadb-client -y

RUN go mod download
