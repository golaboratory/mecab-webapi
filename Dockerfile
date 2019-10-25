FROM ubuntu:18.04

ENV DEBIAN_FRONTEND noninteractive
ENV GOPATH /golang
ENV CGO_LDFLAGS "-L/usr/lib/x86_64-linux-gnu -lmecab -lstdc++"
ENV CGO_CFLAGS "-I /usr/include "

RUN apt-get update
RUN apt-get install -y sudo golang mecab libmecab-dev mecab-ipadic-utf8 git make curl xz-utils file
RUN apt-get clean
RUN git clone --depth 1 https://github.com/neologd/mecab-ipadic-neologd.git
RUN bash /mecab-ipadic-neologd/bin/install-mecab-ipadic-neologd -y

RUN mkdir -p poo/mecab
COPY ./main.go poo/main.go
COPY ./mecab/mecab-bind.go poo/mecab/mecab-bind.go
RUN cd poo
RUN go get github.com/zenazn/goji
RUN go get github.com/rs/cors
RUN go get github.com/shogo82148/go-mecab
RUN ln -s /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so
ADD poo.conf /etc/init/poo.conf
RUN initctl reload-configuration
RUN go build -o poo/poo.bin poo/main.go
ENTRYPOINT ["poo/poo.bin"]
