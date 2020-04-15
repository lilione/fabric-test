FROM hyperledger/fabric-baseimage:amd64-0.4.14

#dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg2 \
    software-properties-common \
    vim

#install Docker & Docker-compose
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
RUN add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
RUN apt-get update && apt-get install -y \
    docker-ce \
    docker-ce-cli \
    containerd.io
RUN curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
RUN chmod +x /usr/local/bin/docker-compose

#fabric golang dependencies
RUN go get github.com/syndtr/goleveldb/leveldb
ENV GOPATH /opt/gopath