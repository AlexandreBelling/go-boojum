FROM ubuntu:16.04

# Install numerous dependencies
RUN apt-get update && apt-get install -y \
    wget \
    build-essential \
    libboost-all-dev \
    libssl-dev \ 
    cmake \
    libprocps-dev \
    libgmp-dev \
    pkg-config \
    software-properties-common \
    git

# Install golang 1.11
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get install -y golang-1.11
ENV GOPATH /usr/src/go
ENV GOROOT /usr/lib/go-1.11
ENV GO11MODULE=on
ENV PATH $PATH:$GOROOT/bin

# Build cpp dependencies
RUN mkdir -p /go-boojum/aggregator
WORKDIR /go-boojum/aggregator
COPY ./aggregator .
RUN make build-all

# Setup the module
WORKDIR /go-boojum/
COPY . .
RUN go mod download /go-boojum


