FROM golang:1.7.3-wheezy

run curl https://glide.sh/get | sh

run mkdir -p /go/src/github.com/jvikstedt/bluemoon
WORKDIR /go/src/github.com/jvikstedt/bluemoon

COPY glide.lock glide.yaml /go/src/github.com/jvikstedt/bluemoon/
run glide install

copy . /go/src/github.com/jvikstedt/bluemoon

ARG app=gate
run go install -v ./cmd/${app}
