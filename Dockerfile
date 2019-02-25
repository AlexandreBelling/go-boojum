FROM ubuntu:16.04

RUN apt-get update && apt-get install -y \
    wget \
    build-essential \
    libboost-all-dev \
    libssl-dev \ 
    cmake \
    libprocps-dev \
    libgmp-dev \
    pkg-config \
    software-properties-common

RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get install -y golang-1.11

RUN mkdir -p /usr/src/go/src/github.com/AlexandreBelling/go-boojum
COPY . /usr/src/go/src/github.com/AlexandreBelling/go-boojum
WORKDIR /usr/src/go/src/github.com/AlexandreBelling/go-boojum

RUN cd aggregator && make build-all
RUN ls

ENV GOPATH /usr/src/go
ENV GOROOT /usr/lib/go-1.11
ENV PATH $PATH:$GOROOT/bin

CMD cd scheduler && go test